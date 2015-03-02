package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// rational stores normalized rational numbers. Integers are stored as {n, 1}
type rational struct {
	n, d int64
}

func newRational(s string) rational {
	p := strings.Split(s, "/")
	if len(p) > 2 {
		log.Fatalf("Cannot convert %s to rational\n", s)
	}
	num, err := strconv.Atoi(p[0])
	if err != nil {
		log.Fatalf("Cannot convert %s to number: %s\n", p[0], err)
	}
	var denom int
	if len(p) == 2 {
		denom, err = strconv.Atoi(p[1])
		if err != nil {
			log.Fatalf("Cannot convert %s to number: %s\n", p[1], err)
		}
	} else {
		denom = 1
	}
	return rational{n: int64(num), d: int64(denom)}
}

func (r rational) String() string {
	if r.d == 1 {
		return strconv.FormatInt(r.n, 10)
	} else {
		return fmt.Sprintf("%s/%s", strconv.FormatInt(r.n, 10), strconv.FormatInt(r.d, 10))
	}
}

// Normalize return {d / gcd(n, d), n / gcd(n, d)}, making sure that denominator positive.
func (r rational) Normalize() rational {
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

func gcd(a, b int64) int64 {
	if a <= 0 && b <= 0 {
		return MaxInt64
	} else if a == 0 {
		return b
	} else if b == 0 {
		return a
	}
	return gcd(b, a%b)
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
		}.Normalize()
	}
}

func (r rational) Sub(r1 rational) rational {
	return r.Add(rational{n: -r1.n, d: r1.d})
}

func (r rational) Mul(r1 rational) rational {
	return rational{
		n: r.n * r1.n,
		d: r.d * r1.d,
	}.Normalize()
}

func (r rational) Div(r1 rational) rational {
	return rational{
		n: r.n * r1.d,
		d: r.d * r1.n,
	}.Normalize()
}

func (r rational) Pow(r1 rational) (rational, error) {
	r1 = r1.Normalize()
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
		return rational{n1, d1}.Normalize(), nil
	} else {
		n2, d2 := root(n1, r1.d), root(d1, r1.d)
		if n2 != MaxInt64 && d2 != MaxInt64 {
			return rational{n: n2, d: d2}.Normalize(), nil
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

func (r rational) Less(r1 rational) bool {
	x := r.n*r1.d < r.d*r1.n
	if r.d*r1.d > 0 {
		return x
	} else {
		return !x
	}
}

func (r rational) Negative() bool {
	r = r.Normalize()
	return r.n < 0
}

func (r rational) Integer() bool {
	r = r.Normalize()
	return r.d == 1
}

func (r rational) Minus() rational {
	return rational{-r.n, r.d}.Normalize()
}

func (r rational) Value() float64 {
	return float64(r.n) / float64(r.d)
}

func (r rational) Equal(r1 rational) bool {
	r = r.Normalize()
	r1 = r1.Normalize()
	return r.n == r1.n && r.d == r1.d
}

func (r rational) Even() bool {
	r = r.Normalize()
	return r.d == 1 && r.n%2 == 0
}

func (r rational) Zero() bool {
	return r.n == 0
}

func (r rational) Sqrt() (rational, error) {
	return r.Pow(rational{1, 2})
}
