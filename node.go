package filter

import (
	"regexp"
	"time"
)

// nodeType represents the type of a node in the expression tree.
type nodeType int

const (
	nodeBinary     nodeType = iota // binary operator node type
	nodeNOT                        // logical NOT node type
	nodeComparison                 // comparison node type
)

// String returns a string representation of the node type.
func (t nodeType) String() string {
	switch t {
	case nodeBinary:
		return "binary node"
	case nodeNOT:
		return "not node"
	case nodeComparison:
		return "comparison node"
	}
	return ""
}

// node represents a node in the expression tree.
type node struct {
	// Node metadata
	typ   nodeType       // type of the node
	left  int            // left child index
	right int            // right child index
	ident token          // identifier token for variable nodes
	op    token          // operator token for binary and comparison nodes
	val   token          // value token for literal nodes
	re    *regexp.Regexp // regular expression for pattern matching

	// Cached values
	num  float64       // cached numeric value
	dur  time.Duration // cached duration value
	time time.Time     // cached time value

	// Cached flags
	hasNum  bool // indicates if num is cached
	hasDur  bool // indicates if dur is cached
	hasTime bool // indicates if time is cached
}

// newNodeBinary creates a new binary expression node.
func newNodeBinary(p *parser, left int, op token, right int) int {
	node := node{
		typ:   nodeBinary,
		left:  left,
		right: right,
		op:    op,
	}
	p.nodes = append(p.nodes, node)
	return len(p.nodes) - 1
}

// newNodeNOT creates a new NOT expression node.
func newNodeNOT(p *parser, child int, op token) int {
	node := node{
		typ:  nodeNOT,
		left: child,
		op:   op,
	}
	p.nodes = append(p.nodes, node)
	return len(p.nodes) - 1
}

// newNodeComparison creates a new comparison expression node.
func newNodeComparison(p *parser, ident token, op token, val token) int {
	node := node{
		typ:   nodeComparison,
		ident: ident,
		op:    op,
		val:   val,
	}
	p.nodes = append(p.nodes, node)
	return len(p.nodes) - 1
}
