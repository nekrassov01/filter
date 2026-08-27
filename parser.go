package filter

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// Epsilon is a small value used to compare numerical equality.
const Epsilon = 1e-9

// MaxParen is the maximum number of opening '(' tokens allowed in one expression.
// Guards against pathological inputs causing excessive work. Counts total openings, not current depth.
const MaxParen = 256

// nodeBufSize is the number of nodes the parser stores inline before moving
// to a heap-allocated slice. Written by index rather than by append so that
// the parser stays on the stack of Parse.
const nodeBufSize = 16

// identBufSize is the number of distinct identifiers the parser stores inline
// before moving to a heap-allocated slice. Written by index for the same
// reason as nodeBufSize.
const identBufSize = 8

// nodeCharsEstimate is the assumed number of input bytes per node when
// sizing the overflow slice. It is an estimate, not a bound; append still
// grows the slice when the guess is low.
const nodeCharsEstimate = 8

// regexMap stores compiled regex patterns to reduce allocations on repeated parses.
// key: pattern string, value: *regexp.Regexp
var regexMap sync.Map

// Parse parses a string expression into an Expr.
func Parse(input string) (*Expr, error) {
	if input == "" {
		return nil, &Error{
			Kind: KindParse,
			Err:  fmt.Errorf("empty input"),
		}
	}
	p := parser{
		lexer: newLexer(input),
	}
	n, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.peek().typ != tokenEOF {
		return nil, &Error{
			Kind: KindParse,
			Err:  fmt.Errorf("unexpected token after parsing: %s", p.peek().v),
		}
	}
	nodes := p.nodes
	if nodes == nil {
		nodes = make([]node, p.nnode)
		copy(nodes, p.nodeBuf[:p.nnode])
	}
	return &Expr{
		nodes:  nodes,
		root:   n,
		nident: p.nident,
		shared: p.shared,
	}, nil
}

// parser represents a parser for the expression.
// It lives on the caller's stack for the duration of Parse; the fixed-size
// buffers let small expressions parse without heap growth.
type parser struct {
	lexer      lexer                // lexer for tokenizing input
	current    token                // current token
	peeked     bool                 // indicates if the next token has been peeked
	parenCount int                  // number of opening parentheses
	nodeBuf    [nodeBufSize]node    // expression tree nodes until nodeBuf is full
	nodes      []node               // all expression tree nodes once nodeBuf overflowed
	nnode      int                  // number of nodes
	identBuf   [identBufSize]string // distinct identifiers until identBuf is full
	idents     []string             // all distinct identifiers once identBuf overflowed
	nident     int                  // number of distinct identifiers
	shared     bool                 // some identifier is referenced more than once
}

// parseExpr parses an expression.
func (p *parser) parseExpr() (int, error) {
	left, err := p.parseAND()
	if err != nil {
		return 0, err
	}
	for {
		if p.peek().typ == tokenOR {
			t, err := p.next()
			if err != nil {
				return 0, err
			}
			right, err := p.parseAND()
			if err != nil {
				return 0, err
			}
			left = newNodeBinary(p, left, t, right)
			continue
		}
		break
	}
	return left, nil
}

// parseAND parses an AND expression.
func (p *parser) parseAND() (int, error) {
	left, err := p.parseNOT()
	if err != nil {
		return 0, err
	}
	for {
		if p.peek().typ == tokenAND {
			t, err := p.next()
			if err != nil {
				return 0, err
			}
			right, err := p.parseNOT()
			if err != nil {
				return 0, err
			}
			left = newNodeBinary(p, left, t, right)
			continue
		}
		break
	}
	return left, nil
}

// parseNOT parses a NOT expression.
func (p *parser) parseNOT() (int, error) {
	if p.peek().typ == tokenNOT {
		t, err := p.next()
		if err != nil {
			return 0, err
		}
		child, err := p.parsePrimary()
		if err != nil {
			return 0, err
		}
		return newNodeNOT(p, child, t), nil
	}
	return p.parsePrimary()
}

// parsePrimary parses a primary expression.
func (p *parser) parsePrimary() (int, error) {
	t := p.peek()
	switch t.typ {
	case tokenLparen:
		if _, err := p.next(); err != nil {
			return 0, err
		}
		p.parenCount++
		if p.parenCount > MaxParen {
			return 0, &Error{
				Kind: KindParse,
				Err:  fmt.Errorf("too many parentheses: exceeded limit %d at %d:%d", MaxParen, t.line, t.col),
			}
		}
		expr, err := p.parseExpr()
		if err != nil {
			return 0, err
		}
		if _, err := p.expect(tokenRparen); err != nil {
			return 0, err
		}
		return expr, nil
	case tokenIdent:
		return p.parseComparison()
	default:
		return 0, &Error{
			Kind: KindParse,
			Err:  fmt.Errorf("expected left parenthesis or identifier, got %s at %d:%d: %q", t.typ, t.line, t.col, t.v),
		}
	}
}

// parseComparison parses a comparison expression.
func (p *parser) parseComparison() (int, error) {
	ident, err := p.expect(tokenIdent)
	if err != nil {
		return 0, err
	}
	ident.idx = p.identIndex(ident.v)
	op, err := p.next()
	if err != nil {
		return 0, err
	}
	if !op.typ.isComparisonOperatorType() {
		return 0, &Error{
			Kind: KindParse,
			Err:  fmt.Errorf("expected comparison operator, got %s at %d:%d: %q", op.typ, op.line, op.col, op.v),
		}
	}
	val, err := p.next()
	if err != nil {
		return 0, err
	}
	if !val.typ.isValueType() {
		return 0, &Error{
			Kind: KindParse,
			Err:  fmt.Errorf("expected value, got %s at %d:%d: %q", val.typ, val.line, val.col, val.v),
		}
	}
	if val.typ == tokenString || val.typ == tokenRawString {
		val.v = unquote(val)
	}
	if op.typ.isCaseInsensitiveRegexOperatorType() {
		val.v = "(?i)" + val.v
	}
	i := newNodeComparison(p, ident, op, val)
	if op.typ.isRegexOperatorType() {
		if err := p.handleRegex(val, i); err != nil {
			return 0, err
		}
	}
	if val.typ == tokenTime {
		if t, err := time.Parse(time.RFC3339, val.v); err == nil {
			p.node(i).time = t
			p.node(i).hasTime = true
		}
	}
	if val.typ == tokenDuration {
		if d, err := time.ParseDuration(val.v); err == nil {
			p.node(i).dur = d
			p.node(i).hasDur = true
		}
	}
	if val.typ == tokenNumber {
		if f, err := strconv.ParseFloat(val.v, 64); err == nil {
			p.node(i).num = f
			p.node(i).hasNum = true
		}
	}
	return i, nil
}

// handleRegex processes a regex token and associates it with a node.
// Caches compiled regex patterns to reduce allocations on repeated parses.
func (p *parser) handleRegex(t token, i int) error {
	if t.v == "" {
		return &Error{
			Kind: KindParse,
			Err:  fmt.Errorf("invalid regex %q at %d:%d: empty pattern", t.v, t.line, t.col),
		}
	}
	if cached, ok := regexMap.Load(t.v); ok {
		p.node(i).re = cached.(*regexp.Regexp)
	} else {
		re, err := regexp.Compile(t.v)
		if err != nil {
			return &Error{
				Kind: KindParse,
				Err:  fmt.Errorf("invalid regex %q at %d:%d: %w", t.v, t.line, t.col, err),
			}
		}
		regexMap.Store(t.v, re)
		p.node(i).re = re
	}
	return nil
}

// identIndex returns the index of the identifier, registering it on first use.
func (p *parser) identIndex(name string) int {
	idents := p.idents
	if idents == nil {
		idents = p.identBuf[:p.nident]
	}
	for i, s := range idents {
		if s == name {
			p.shared = true
			return i
		}
	}
	i := p.nident
	switch {
	case p.idents != nil:
		p.idents = append(p.idents, name)
	case i < identBufSize:
		p.identBuf[i] = name
	default:
		p.idents = make([]string, i, 2*identBufSize)
		copy(p.idents, p.identBuf[:])
		p.idents = append(p.idents, name)
	}
	p.nident++
	return i
}

// addNode stores a node and returns its index.
func (p *parser) addNode(n node) int {
	i := p.nnode
	switch {
	case p.nodes != nil:
		p.nodes = append(p.nodes, n)
	case i < nodeBufSize:
		p.nodeBuf[i] = n
	default:
		// Estimate the final node count from the unread input so that large
		// expressions grow in one step instead of doubling repeatedly.
		remaining := len(p.lexer.input) - p.lexer.pos
		p.nodes = make([]node, i, max(2*nodeBufSize, i+remaining/nodeCharsEstimate))
		copy(p.nodes, p.nodeBuf[:])
		p.nodes = append(p.nodes, n)
	}
	p.nnode++
	return i
}

// node returns the node at index i.
func (p *parser) node(i int) *node {
	if p.nodes != nil {
		return &p.nodes[i]
	}
	return &p.nodeBuf[i]
}

// expect returns the next token and consumes it if it matches the expected type.
func (p *parser) expect(typ tokenType) (token, error) {
	t, err := p.next()
	if err != nil {
		return t, err
	}
	if t.typ != typ {
		return t, &Error{
			Kind: KindParse,
			Err:  fmt.Errorf("expected %s, got %s at %d:%d: %q", typ, t.typ, t.line, t.col, t.v),
		}
	}
	return t, nil
}

// next returns the next token from the lexer.
func (p *parser) next() (token, error) {
	if p.peeked {
		p.peeked = false
		if p.current.typ == tokenError {
			return p.current, &Error{
				Kind: KindLex,
				Err:  errors.New(p.current.v),
			}
		}
		return p.current, nil
	}
	p.current = p.lexer.nextToken()
	if p.current.typ == tokenError {
		return p.current, &Error{
			Kind: KindLex,
			Err:  errors.New(p.current.v),
		}
	}
	return p.current, nil
}

// peek returns the next token without consuming it.
func (p *parser) peek() token {
	if !p.peeked {
		p.current = p.lexer.nextToken()
		p.peeked = true
	}
	return p.current
}

// unquote removes the surrounding quotes from a string token.
func unquote(t token) string {
	n := len(t.v)
	if t.typ.isStringType() && n >= 2 {
		return t.v[1 : n-1]
	}
	return t.v
}
