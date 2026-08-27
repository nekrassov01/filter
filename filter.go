package filter

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// fieldCacheSize is the number of field values cached on the stack per evaluation.
// Expressions with more distinct identifiers fall back to a heap-allocated cache.
const fieldCacheSize = 8

// Target implements the entity to be evaluated.
type Target interface {
	Value(key string) (any, error)
}

// Expr represents a parsed expression.
type Expr struct {
	nodes  []node // expression tree nodes
	root   int    // index of the root node
	nident int    // number of distinct identifiers
	shared bool   // some identifier is referenced more than once
}

// Eval evaluates the expression against a target.
func (e *Expr) Eval(t Target) (bool, error) {
	if !e.shared {
		return eval(e.nodes, e.root, t, nil)
	}
	var buf [fieldCacheSize]field
	cache := buf[:]
	if e.nident > fieldCacheSize {
		cache = make([]field, e.nident)
	}
	return eval(e.nodes, e.root, t, cache)
}

// field holds a fetched field value for reuse within one evaluation.
type field struct {
	v  any
	ok bool
}

func eval(nodes []node, i int, t Target, cache []field) (bool, error) {
	n := &nodes[i]
	switch n.typ {
	case nodeBinary:
		switch n.op.typ {
		case tokenAND:
			left, err := eval(nodes, n.left, t, cache)
			if err != nil {
				return false, err
			}
			if !left {
				return false, nil
			}
			return eval(nodes, n.right, t, cache)
		case tokenOR:
			left, err := eval(nodes, n.left, t, cache)
			if err != nil {
				return false, err
			}
			if left {
				return true, nil
			}
			return eval(nodes, n.right, t, cache)
		default:
			return false, &Error{
				Kind: KindEval,
				Err:  fmt.Errorf("invalid logical operator at %d:%d: %q", n.op.line, n.op.col, n.op.typ.literal()),
			}
		}
	case nodeNOT:
		v, err := eval(nodes, n.left, t, cache)
		if err != nil {
			return false, err
		}
		return !v, nil
	case nodeComparison:
		if cache != nil && cache[n.ident.idx].ok {
			return evalComparison(n, cache[n.ident.idx].v)
		}
		v, err := t.Value(n.ident.v)
		if err != nil {
			return false, &Error{
				Kind: KindEval,
				Err:  err,
			}
		}
		if cache != nil {
			cache[n.ident.idx] = field{v: v, ok: true}
		}
		return evalComparison(n, v)
	}
	return false, &Error{
		Kind: KindEval,
		Err:  fmt.Errorf("invalid node type at %d:%d: %q", n.op.line, n.op.col, n.op.typ),
	}
}

// evalComparison evaluates a comparison expression against a target field.
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

// evalString evaluates a string expression against a target.
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
		return false, &Error{
			Kind: KindEval,
			Err:  fmt.Errorf("invalid operator for string field at %d:%d: %q", n.op.line, n.op.col, n.op.typ.literal()),
		}
	}
}

// evalNumber evaluates a number expression against a target.
func evalNumber(n *node, v float64) (bool, error) {
	f := n.num
	if !n.hasNum {
		parsed, err := strconv.ParseFloat(n.val.v, 64)
		if err != nil {
			return false, &Error{
				Kind: KindEval,
				Err:  fmt.Errorf("invalid number at %d:%d: %q", n.val.line, n.val.col, n.val.v),
			}
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
		return false, &Error{
			Kind: KindEval,
			Err:  fmt.Errorf("invalid operator for number field at %d:%d: %q", n.op.line, n.op.col, n.op.typ.literal()),
		}
	}
}

// evalTime evaluates a time expression against a target.
func evalTime(n *node, v time.Time) (bool, error) {
	t := n.time
	if !n.hasTime {
		parsed, err := time.Parse(time.RFC3339, n.val.v)
		if err != nil {
			return false, &Error{
				Kind: KindEval,
				Err:  fmt.Errorf("invalid time at %d:%d: %q", n.val.line, n.val.col, n.val.v),
			}
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
		return false, &Error{
			Kind: KindEval,
			Err:  fmt.Errorf("invalid operator for time field at %d:%d: %q", n.op.line, n.op.col, n.op.typ.literal()),
		}
	}
}

// evalDuration evaluates a duration expression against a target.
func evalDuration(n *node, v time.Duration) (bool, error) {
	d := n.dur
	if !n.hasDur {
		parsed, err := time.ParseDuration(n.val.v)
		if err != nil {
			return false, &Error{
				Kind: KindEval,
				Err:  fmt.Errorf("invalid duration at %d:%d: %q", n.val.line, n.val.col, n.val.v),
			}
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
		return false, &Error{
			Kind: KindEval,
			Err:  fmt.Errorf("invalid operator for duration field at %d:%d: %q", n.op.line, n.op.col, n.op.typ.literal()),
		}
	}
}
