package main

import "math"

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
	if a < 0 && b%2 == 1 {
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
	if b > 15 && (a > 15 || a < -15) {
		return MaxInt64
	}
	if p := math.Pow(float64(a), float64(b)); math.Abs(p) > MaxInt64 {
		return MaxInt64
	} else {
		return int64(p)
	}
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
