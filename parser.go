package filter

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// Epsilon is a small value used to compare numerical equality.
const Epsilon = 1e-9

// MaxParen is the maximum number of opening '(' tokens allowed in one expression.
// It guards against pathological inputs and counts total openings, not the
// current nesting depth.
const MaxParen = 256

// MaxInput is the maximum input length in bytes accepted by Parse.
// Positions and node indices are stored as int32; bounding the input keeps
// every int to int32 conversion in the lexer and parser within range.
const MaxInput = 1 << 20

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

// regexMap caches compiled regular expressions across parses. Keys are
// pattern strings and values are *regexp.Regexp.
var regexMap sync.Map

// Parse parses input into an Expr that can be evaluated repeatedly.
func Parse(input string) (*Expr, error) {
	if input == "" {
		return nil, &Error{
			Kind: KindParse,
			Err:  fmt.Errorf("empty input"),
		}
	}
	if len(input) > MaxInput {
		return nil, &Error{
			Kind: KindParse,
			Err:  fmt.Errorf("input too long: %d bytes exceeds limit %d", len(input), MaxInput),
		}
	}
	p := parser{
		lexer: newLexer(input),
	}
	n, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if t := p.peek(); t.typ != tokenEOF {
		return nil, newError(KindParse, t, "unexpected token after parsing: %s", t.v)
	}
	nodes := p.nodes
	if nodes == nil {
		nodes = make([]node, p.nnode)
		copy(nodes, p.nodeBuf[:p.nnode])
	}
	return &Expr{
		nodes:  nodes,
		root:   int32(n), //nolint:gosec // bounded by MaxInput
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

// parseExpr parses OR expressions, the lowest precedence level.
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

// parseAND parses AND expressions, which bind tighter than OR.
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

// parseNOT parses an optional NOT prefix followed by a primary expression.
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

// parsePrimary parses a parenthesized expression or a comparison.
func (p *parser) parsePrimary() (int, error) {
	t := p.peek()
	switch t.typ {
	case tokenLparen:
		if _, err := p.next(); err != nil {
			return 0, err
		}
		p.parenCount++
		if p.parenCount > MaxParen {
			return 0, newError(KindParse, t, "too many parentheses: exceeded limit %d", MaxParen)
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
		return 0, newError(KindParse, t, "expected left parenthesis or identifier, got %s: %q", t.typ, t.v)
	}
}

// parseComparison parses an identifier, a comparison operator, and a literal,
// validating typed literals as it goes.
func (p *parser) parseComparison() (int, error) {
	ident, err := p.expect(tokenIdent)
	if err != nil {
		return 0, err
	}
	ident.idx = int32(p.identIndex(ident.v)) //nolint:gosec // bounded by MaxInput
	op, err := p.next()
	if err != nil {
		return 0, err
	}
	if !op.typ.isComparisonOperatorType() {
		return 0, newError(KindParse, op, "expected comparison operator, got %s: %q", op.typ, op.v)
	}
	val, err := p.next()
	if err != nil {
		return 0, err
	}
	if !val.typ.isValueType() {
		return 0, newError(KindParse, val, "expected value, got %s: %q", val.typ, val.v)
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
	// Typed literals are validated here, once; string literals compared
	// against typed values are converted at evaluation time instead.
	switch val.typ {
	case tokenTime:
		t, err := time.Parse(time.RFC3339, val.v)
		if err != nil {
			return 0, newError(KindParse, val, "invalid time %q", val.v)
		}
		p.node(i).time = t
		p.node(i).hasTime = true
	case tokenDuration:
		d, err := time.ParseDuration(val.v)
		if err != nil {
			return 0, newError(KindParse, val, "invalid duration %q", val.v)
		}
		p.node(i).dur = d
		p.node(i).hasDur = true
	case tokenNumber:
		f, err := strconv.ParseFloat(val.v, 64)
		if err != nil {
			return 0, newError(KindParse, val, "invalid number %q", val.v)
		}
		p.node(i).num = f
		p.node(i).hasNum = true
	}
	return i, nil
}

// handleRegex compiles the regex pattern in t and stores it on node i.
// Compiled patterns are cached in regexMap to reduce allocations on repeated parses.
func (p *parser) handleRegex(t token, i int) error {
	if t.v == "" {
		return newError(KindParse, t, "invalid regex %q: empty pattern", t.v)
	}
	if cached, ok := regexMap.Load(t.v); ok {
		p.node(i).re = cached.(*regexp.Regexp)
	} else {
		re, err := regexp.Compile(t.v)
		if err != nil {
			return newError(KindParse, t, "invalid regex %q: %w", t.v, err)
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
		return t, newError(KindParse, t, "expected %s, got %s: %q", typ, t.typ, t.v)
	}
	return t, nil
}

// next consumes and returns the next token, reporting lexer errors as Error.
func (p *parser) next() (token, error) {
	if p.peeked {
		p.peeked = false
		if p.current.typ == tokenError {
			return p.current, newError(KindLex, p.current, "%s", p.current.v)
		}
		return p.current, nil
	}
	p.current = p.lexer.nextToken()
	if p.current.typ == tokenError {
		return p.current, newError(KindLex, p.current, "%s", p.current.v)
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

// unquote returns the text of a string token without its surrounding quotes.
func unquote(t token) string {
	n := len(t.v)
	if t.typ.isStringType() && n >= 2 {
		return t.v[1 : n-1]
	}
	return t.v
}
