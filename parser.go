package filter

// Experimental arena-based parser to reduce allocations.
// Keeps existing Parse() unchanged; use Parse() to try it.

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Epsilon is a small value used for floating-point comparison of numeric equality.
// Kept identical to previous evaluator implementation.
const Epsilon = 1e-9

// MaxParen is the maximum number of opening '(' tokens allowed in one expression.
// Guards against pathological inputs causing excessive work. Counts total openings, not current depth.
const MaxParen = 256

// regexMap stores compiled regex patterns to reduce allocations on repeated parses.
// key: pattern string, value: *regexp.Regexp
var regexMap sync.Map

// nodeType represents the type of a node in the expression tree.
type nodeType int

const (
	nodeComparison nodeType = iota // comparison node type
	nodeNot                        // logical NOT node type
	nodeBinary                     // binary operator node type
)

// String returns a string representation of the node type.
func (t nodeType) String() string {
	switch t {
	case nodeComparison:
		return "comparison node"
	case nodeNot:
		return "not node"
	case nodeBinary:
		return "binary node"
	}
	return ""
}

// node represents a node in the expression tree.
type node struct {
	typ    nodeType       // type of the node
	op     tokenType      // operator for binary and comparison nodes
	left   int            // left child index
	right  int            // right child index
	ident  string         // identifier for variable nodes
	val    string         // value for literal nodes
	re     *regexp.Regexp // regular expression for pattern matching
	num    float64        // cached numeric value
	dur    time.Duration  // cached duration value
	hasNum bool           // indicates if num is cached
	hasDur bool           // indicates if dur is cached
}

// parser represents a parser for the expression.
type parser struct {
	lexer   *lexer
	nodes   []node
	current token
	peeked  bool
	depth   int
	idents  map[string]struct{} // unique identifiers encountered (for field cache sizing)
}

// expr represents an expression in the parser.
type expr struct {
	parser *parser
	root   int
}

// Eval evaluates the expression against a target (per-call map cache for fields).
func (e *expr) Eval(t Target) (bool, error) {
	var cache map[string]any
	if len(e.parser.idents) > 0 {
		cache = make(map[string]any, len(e.parser.idents))
	}
	return e.parser.eval(e.root, t, cache)
}

// next returns the next token from the lexer.
func (p *parser) next() (token, error) {
	if p.peeked {
		p.peeked = false
		if p.current.typ == tokenError {
			return p.current, lexError(p.current.val)
		}
		return p.current, nil
	}
	p.current = p.lexer.nextToken()
	if p.current.typ == tokenError {
		return p.current, lexError(p.current.val)
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
		return t, parseError("expected %s, got %s at %d:%d: %q", typ, t.typ, t.line, t.col, t.val)
	}
	return t, nil
}

// parseExpr parses an expression.
func (p *parser) parseExpr() (int, error) {
	left, err := p.parseAND()
	if err != nil {
		return 0, err
	}
	for {
		if p.peek().typ == tokenOR {
			if _, err := p.next(); err != nil {
				return 0, err
			}
			right, err := p.parseAND()
			if err != nil {
				return 0, err
			}
			left = p.newNodeBinary(left, right, tokenOR)
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
			if _, err := p.next(); err != nil {
				return 0, err
			}
			right, err := p.parseNOT()
			if err != nil {
				return 0, err
			}
			left = p.newNodeBinary(left, right, tokenAND)
			continue
		}
		break
	}
	return left, nil
}

// parseNOT parses a NOT expression.
func (p *parser) parseNOT() (int, error) {
	if p.peek().typ == tokenNOT {
		if _, err := p.next(); err != nil {
			return 0, err
		}
		child, err := p.parsePrimary()
		if err != nil {
			return 0, err
		}
		return p.newNodeNot(child), nil
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
		p.depth++
		if p.depth > MaxParen {
			return 0, parseError("too many parentheses: exceeded limit %d at %d:%d", MaxParen, t.line, t.col)
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
		return 0, parseError("expected left parenthesis or identifier, got %s at %d:%d: %q", t.typ, t.line, t.col, t.val)
	}
}

// parseComparison parses a comparison expression.
func (p *parser) parseComparison() (int, error) {
	key, err := p.expect(tokenIdent)
	if err != nil {
		return 0, err
	}
	if p.idents != nil {
		p.idents[key.val] = struct{}{}
	}
	op, err := p.next()
	if err != nil {
		return 0, err
	}
	if !op.typ.isComparisonOperatorType() {
		return 0, parseError("expected comparison operator, got %s at %d:%d: %q", op.typ, op.line, op.col, op.val)
	}
	v, err := p.next()
	if err != nil {
		return 0, err
	}
	if !v.typ.isValueType() {
		return 0, parseError("expected value, got %s at %d:%d: %q", v.typ, v.line, v.col, v.val)
	}
	if op.typ.isCaseInsensitiveOperatorType() && !v.typ.isStringType() {
		return 0, parseError("expected numeric comparison operator, got string-only operator at %d:%d: %q", op.line, op.col, op.val)
	}
	val := v.val
	if v.typ == tokenString || v.typ == tokenRawString {
		val = unquote(v)
	}
	i := p.newNodeComparison(key.val, op.typ, val)
	if op.typ.isRegexOperatorType() {
		if val == "" {
			return 0, parseError("invalid regex %q at %d:%d: empty pattern", val, v.line, v.col)
		}
		if cached, ok := regexMap.Load(val); ok {
			p.nodes[i].re = cached.(*regexp.Regexp)
		} else {
			re, err := regexp.Compile(val)
			if err != nil {
				return 0, parseError("invalid regex %q at %d:%d: %w", val, v.line, v.col, err)
			}
			regexMap.Store(val, re)
			p.nodes[i].re = re
		}
	}
	if v.typ == tokenNumber {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			p.nodes[i].num = f
			p.nodes[i].hasNum = true
		}
	}
	if v.typ == tokenDuration {
		if d, err := time.ParseDuration(val); err == nil {
			p.nodes[i].dur = d
			p.nodes[i].hasDur = true
		}
	}
	return i, nil
}

// newNodeBinary creates a new binary expression node.
func (p *parser) newNodeBinary(left, right int, op tokenType) int {
	node := node{
		typ:   nodeBinary,
		op:    op,
		left:  left,
		right: right,
	}
	p.nodes = append(p.nodes, node)
	return len(p.nodes) - 1
}

// newNodeNot creates a new NOT expression node.
func (p *parser) newNodeNot(child int) int {
	node := node{
		typ:  nodeNot,
		op:   tokenNOT,
		left: child,
	}
	p.nodes = append(p.nodes, node)
	return len(p.nodes) - 1
}

// newNodeComparison creates a new comparison expression node.
func (p *parser) newNodeComparison(ident string, op tokenType, val string) int {
	node := node{
		typ:   nodeComparison,
		op:    op,
		ident: ident,
		val:   val,
	}
	p.nodes = append(p.nodes, node)
	return len(p.nodes) - 1
}

// eval evaluates the expression against a target.
func (p *parser) eval(i int, t Target, cache map[string]any) (bool, error) {
	n := p.nodes[i]
	switch n.typ {
	case nodeBinary:
		switch n.op {
		case tokenAND:
			left, err := p.eval(n.left, t, cache)
			if err != nil {
				return false, err
			}
			if !left {
				return false, nil
			}
			return p.eval(n.right, t, cache)
		case tokenOR:
			left, err := p.eval(n.left, t, cache)
			if err != nil {
				return false, err
			}
			if left {
				return true, nil
			}
			return p.eval(n.right, t, cache)
		default:
			return false, evalError("unsupported logical operator: %q", operators[n.op])
		}
	case nodeNot:
		v, err := p.eval(n.left, t, cache)
		if err != nil {
			return false, err
		}
		return !v, nil
	case nodeComparison:
		var field any
		if cache != nil {
			if v, ok := cache[n.ident]; ok {
				field = v
			} else {
				var err error
				field, err = t.GetField(n.ident)
				if err != nil {
					return false, evalError("%w", err)
				}
				cache[n.ident] = field
			}
		} else {
			var err error
			field, err = t.GetField(n.ident)
			if err != nil {
				return false, evalError("%w", err)
			}
		}
		switch v := field.(type) {
		case string:
			return p.evalString(n, v)
		case int:
			return p.evalNumber(n, float64(v))
		case int8:
			return p.evalNumber(n, float64(v))
		case int16:
			return p.evalNumber(n, float64(v))
		case int32:
			return p.evalNumber(n, float64(v))
		case int64:
			return p.evalNumber(n, float64(v))
		case uint:
			return p.evalNumber(n, float64(v))
		case uint8:
			return p.evalNumber(n, float64(v))
		case uint16:
			return p.evalNumber(n, float64(v))
		case uint32:
			return p.evalNumber(n, float64(v))
		case uint64:
			return p.evalNumber(n, float64(v))
		case float32:
			return p.evalNumber(n, float64(v))
		case float64:
			return p.evalNumber(n, v)
		case time.Duration:
			return p.evalDuration(n, v)
		default:
			return p.evalString(n, fmt.Sprint(v))
		}
	}
	return false, evalError("unsupported node type: %q", n.typ)
}

// evalString evaluates a string expression against a target.
func (p *parser) evalString(n node, v string) (bool, error) {
	switch n.op {
	case tokenEQ:
		return v == n.val, nil
	case tokenEQI:
		return strings.EqualFold(v, n.val), nil
	case tokenNEQ:
		return v != n.val, nil
	case tokenNEQI:
		return !strings.EqualFold(v, n.val), nil
	case tokenREQ:
		return n.re.MatchString(v), nil
	case tokenNREQ:
		return !n.re.MatchString(v), nil
	default:
		return false, evalError("unsupported operator for string: %q", operators[n.op])
	}
}

// evalNumber evaluates a number expression against a target.
func (p *parser) evalNumber(n node, v float64) (bool, error) {
	f := n.num
	if !n.hasNum {
		parsed, err := strconv.ParseFloat(n.val, 64)
		if err != nil {
			return false, evalError("invalid number: %q", n.val)
		}
		f = parsed
	}
	switch n.op {
	case tokenGT:
		return v > f, nil
	case tokenGTE:
		return v >= f, nil
	case tokenLT:
		return v < f, nil
	case tokenLTE:
		return v <= f, nil
	case tokenEQ:
		return math.Abs(v-f) <= Epsilon, nil
	case tokenNEQ:
		return math.Abs(v-f) > Epsilon, nil
	default:
		return false, evalError("unsupported operator for number: %q", operators[n.op])
	}
}

// evalDuration evaluates a duration expression against a target.
func (p *parser) evalDuration(n node, v time.Duration) (bool, error) {
	d := n.dur
	if !n.hasDur {
		parsed, err := time.ParseDuration(n.val)
		if err != nil {
			return false, evalError("invalid duration: %q", n.val)
		}
		d = parsed
	}
	switch n.op {
	case tokenGT:
		return v > d, nil
	case tokenGTE:
		return v >= d, nil
	case tokenLT:
		return v < d, nil
	case tokenLTE:
		return v <= d, nil
	case tokenEQ:
		return v == d, nil
	case tokenNEQ:
		return v != d, nil
	default:
		return false, evalError("unsupported operator for duration: %q", operators[n.op])
	}
}

// unquote removes the surrounding quotes from a string token.
func unquote(t token) string {
	var v string
	switch t.typ {
	case tokenString:
		if len(t.val) >= 2 {
			v = t.val[1 : len(t.val)-1]
		}
	case tokenRawString:
		if len(t.val) >= 2 {
			v = t.val[1 : len(t.val)-1]
		}
	}
	return v
}
