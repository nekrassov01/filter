package filter

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// cacheSize is the number of resolved values cached on the stack per evaluation.
// Expressions with more distinct identifiers fall back to a heap-allocated cache.
const cacheSize = 8

// Resolver resolves an identifier in an expression to its value.
// The second result reports whether the identifier is known; unknown
// identifiers are reported by Eval as an error with the position.
type Resolver interface {
	Resolve(name string) (any, bool)
}

// Expr represents a parsed expression.
type Expr struct {
	nodes  []node // expression tree nodes
	root   int32  // index of the root node
	nident int    // number of distinct identifiers
	shared bool   // some identifier is referenced more than once
}

// Eval evaluates the expression against the values provided by r.
func (e *Expr) Eval(r Resolver) (bool, error) {
	if !e.shared {
		return eval(e.nodes, e.root, r, nil)
	}
	var buf [cacheSize]cached
	cache := buf[:]
	if e.nident > cacheSize {
		cache = make([]cached, e.nident)
	}
	return eval(e.nodes, e.root, r, cache)
}

// cached holds a resolved value for reuse within one evaluation.
type cached struct {
	v  any
	ok bool
}

// eval evaluates the node at index i, resolving identifiers through r and
// reusing values from cache when it is non-nil.
func eval(nodes []node, i int32, r Resolver, cache []cached) (bool, error) {
	n := &nodes[i]
	switch n.typ {
	case nodeBinary:
		switch n.op.typ {
		case tokenAND:
			left, err := eval(nodes, n.left, r, cache)
			if err != nil {
				return false, err
			}
			if !left {
				return false, nil
			}
			return eval(nodes, n.right, r, cache)
		case tokenOR:
			left, err := eval(nodes, n.left, r, cache)
			if err != nil {
				return false, err
			}
			if left {
				return true, nil
			}
			return eval(nodes, n.right, r, cache)
		default:
			return false, newError(KindEval, n.op, "invalid logical operator %q", n.op.typ.literal())
		}
	case nodeNOT:
		v, err := eval(nodes, n.left, r, cache)
		if err != nil {
			return false, err
		}
		return !v, nil
	case nodeComparison:
		if cache != nil && cache[n.ident.idx].ok {
			return evalComparison(n, cache[n.ident.idx].v)
		}
		v, ok := r.Resolve(n.ident.v)
		if !ok {
			return false, newError(KindEval, n.ident, "unknown identifier %q", n.ident.v)
		}
		if cache != nil {
			cache[n.ident.idx] = cached{v: v, ok: true}
		}
		return evalComparison(n, v)
	}
	return false, newError(KindEval, n.op, "invalid node type %q", n.op.typ)
}

// evalComparison evaluates a comparison expression against a resolved value.
func evalComparison(n *node, v any) (bool, error) {
	switch v := v.(type) {
	case string:
		return evalString(n, v)
	case int:
		return evalNumber(n, float64(v))
	case int8:
		return evalNumber(n, float64(v))
	case int16:
		return evalNumber(n, float64(v))
	case int32:
		return evalNumber(n, float64(v))
	case int64:
		return evalNumber(n, float64(v))
	case uint:
		return evalNumber(n, float64(v))
	case uint8:
		return evalNumber(n, float64(v))
	case uint16:
		return evalNumber(n, float64(v))
	case uint32:
		return evalNumber(n, float64(v))
	case uint64:
		return evalNumber(n, float64(v))
	case float32:
		return evalNumber(n, float64(v))
	case float64:
		return evalNumber(n, v)
	case time.Time:
		return evalTime(n, v)
	case time.Duration:
		return evalDuration(n, v)
	default:
		return evalString(n, fmt.Sprint(v))
	}
}

// evalString evaluates a comparison against a string value.
func evalString(n *node, v string) (bool, error) {
	switch n.op.typ {
	case tokenEQ:
		return v == n.val.v, nil
	case tokenEQI:
		return strings.EqualFold(v, n.val.v), nil
	case tokenNEQ:
		return v != n.val.v, nil
	case tokenNEQI:
		return !strings.EqualFold(v, n.val.v), nil
	case tokenREQ, tokenREQI:
		return n.re.MatchString(v), nil
	case tokenNREQ, tokenNREQI:
		return !n.re.MatchString(v), nil
	default:
		return false, newError(KindEval, n.op, "invalid operator for string value %q", n.op.typ.literal())
	}
}

// evalNumber evaluates a comparison against a numeric value.
func evalNumber(n *node, v float64) (bool, error) {
	f := n.num
	if !n.hasNum {
		parsed, err := strconv.ParseFloat(n.val.v, 64)
		if err != nil {
			return false, newError(KindEval, n.val, "invalid number %q", n.val.v)
		}
		f = parsed
	}
	switch n.op.typ {
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
		return false, newError(KindEval, n.op, "invalid operator for number value %q", n.op.typ.literal())
	}
}

// evalTime evaluates a comparison against a time value.
func evalTime(n *node, v time.Time) (bool, error) {
	t := n.time
	if !n.hasTime {
		parsed, err := parseTime(n.val.v)
		if err != nil {
			return false, newError(KindEval, n.val, "invalid time %q", n.val.v)
		}
		t = parsed
	}
	switch n.op.typ {
	case tokenGT:
		return v.After(t), nil
	case tokenGTE:
		return v.Equal(t) || v.After(t), nil
	case tokenLT:
		return v.Before(t), nil
	case tokenLTE:
		return v.Equal(t) || v.Before(t), nil
	case tokenEQ:
		return v.Equal(t), nil
	case tokenNEQ:
		return !v.Equal(t), nil
	default:
		return false, newError(KindEval, n.op, "invalid operator for time value %q", n.op.typ.literal())
	}
}

// evalDuration evaluates a comparison against a duration value.
func evalDuration(n *node, v time.Duration) (bool, error) {
	d := n.dur
	if !n.hasDur {
		parsed, err := time.ParseDuration(n.val.v)
		if err != nil {
			return false, newError(KindEval, n.val, "invalid duration %q", n.val.v)
		}
		d = parsed
	}
	switch n.op.typ {
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
		return false, newError(KindEval, n.op, "invalid operator for duration value %q", n.op.typ.literal())
	}
}
