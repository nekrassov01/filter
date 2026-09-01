package filter

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
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

// timeLayouts are the layouts a time literal may use, most common first.
var timeLayouts = [...]string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	time.DateTime,
	time.DateOnly,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC850,
	time.RFC822,
	time.RFC822Z,
}

// regexMap holds compiled regular expressions keyed by pattern.
var regexMap sync.Map

// parse parses input into an expression tree.
func parse(input string) (expr, error) {
	if input == "" {
		return expr{}, &Error{
			Kind: KindParse,
			Err:  fmt.Errorf("empty input"),
		}
	}
	if len(input) > MaxInput {
		return expr{}, &Error{
			Kind: KindParse,
			Err:  fmt.Errorf("input too long: %d bytes exceeds limit %d", len(input), MaxInput),
		}
	}
	p := newParser(input)
	n, err := p.parseExpr()
	if err != nil {
		return expr{}, err
	}
	if t := p.peek(); t.typ != tokenEOF {
		return expr{}, newError(KindParse, t, "unexpected token after parsing: %s", t.v)
	}
	nodes := p.nodes
	if nodes == nil {
		nodes = make([]node, p.nnode)
		copy(nodes, p.nodeBuf[:p.nnode])
	}
	return expr{
		nodes:  nodes,
		root:   n,
		nident: int(p.nident),
		shared: p.shared,
	}, nil
}

// parser builds the expression tree for one input. Nodes and identifiers
// are written into inline buffers by index so that the parser stays on the
// stack of parse.
type parser struct {
	lexer      lexer                // lexer for tokenizing input
	current    token                // current token
	peeked     bool                 // indicates if the next token has been peeked
	parenCount int                  // number of opening parentheses
	nodeBuf    [nodeBufSize]node    // expression tree nodes until nodeBuf is full
	nodes      []node               // all expression tree nodes once nodeBuf overflowed
	nnode      int32                // number of nodes
	identBuf   [identBufSize]string // distinct identifiers until identBuf is full
	idents     []string             // all distinct identifiers once identBuf overflowed
	nident     int32                // number of distinct identifiers
	shared     bool                 // some identifier is referenced more than once
}

// newParser creates a new parser for the input string.
func newParser(input string) parser {
	return parser{
		lexer: newLexer(input),
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
	for {
		if p.peek().typ == tokenOR {
			t, err := p.next()
			if err != nil {
				return 0, err
			}
			right, err := p.parseLogicalAnd()
			if err != nil {
				return 0, err
			}
			left = p.addNode(newNodeBinary(left, t, right))
			continue
		}
		break
	}
	return left, nil
}

// parseLogicalAnd parses AND expressions, which bind tighter than OR.
func (p *parser) parseLogicalAnd() (int32, error) {
	left, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	for {
		if p.peek().typ == tokenAND {
			t, err := p.next()
			if err != nil {
				return 0, err
			}
			right, err := p.parseUnary()
			if err != nil {
				return 0, err
			}
			left = p.addNode(newNodeBinary(left, t, right))
			continue
		}
		break
	}
	return left, nil
}

// parseUnary parses an optional NOT prefix followed by a primary expression.
func (p *parser) parseUnary() (int32, error) {
	if p.peek().typ == tokenNOT {
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
	return p.parsePrimary()
}

// parsePrimary parses a parenthesized expression or a predicate.
func (p *parser) parsePrimary() (int32, error) {
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
	ident.idx = p.identIndex(ident.v)
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
	i := p.addNode(newNodePredicate(ident, op, val))
	if op.typ.isRegexOperatorType() {
		if err := p.cacheRegex(i, val); err != nil {
			return 0, err
		}
	}
	switch val.typ {
	case tokenString, tokenRawString:
		p.cacheValues(i, val.v)
	case tokenTime:
		if !p.cacheTime(i, val.v) {
			return 0, newError(KindParse, val, "invalid time %q", val.v)
		}
	case tokenDuration:
		if !p.cacheDuration(i, val.v) {
			return 0, newError(KindParse, val, "invalid duration %q", val.v)
		}
	case tokenNumber:
		if !p.cacheNumber(i, val.v) {
			return 0, newError(KindParse, val, "invalid number %q", val.v)
		}
		p.cacheTime(i, val.v)
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

// cacheValues stores on node i every time, duration, or number that the
// string literal s also spells.
func (p *parser) cacheValues(i int32, s string) {
	r, _ := utf8.DecodeRuneInString(s)
	switch {
	case isNumberStart(r):
		l := newLexer(s)
		tok := l.nextToken()
		if l.nextToken().typ != tokenEOF {
			// A time whose layout contains spaces.
			p.cacheTime(i, s)
			return
		}
		switch tok.typ {
		case tokenTime:
			p.cacheTime(i, s)
		case tokenDuration:
			p.cacheDuration(i, s)
		case tokenNumber:
			if !strings.ContainsAny(s, "0123456789") {
				return
			}
			p.cacheNumber(i, s)
			p.cacheTime(i, s)
		}
	case strings.Contains(s, ", "):
		// A time whose layout starts with a weekday name.
		p.cacheTime(i, s)
	}
}

// cacheTime stores the time that s spells on node i and reports whether it did.
func (p *parser) cacheTime(i int32, s string) bool {
	t, err := parseTime(s)
	if err != nil {
		return false
	}
	p.node(i).time = t
	p.node(i).hasTime = true
	return true
}

// cacheDuration stores the duration that s spells on node i and reports whether it did.
func (p *parser) cacheDuration(i int32, s string) bool {
	d, err := time.ParseDuration(s)
	if err != nil {
		return false
	}
	p.node(i).dur = d
	p.node(i).hasDur = true
	return true
}

// cacheNumber stores the number that s spells on node i and reports whether it did.
func (p *parser) cacheNumber(i int32, s string) bool {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false
	}
	p.node(i).num = f
	p.node(i).hasNum = true
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
		remaining := len(p.lexer.input) - int(p.lexer.pos)
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

// parseTime converts Unix seconds or a literal in one of timeLayouts to a
// UTC time. A zone abbreviation other than UTC or GMT is rejected, since
// time.Parse resolves no other abbreviation to an offset.
func parseTime(s string) (time.Time, error) {
	// Unix seconds.
	digits := s
	if digits != "" && (digits[0] == '-' || digits[0] == '+') {
		digits = digits[1:]
	}
	if digits != "" && digits[0] != '_' && digits[len(digits)-1] != '_' {
		var sec int64
		integer := true
		for i := 0; i < len(digits) && integer; i++ {
			switch c := digits[i]; {
			case c == '_':
			case '0' <= c && c <= '9':
				d := int64(c - '0')
				if sec > (math.MaxInt64-d)/10 {
					return time.Time{}, fmt.Errorf("unix seconds out of range %q", s)
				}
				sec = sec*10 + d
			default:
				integer = false
			}
		}
		if integer {
			if s[0] == '-' {
				sec = -sec
			}
			return time.Unix(sec, 0).UTC(), nil
		}
	}
	for _, layout := range timeLayouts {
		t, err := time.ParseInLocation(layout, s, time.UTC)
		if err != nil {
			continue
		}
		if name, _ := t.Zone(); name != "" && name != "UTC" && name != "GMT" {
			return time.Time{}, fmt.Errorf("unknown time zone %q", name)
		}
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("unrecognized time %q", s)
}

// unquote returns the text of a string token without its surrounding quotes.
func unquote(t token) string {
	n := len(t.v)
	if t.typ.isStringType() && n >= 2 {
		return t.v[1 : n-1]
	}
	return t.v
}
