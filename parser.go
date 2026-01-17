package filter

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// Parse parses a string expression into an Expr.
func Parse(input string) (Expr, error) {
	p, err := newParser(input)
	if err != nil {
		return Expr{}, err
	}
	n, err := p.parseExpr()
	if err != nil {
		return Expr{}, err
	}
	if p.peek().typ != tokenEOF {
		return Expr{}, &FilterError{
			Kind: KindParse,
			Err:  fmt.Errorf("unexpected token after parsing: %s", p.peek().v),
		}
	}
	return Expr{
		parser: p,
		root:   n,
	}, nil
}

// Epsilon is a small value used to compare numerical equality.
const Epsilon = 1e-9

// MaxParen is the maximum number of opening '(' tokens allowed in one expression.
// Guards against pathological inputs causing excessive work. Counts total openings, not current depth.
const MaxParen = 256

// regexMap stores compiled regex patterns to reduce allocations on repeated parses.
// key: pattern string, value: *regexp.Regexp
var regexMap sync.Map

// parser represents a parser for the expression.
type parser struct {
	lexer      lexer               // lexer for tokenizing input
	nodes      []node              // expression tree nodes
	current    token               // current token
	peeked     bool                // indicates if the next token has been peeked
	parenCount int                 // Number of opening parentheses
	idents     map[string]struct{} // Unique identifier encountered in field cache size settings
}

// newParser creates a new parser for the given input.
func newParser(input string) (parser, error) {
	if input == "" {
		return parser{}, &FilterError{
			Kind: KindParse,
			Err:  fmt.Errorf("empty input"),
		}
	}
	return parser{
		lexer:  newLexer(input),
		nodes:  make([]node, 0, 16),
		idents: make(map[string]struct{}),
	}, nil
}

// next returns the next token from the lexer.
func (p *parser) next() (token, error) {
	if p.peeked {
		p.peeked = false
		if p.current.typ == tokenError {
			return p.current, &FilterError{
				Kind: KindLex,
				Err:  errors.New(p.current.v),
			}
		}
		return p.current, nil
	}
	p.current = p.lexer.nextToken()
	if p.current.typ == tokenError {
		return p.current, &FilterError{
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

// expect returns the next token and consumes it if it matches the expected type.
func (p *parser) expect(typ tokenType) (token, error) {
	t, err := p.next()
	if err != nil {
		return t, err
	}
	if t.typ != typ {
		return t, &FilterError{
			Kind: KindParse,
			Err:  fmt.Errorf("expected %s, got %s at %d:%d: %q", typ, t.typ, t.line, t.col, t.v),
		}
	}
	return t, nil
}

// unquote removes the surrounding quotes from a string token.
func unquote(t token) string {
	n := len(t.v)
	if t.typ.isStringType() && n >= 2 {
		return t.v[1 : n-1]
	}
	return t.v
}

// handleRegex processes a regex token and associates it with a node.
// Caches compiled regex patterns to reduce allocations on repeated parses.
func (p *parser) handleRegex(t token, i int) error {
	if t.v == "" {
		return &FilterError{
			Kind: KindParse,
			Err:  fmt.Errorf("invalid regex %q at %d:%d: empty pattern", t.v, t.line, t.col),
		}
	}
	if cached, ok := regexMap.Load(t.v); ok {
		p.nodes[i].re = cached.(*regexp.Regexp)
	} else {
		re, err := regexp.Compile(t.v)
		if err != nil {
			return &FilterError{
				Kind: KindParse,
				Err:  fmt.Errorf("invalid regex %q at %d:%d: %w", t.v, t.line, t.col, err),
			}
		}
		regexMap.Store(t.v, re)
		p.nodes[i].re = re
	}
	return nil
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
			return 0, &FilterError{
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
		return 0, &FilterError{
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
	if p.idents != nil {
		p.idents[ident.v] = struct{}{}
	}
	op, err := p.next()
	if err != nil {
		return 0, err
	}
	if !op.typ.isComparisonOperatorType() {
		return 0, &FilterError{
			Kind: KindParse,
			Err:  fmt.Errorf("expected comparison operator, got %s at %d:%d: %q", op.typ, op.line, op.col, op.v),
		}
	}
	val, err := p.next()
	if err != nil {
		return 0, err
	}
	if !val.typ.isValueType() {
		return 0, &FilterError{
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
			p.nodes[i].time = t
			p.nodes[i].hasTime = true
		}
	}
	if val.typ == tokenDuration {
		if d, err := time.ParseDuration(val.v); err == nil {
			p.nodes[i].dur = d
			p.nodes[i].hasDur = true
		}
	}
	if val.typ == tokenNumber {
		if f, err := strconv.ParseFloat(val.v, 64); err == nil {
			p.nodes[i].num = f
			p.nodes[i].hasNum = true
		}
	}
	return i, nil
}
