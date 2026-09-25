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
	ident token          // identifier token for variable nodes
	op    token          // operator token for binary and predicate nodes
	val   token          // value token for literal nodes
	re    *regexp.Regexp // regular expression for pattern matching

	valInt      int64         // cached signed integer
	valUint     uint64        // cached unsigned integer
	valFloat    float64       // cached floating-point value
	valDuration time.Duration // cached duration value
	valTime     time.Time     // cached time value

	hasInt      bool // indicates if valInt is cached
	hasUint     bool // indicates if valUint is cached
	hasFloat    bool // indicates if valFloat is cached
	hasDuration bool // indicates if valDuration is cached
	hasTime     bool // indicates if valTime is cached

	typ   nodeType // type of the node
	left  int32    // left child index
	right int32    // right child index
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
