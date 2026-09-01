package filter

import (
	"regexp"
	"time"
)

// nodeType represents the type of a node in the expression tree.
type nodeType uint8

// Node types of the expression tree.
const (
	nodeBinary    nodeType = iota // binary operator node type
	nodeUnary                     // unary NOT node type
	nodePredicate                 // predicate node type
)

// String returns a human-readable name for the node type.
func (t nodeType) String() string {
	switch t {
	case nodeBinary:
		return "binary node"
	case nodeUnary:
		return "unary node"
	case nodePredicate:
		return "predicate node"
	}
	return ""
}

// node represents a node in the expression tree.
type node struct {
	// Node metadata
	typ   nodeType       // type of the node
	left  int32          // left child index
	right int32          // right child index
	ident token          // identifier token for variable nodes
	op    token          // operator token for binary and predicate nodes
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
func newNodeBinary(left int32, op token, right int32) node {
	return node{
		typ:   nodeBinary,
		left:  left,
		right: right,
		op:    op,
	}
}

// newNodeUnary creates a new unary NOT node.
func newNodeUnary(child int32, op token) node {
	return node{
		typ:  nodeUnary,
		left: child,
		op:   op,
	}
}

// newNodePredicate creates a new predicate node.
func newNodePredicate(ident token, op token, val token) node {
	return node{
		typ:   nodePredicate,
		ident: ident,
		op:    op,
		val:   val,
	}
}
