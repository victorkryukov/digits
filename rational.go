package main

import (
	"fmt"
	"strconv"
	"strings"
)

// rational stores normalized rational numbers. Integers are stored as {n, 1}
type rational struct {
	n, d int64
}

// newRational creates a normalized rational for a/b, and returns an error
// if b == 0.
func newRational(a, b int64) (rational, error) {
	if b == 0 {
		return rational{}, fmt.Errorf("%d/0 is not a proper rational", a)
	} else {
		return rational{a, b}.normalize(), nil
	}
}

// normalize return {d / gcd(n, d), n / gcd(n, d)}, making sure that denominator positive.
func (r rational) normalize() rational {
	n1, d1 := r.n, r.d
	if r.d < 0 {
		n1, d1 = -n1, -d1
	}
	var g int64
	if n1 < 0 {
		g = gcd(-n1, d1)
	} else {
		g = gcd(n1, d1)
	}
	return rational{n: n1 / g, d: d1 / g}
}

// newRationalFromString creates a normalized rational from a string "a/b", and
// returns an error if it cannot be parsed.
func newRationalFromString(s string) (rational, error) {
	p := strings.Split(s, "/")
	if len(p) > 2 {
		return rational{}, fmt.Errorf("cannot convert %s to rational\n", s)
	}
	num, err := strconv.Atoi(p[0])
	if err != nil {
		return rational{}, fmt.Errorf("cannot convert %s to number: %s\n", p[0], err)
	}
	var denom int
	if len(p) == 2 {
		denom, err = strconv.Atoi(p[1])
		if err != nil {
			return rational{}, fmt.Errorf("cannot convert %s to number: %s\n", p[1], err)
		}
	} else {
		denom = 1
	}
	return newRational(int64(num), int64(denom))
}

func (r rational) String() string {
	if r.d == 1 {
		return strconv.FormatInt(r.n, 10)
	} else {
		return fmt.Sprintf("%s/%s", strconv.FormatInt(r.n, 10), strconv.FormatInt(r.d, 10))
	}
}

// PerformUnary is an implementation of Value.PermormUnary
func (r rational) PerformUnary(op Op) (Value, error) {
	switch op {
	case OpFact:
		return r.Fact()
	case OpSqrt:
		return r.Sqrt()
	case OpMinus:
		return r.Minus(), nil
	default:
		return rational{}, fmt.Errorf("%s is not unary operator", op)
	}
}

// PerformBinary is an implementation of Value.PerformBinary
func (r rational) PerformBinary(op Op, v Value) (Value, error) {
	r1, ok := v.(rational)
	if !ok {
		return rational{}, fmt.Errorf("%v is not rational", v)
	}
	switch op {
	case OpAdd:
		return r.Add(r1), nil
	case OpSub:
		return r.Sub(r1), nil
	case OpMul:
		return r.Mul(r1), nil
	case OpDiv:
		return r.Div(r1)
	case OpPow:
		return r.Pow(r1)
	default:
		return rational{}, fmt.Errorf("%s is not binary operator", op)
	}
}

// Equal is an implementation of Value.Equal
func (r rational) Equal(v Value) bool {
	r1, ok := v.(rational)
	if !ok {
		return false
	}
	return r.isEqual(r1)
}

// Add returns r + r1
func (r rational) Add(r1 rational) rational {
	if r.d == 1 && r1.d == 1 {
		return rational{
			n: r.n + r1.n,
			d: 1,
		}
	} else {
		return rational{
			n: r.n*r1.d + r1.n*r.d,
			d: r.d * r1.d,
		}.normalize()
	}
}

func (r rational) Sub(r1 rational) rational {
	return r.Add(rational{n: -r1.n, d: r1.d})
}

func (r rational) Mul(r1 rational) rational {
	return rational{
		n: r.n * r1.n,
		d: r.d * r1.d,
	}.normalize()
}

func (r rational) Div(r1 rational) (rational, error) {
	if r1.n == 0 {
		return rational{}, fmt.Errorf("division by 0: %s / %s", r, r1)
	}
	return rational{
		n: r.n * r1.d,
		d: r.d * r1.n,
	}.normalize(), nil
}

func (r rational) Pow(r1 rational) (rational, error) {
	r1 = r1.normalize()
	if r1.n < 0 {
		if r.n == 0 {
			return rational{}, fmt.Errorf("Cannot raise 0 to %d", r1.n)
		}
		return rational{n: r.d, d: r.n}.Pow(r1.Minus())
	}
	n1 := pow(r.n, r1.n)
	if n1 == MaxInt64 {
		return rational{}, fmt.Errorf("Cannot calculate %d^%d", r.n, r1.n)
	}
	d1 := pow(r.d, r1.n)
	if d1 == MaxInt64 {
		return rational{}, fmt.Errorf("Cannot calculate %d^%d", r.d, r1.n)
	}
	if r1.d == 1 {
		return rational{n1, d1}.normalize(), nil
	} else {
		n2, d2 := root(n1, r1.d), root(d1, r1.d)
		if n2 != MaxInt64 && d2 != MaxInt64 {
			return rational{n: n2, d: d2}.normalize(), nil
		} else {
			return rational{}, fmt.Errorf("Cannot calculate root[%d] of %s (2), n2 = %d, d2 = %d", r1.d, rational{n: n1, d: d1}, n2, d2)
		}
	}
}

func (r rational) Fact() (rational, error) {
	if r.d != 1 {
		return rational{}, fmt.Errorf("Cannot calculate %s!", r)
	}
	if r.n == 1 || r.n == 2 {
		return r, nil
	}
	if f := fact(r.n); f == MaxInt64 {
		return rational{}, fmt.Errorf("Cannot calculate %d!", r)
	} else {
		return rational{f, 1}, nil
	}
}

func (r rational) isLess(r1 rational) bool {
	x := r.n*r1.d < r.d*r1.n
	if r.d*r1.d > 0 {
		return x
	} else {
		return !x
	}
}

func (r rational) Less(v Value) bool {
	r1, ok := v.(rational)
	if !ok {
		return false
	}
	return r.isLess(r1)
}

func (r rational) IsInteger() bool {
	r = r.normalize()
	return r.d == 1
}

func (r rational) Minus() rational {
	return rational{-r.n, r.d}.normalize()
}

func (r rational) Value() float64 {
	return float64(r.n) / float64(r.d)
}

func (r rational) isEqual(r1 rational) bool {
	r = r.normalize()
	r1 = r1.normalize()
	return r.n == r1.n && r.d == r1.d
}

func (r rational) Even() bool {
	r = r.normalize()
	return r.d == 1 && r.n%2 == 0
}

func (r rational) Zero() bool {
	return r.n == 0
}

func (r rational) One() bool {
	return r.isEqual(rational{1, 1})
}

func (r rational) MinusOne() bool {
	return r.isEqual(rational{-1, 1})
}

func (r rational) Sqrt() (rational, error) {
	return r.Pow(rational{1, 2})
}

func (r rational) Negative() bool {
	return (r.n < 0 && r.d > 0) || (r.n > 0 && r.d < 0)
}
