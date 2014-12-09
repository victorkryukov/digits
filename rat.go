package main

import (
	"fmt"
	"math"
	"strconv"
)

// Lookup tables for fast factorial and integer square root
var factLookup, sqrtLookup map[int64]int64

const (
	maxFactorial = 20      // Pre-calculate n! up to this
	maxSqrt      = 1000000 // Pre-calculate sqrt[n] up to this
	maxSqrt2     = maxSqrt * maxSqrt
)

func init() {
	factLookup = make(map[int64]int64)
	sqrtLookup = make(map[int64]int64)
	factLookup[0] = 1
	var i, fact int64
	fact = 2
	// We are starting from 3, since 1! = 1 and 2! = 2
	for i = 3; i <= maxFactorial; i++ {
		fact *= i
		factLookup[i] = fact
	}
	for i = 1; i <= maxSqrt; i++ {
		sqrtLookup[i*i] = i
	}
}

// fact calculates n! using lookup table, and returns MaxInt64 for invalid inputs
func fact(n int64) int64 {
	if n < 0 || n > maxFactorial || (n < 3 && n != 0) {
		return MaxInt64
	}
	return factLookup[n]
}

// sqrt calculates sqrt[n] using lookup table, and returns MaxInt64 for invalid inputs
// or non-integer square roots
func sqrt(n int64) int64 {
	if n < 2 {
		return MaxInt64
	} else if n > maxSqrt2 {
		s := int64(math.Floor(math.Sqrt(float64(n))))
		if s*s == n {
			return s
		} else {
			return MaxInt64
		}
	}
	// n <= maxSqrt2, so use the lookup table
	if r, ok := sqrtLookup[n]; ok {
		return r
	}
	return MaxInt64
}

// root calculate a ^ (1/b) for positive b, and returns -1 for non-integer results or invalid inputs.
func root(a, b int64) int64 {
	if a == 0 && b == 0 {
		return MaxInt64
	} else if b == 1 || a == 0 || a == 1 {
		return a
	} else if b == 2 {
		return sqrt(a)
	}
	// a != 0, b > 2
	var isOdd int64 = 1
	if a < 0 && b%2 == 1 { // math.Pow doesn't support e.g. -8 ^ (1/3)
		isOdd = -1
		a = -a
	} else if a < 0 {
		return MaxInt64
	}
	r := int64(math.Floor(math.Pow(float64(a), 1/float64(b)) + 0.5))
	if pow(r, b) == a {
		return isOdd * r
	} else {
		return MaxInt64
	}
}

const MaxInt64 = 9223372036854775807

// pow returns a^b, if both of them are small enough, or MaxInt64 if the result is invalid
// FIXME: Once we support ratios, we should support a^r where r is a ratio, too.
func pow(a, b int64) int64 {
	if a == 0 && b <= 0 {
		return MaxInt64
	} else if a == 0 || b == 1 {
		return a
	} else if b == 0 {
		return 1
	}
	if b > 15 && (a > 15 || a < -15) { // 16 ^ 16 > MaxInt64
		return MaxInt64
	}
	if p := math.Pow(float64(a), float64(b)); math.Abs(p) > MaxInt64 {
		return MaxInt64
	} else {
		return int64(p)
	}
}

// Rat stores normalized rational numbers. Integers are stored as {n, 1}
type Rat struct {
	n, d int64
}

var BadRat Rat

func init() {
	BadRat = Rat{0, 0}
}

func (r Rat) String() string {
	if r.d == 1 {
		return strconv.FormatInt(r.n, 10)
	} else {
		return fmt.Sprintf("%s/%s", strconv.FormatInt(r.n, 10), strconv.FormatInt(r.d, 10))
	}
}

// Normalize return {d / gcd(n, d), n / gcd(n, d)}, making sure that denominator positive.
func (r Rat) Normalize() Rat {
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
	return Rat{n: n1 / g, d: d1 / g}
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
func (r Rat) Add(r1 Rat) Rat {
	if r.d == 1 && r1.d == 1 {
		return Rat{
			n: r.n + r1.n,
			d: 1,
		}
	} else {
		return Rat{
			n: r.n*r1.d + r1.n*r.d,
			d: r.d * r1.d,
		}.Normalize()
	}
}

func (r Rat) Sub(r1 Rat) Rat {
	return r.Add(Rat{n: -r1.n, d: r1.d})
}

func (r Rat) Mul(r1 Rat) Rat {
	return Rat{
		n: r.n * r1.n,
		d: r.d * r1.d,
	}.Normalize()
}

func (r Rat) Div(r1 Rat) Rat {
	return Rat{
		n: r.n * r1.d,
		d: r.d * r1.n,
	}.Normalize()
}

func (r Rat) Pow(r1 Rat) (Rat, error) {
	r1 = r1.Normalize()
	if r1.n < 0 {
		if r.n == 0 {
			return BadRat, fmt.Errorf("Cannot raise 0 to %d", r1.n)
		}
		return Rat{n: r.d, d: r.n}.Pow(r1.Minus())
	}
	n1 := pow(r.n, r1.n)
	if n1 == MaxInt64 {
		return BadRat, fmt.Errorf("Cannot calculate %d^%d", r.n, r1.n)
	}
	d1 := pow(r.d, r1.n)
	if d1 == MaxInt64 {
		return BadRat, fmt.Errorf("Cannot calculate %d^%d", r.d, r1.n)
	}
	if r1.d == 1 {
		return Rat{n1, d1}.Normalize(), nil
	} else {
		n2, d2 := root(n1, r1.d), root(d1, r1.d)
		if n2 != MaxInt64 && d2 != MaxInt64 {
			return Rat{n: n2, d: d2}.Normalize(), nil
		} else {
			return BadRat, fmt.Errorf("Cannot calculate root[%d] of %s (2), n2 = %d, d2 = %d", r1.d, Rat{n: n1, d: d1}, n2, d2)
		}
	}
}

func (r Rat) Fact() (Rat, error) {
	if r.d != 1 {
		return BadRat, fmt.Errorf("Cannot calculate %s!", r)
	}
	if r.n == 1 || r.n == 2 {
		return r, nil
	}
	if f := fact(r.n); f == MaxInt64 {
		return BadRat, fmt.Errorf("Cannot calculate %d!", r)
	} else {
		return Rat{f, 1}, nil
	}
}

func (r Rat) Less(r1 Rat) bool {
	x := r.n*r1.d < r.d*r1.n
	if r.d*r1.d > 0 {
		return x
	} else {
		return !x
	}
}

func (r Rat) Negative() bool {
	r = r.Normalize()
	return r.n < 0
}

func (r Rat) Integer() bool {
	r = r.Normalize()
	return r.d == 1
}

func (r Rat) Minus() Rat {
	return Rat{-r.n, r.d}.Normalize()
}

func (r Rat) Value() float64 {
	return float64(r.n) / float64(r.d)
}

func (r Rat) Equal(r1 Rat) bool {
	r = r.Normalize()
	r1 = r1.Normalize()
	return r.n == r1.n && r.d == r1.d
}

func (r Rat) Even() bool {
	r = r.Normalize()
	return r.d == 1 && r.n%2 == 0
}

func (r Rat) Zero() bool {
	return r.n == 0
}

func (r Rat) Sqrt() (Rat, error) {
	return r.Pow(Rat{1, 2})
}
