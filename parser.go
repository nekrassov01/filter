package filter

import (
	"net/netip"
	"regexp"
	"strings"
	"sync"
	"time"
)

// MaxParen is the maximum number of opening parentheses in one expression.
const MaxParen = 256

// MaxInput is the maximum input length in bytes accepted by Parse. Positions,
// columns, and node indices are int32; this bound keeps them in range.
const MaxInput = 1 << 20

// nodeBufSize is the number of nodes a parser holds inline.
const nodeBufSize = 16

// identBufSize is the number of distinct identifiers a parser holds inline.
const identBufSize = 8

// nodeCharsEstimate is the assumed number of input bytes per node.
const nodeCharsEstimate = 8

// regexMap holds compiled regular expressions keyed by pattern.
var regexMap sync.Map

// parse parses input into an expression tree.
func parse(input string) (expr, error) {
	if input == "" {
		return expr{}, newError(KindParse, token{}, "empty input")
	}
	if len(input) > MaxInput {
		return expr{}, newError(KindParse, token{}, "input too long: %d bytes exceeds limit %d", len(input), MaxInput)
	}
	p := newParser(input)
	n, err := p.parseExpr()
	if err != nil {
		return expr{}, err
	}
	if t := p.peek(); t.typ != tokenEOF {
		if t.typ == tokenError {
			_, err := p.next()
			return expr{}, err
		}
		return expr{}, newError(KindParse, t, "unexpected token after parsing: %s", t.v)
	}
	nodes := p.nodes
	if nodes == nil {
		nodes = make([]node, p.nnode)
		copy(nodes, p.nodeBuf[:p.nnode])
	}
	return expr{
		nodes:  nodes,
		nident: int(p.nident),
		root:   n,
		shared: p.shared,
	}, nil
}

// parser builds the expression tree for one input. Nodes and identifiers
// are written into inline buffers by index so that the parser stays on the
// stack of parse.
type parser struct {
	lexer      lexer // lexer for tokenizing input
	current    token // current token
	parenCount int   // number of opening parentheses
	inputLen   int32 // input length in bytes for node capacity estimates
	peeked     bool  // indicates if the next token has been peeked

	nodeBuf [nodeBufSize]node // expression tree nodes until nodeBuf is full
	nodes   []node            // all expression tree nodes once nodeBuf overflowed
	nnode   int32             // number of nodes

	identBuf [identBufSize]string // distinct identifiers until identBuf is full
	idents   []string             // all distinct identifiers once identBuf overflowed
	nident   int32                // number of distinct identifiers
	shared   bool                 // some identifier is referenced more than once
}

// newParser creates a new parser for the input string.
func newParser(input string) parser {
	//nolint:gosec // Parse bounds the input length by MaxInput.
	return parser{
		lexer:    newLexer(input),
		inputLen: int32(len(input)),
	}
}

// parseExpr parses an expression.
func (p *parser) parseExpr() (int32, error) {
	return p.parseLogicalOr()
}

// parseLogicalOr parses OR expressions, the lowest precedence level.
func (p *parser) parseLogicalOr() (int32, error) {
	left, err := p.parseLogicalAnd()
	if err != nil {
		return 0, err
	}
	for p.peek().typ == tokenOR {
		t, err := p.next()
		if err != nil {
			return 0, err
		}
		right, err := p.parseLogicalAnd()
		if err != nil {
			return 0, err
		}
		left = p.addNode(newNodeBinary(left, t, right))
	}
	return left, nil
}

// parseLogicalAnd parses AND expressions, which bind tighter than OR.
func (p *parser) parseLogicalAnd() (int32, error) {
	left, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	for p.peek().typ == tokenAND {
		t, err := p.next()
		if err != nil {
			return 0, err
		}
		right, err := p.parseUnary()
		if err != nil {
			return 0, err
		}
		left = p.addNode(newNodeBinary(left, t, right))
	}
	return left, nil
}

// parseUnary parses an optional NOT prefix followed by a primary expression.
func (p *parser) parseUnary() (int32, error) {
	if p.peek().typ != tokenNOT {
		return p.parsePrimary()
	}
	t, err := p.next()
	if err != nil {
		return 0, err
	}
	child, err := p.parsePrimary()
	if err != nil {
		return 0, err
	}
	return p.addNode(newNodeUnary(child, t)), nil
}

// parsePrimary parses a parenthesized expression or a predicate.
func (p *parser) parsePrimary() (int32, error) {
	t := p.peek()
	switch t.typ {
	case tokenError:
		_, err := p.next()
		return 0, err
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
		return p.parsePredicate()
	default:
		return 0, newError(KindParse, t, "expected left parenthesis or identifier, got %s: %q", t.typ, t.v)
	}
}

// parsePredicate parses an identifier, a predicate operator, and a literal.
func (p *parser) parsePredicate() (int32, error) {
	ident, err := p.expect(tokenIdent)
	if err != nil {
		return 0, err
	}
	identIdx := p.identIndex(ident.v)
	op, err := p.next()
	if err != nil {
		return 0, err
	}
	if !op.typ.isPredicateOperatorType() {
		return 0, newError(KindParse, op, "expected predicate operator, got %s: %q", op.typ, op.v)
	}
	val, err := p.next()
	if err != nil {
		return 0, err
	}
	if !val.typ.isValueType() {
		return 0, newError(KindParse, val, "expected value, got %s: %q", val.typ, val.v)
	}
	if op.typ.isRegexOperatorType() && !val.typ.isStringType() {
		return 0, newError(KindParse, val, "expected string pattern, got %s: %q", val.typ, val.v)
	}
	switch val.typ {
	case tokenString, tokenRawString:
		val.v = unquote(val)
	case tokenBool:
		if val.v[0] == 't' || val.v[0] == 'T' {
			val.v = "true"
		} else {
			val.v = "false"
		}
	}
	i := p.addNode(newNodePredicate(ident, op, val, identIdx))
	if op.typ.isRegexOperatorType() {
		if err := p.cacheRegex(i, val); err != nil {
			return 0, err
		}
	}
	switch val.typ {
	case tokenString, tokenRawString:
		p.cacheValues(i, val.v)
	case tokenNumber:
		if !p.cacheNumber(i, val.v) {
			return 0, newError(KindParse, val, "invalid number %q", val.v)
		}
		p.cacheTime(i, val.v)
	case tokenTime:
		if !p.cacheTime(i, val.v) {
			return 0, newError(KindParse, val, "invalid time %q", val.v)
		}
	case tokenDuration:
		if !p.cacheDuration(i, val.v) {
			return 0, newError(KindParse, val, "invalid duration %q", val.v)
		}
	case tokenAddr:
		if !p.cacheAddr(i, val.v) {
			return 0, newError(KindParse, val, "invalid address %q", val.v)
		}
	}
	return i, nil
}

// cacheRegex compiles the pattern in t through regexMap and stores it on node i.
func (p *parser) cacheRegex(i int32, t token) error {
	if t.v == "" {
		return newError(KindParse, t, "invalid regex %q: empty pattern", t.v)
	}
	if cached, ok := regexMap.Load(t.v); ok {
		p.node(i).re = cached.(*regexp.Regexp)
		return nil
	}
	re, err := regexp.Compile(t.v)
	if err != nil {
		return newError(KindParse, t, "invalid regex %q: %w", t.v, err)
	}
	regexMap.Store(t.v, re)
	p.node(i).re = re
	return nil
}

// cacheValues stores on node i every number, time, duration, or address that the
// string literal s also spells.
func (p *parser) cacheValues(i int32, s string) {
	if strings.ContainsAny(s, ".:") && p.cacheAddr(i, s) {
		return
	}
	if len(s) == 0 || strings.IndexByte("0123456789+.-", s[0]) < 0 {
		if strings.Contains(s, ", ") {
			// A time whose layout starts with a weekday name.
			p.cacheTime(i, s)
		}
		return
	}
	l := newLexer(s)
	t := l.nextToken()
	if l.nextToken().typ != tokenEOF {
		// A time whose layout contains spaces.
		p.cacheTime(i, s)
		return
	}
	switch t.typ {
	case tokenNumber:
		if !strings.ContainsAny(s, "0123456789") {
			return
		}
		p.cacheNumber(i, s)
		p.cacheTime(i, s)
	case tokenTime:
		p.cacheTime(i, s)
	case tokenDuration:
		p.cacheDuration(i, s)
	}
}

// cacheNumber stores the number that s spells on node i and reports whether it did.
func (p *parser) cacheNumber(i int32, s string) bool {
	if v, err := parseNumber[int64](s); err == nil {
		p.node(i).valInt = v
		p.node(i).hasInt = true
		return true
	}
	if v, err := parseNumber[uint64](s); err == nil {
		p.node(i).valUint = v
		p.node(i).hasUint = true
		return true
	}
	v, err := parseNumber[float64](s)
	if err != nil {
		return false
	}
	p.node(i).valFloat = v
	p.node(i).hasFloat = true
	return true
}

// cacheTime stores the time that s spells on node i and reports whether it did.
func (p *parser) cacheTime(i int32, s string) bool {
	t, err := parseTime(s)
	if err != nil {
		return false
	}
	p.node(i).valTime = t
	p.node(i).hasTime = true
	return true
}

// cacheDuration stores the duration that s spells on node i and reports whether it did.
func (p *parser) cacheDuration(i int32, s string) bool {
	d, err := time.ParseDuration(s)
	if err != nil {
		return false
	}
	p.node(i).valDuration = d
	p.node(i).hasDuration = true
	return true
}

// cacheAddr stores the IP address that s spells on node i and reports whether it did.
func (p *parser) cacheAddr(i int32, s string) bool {
	v, err := netip.ParseAddr(s)
	if err != nil {
		return false
	}
	p.node(i).valAddr = v
	p.node(i).hasAddr = true
	return true
}

// identIndex returns the index of the identifier, registering it on first use.
func (p *parser) identIndex(name string) int32 {
	idents := p.idents
	if idents == nil {
		idents = p.identBuf[:p.nident]
	}
	for i := range p.nident {
		if idents[i] == name {
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
func (p *parser) addNode(n node) int32 {
	i := p.nnode
	switch {
	case p.nodes != nil:
		p.nodes = append(p.nodes, n)
	case i < nodeBufSize:
		p.nodeBuf[i] = n
	default:
		// current retains the original token text, including quotes.
		remaining := int(p.inputLen-p.current.pos) - len(p.current.v)
		p.nodes = make([]node, i, max(2*nodeBufSize, int(i)+remaining/nodeCharsEstimate))
		copy(p.nodes, p.nodeBuf[:])
		p.nodes = append(p.nodes, n)
	}
	p.nnode++
	return i
}

// node returns the node at index i.
func (p *parser) node(i int32) *node {
	if p.nodes != nil {
		return &p.nodes[i]
	}
	return &p.nodeBuf[i]
}

// expect consumes and returns the next token, reporting an error if its type differs.
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
	} else {
		p.current = p.lexer.nextToken()
	}
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
