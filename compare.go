package filter

import (
	"cmp"
	"math"
)

// Epsilon is the tolerance within which two floating-point values compare as equal.
const Epsilon = 1e-9

// compareNumber returns ordering, equality, and whether the operands are ordered.
// Floating-point equality uses Epsilon; mixed comparisons preserve integer precision.
func compareNumber[L int64 | uint64 | float64, R int64 | uint64 | float64](left L, right R) (int, bool, bool) {
	if left != left || right != right {
		return 0, false, false
	}
	var c int
	switch left := any(left).(type) {
	case int64:
		switch right := any(right).(type) {
		case int64:
			c = cmp.Compare(left, right)
		case uint64:
			if left < 0 {
				c = -1
			} else {
				c = cmp.Compare(uint64(left), right)
			}
		case float64:
			c = compareIntFloat(left, right)
		}
	case uint64:
		switch right := any(right).(type) {
		case int64:
			if right < 0 {
				c = 1
			} else {
				c = cmp.Compare(left, uint64(right))
			}
		case uint64:
			c = cmp.Compare(left, right)
		case float64:
			c = compareUintFloat(left, right)
		}
	case float64:
		switch right := any(right).(type) {
		case int64:
			c = -compareIntFloat(right, left)
		case uint64:
			c = -compareUintFloat(right, left)
		case float64:
			return cmp.Compare(left, right), math.Abs(left-right) <= Epsilon, true
		}
	}
	return c, c == 0, true
}

// compareIntFloat returns -1, 0, or 1 when i is less than, equal to, or greater
// than f, without rounding i to float64. The caller must exclude NaN.
func compareIntFloat(i int64, f float64) int {
	// Check the exclusive upper bound before converting: float64(MaxInt64)
	// rounds up to 2^63.
	if f < -0x1p63 {
		return 1
	}
	if f >= 0x1p63 {
		return -1
	}
	truncated := int64(f)
	if c := cmp.Compare(i, truncated); c != 0 {
		return c
	}
	// The truncated part of an in-range float is exactly representable.
	return cmp.Compare(float64(truncated), f)
}

// compareUintFloat returns -1, 0, or 1 when u is less than, equal to, or greater
// than f, without rounding u to float64. The caller must exclude NaN.
func compareUintFloat(u uint64, f float64) int {
	// Check the exclusive upper bound before converting: float64(MaxUint64)
	// rounds up to 2^64.
	if f < 0 {
		return 1
	}
	if f >= 0x1p64 {
		return -1
	}
	truncated := uint64(f)
	if c := cmp.Compare(u, truncated); c != 0 {
		return c
	}
	// The truncated part of an in-range float is exactly representable.
	return cmp.Compare(float64(truncated), f)
}
