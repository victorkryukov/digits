package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	a  string
	op Op
	b  string
	r  string
}

func TestRationalOps(t *testing.T) {
	assert := assert.New(t)
	cases := []testCase{
		{"1/2", OpAdd, "1/2", "1"},
		{"1/3", OpAdd, "-1/2", "-1/6"},
		{"2/4", OpAdd, "1/3", "5/6"},
		{"3/1", OpAdd, "-4/-1", "7"},
		{"3/1", OpAdd, "4/1", "7"},
		{"0/1", OpAdd, "-4/-6", "2/3"},

		{"1/2", OpSub, "1/2", "0"},
		{"1/2", OpSub, "1/3", "1/6"},

		{"1/2", OpMul, "2/3", "1/3"},
		{"2/5", OpMul, "120", "48"},

		{"1/2", OpDiv, "2/1", "1/4"},
		{"3/5", OpDiv, "5/3", "9/25"},
		{"3/5", OpDiv, "3/5", "1"},

		{"1/2", OpPow, "2", "1/4"},
		{"2/3", OpPow, "3", "8/27"},
		{"0", OpPow, "1", "0"},
		{"1", OpPow, "0", "1"},
		{"-1", OpPow, "0", "1"},
		{"1/2", OpPow, "-1", "2"},
		{"1/4", OpPow, "-1/2", "2"},
		{"4/9", OpPow, "3/2", "8/27"},
		{"4/9", OpPow, "3/-2", "27/8"},
		{"8/27", OpPow, "2/3", "4/9"},
		{"8/27", OpPow, "-2/3", "9/4"},
		{"8/-27", OpPow, "1/3", "-2/3"},
		{"-1", OpPow, "3", "-1"},
		{"-1", OpPow, "0", "1"},
		{"-1", OpPow, "-1/3", "-1"},
		{"-1", OpPow, "1/3", "-1"},

		{"0", OpFact, "", "1"},
		{"1", OpFact, "", "1"},
		{"2", OpFact, "", "2"},
		{"3", OpFact, "", "6"},
		{"4", OpFact, "", "24"},
		{"5", OpFact, "", "120"},

		{"9", OpSqrt, "", "3"},
		{"9/4", OpSqrt, "", "3/2"},
		{"1", OpSqrt, "", "1"},
		{"1000002000001", OpSqrt, "", "1000001"},
	}
	for _, tc := range cases {
		r1, _ := newRationalFromString(tc.a)
		r2, _ := newRationalFromString(tc.b)
		r, _ := newRationalFromString(tc.r)
		var v Value
		var err error
		if tc.op.binary() {
			v, err = r1.PerformBinary(tc.op, r2)
		} else {
			v, err = r1.PerformUnary(tc.op)
		}
		assert.NoError(err)
		assert.True(r.Equal(v))
	}
}

func rat(s string) rational {
	r, _ := newRationalFromString(s)
	return r
}

func TestRationalErrors(t *testing.T) {
	assert := assert.New(t)
	cases := []struct{ a, op, b string }{
		{"1/2", "!", ""},
		{"-1", "!", ""},
		{"21", "!", ""},

		{"0", "^", "0"},
		{"0", "^", "-1"},
		{"1/2", "^", "1/2"},
		{"20/3", "^", "20"},
		{"-8", "^", "1/4"},
		{"30", "^", "14"},
		{"1/4", "^", "1/3"},
		{"16", "^", "16"},
		{"1/16", "^", "16"},
		{"1/30", "^", "14"},

		{"1/3", "sqrt", ""},
		{"-5", "sqrt", ""},
		{"5", "sqrt", ""},
		{"1000002000002", "sqrt", ""},
	}
	for _, tc := range cases {
		switch tc.op {
		case "!":
			_, err := rat(tc.a).Fact()
			assert.Error(err, "%s!", rat(tc.a))
		case "^":
			_, err := rat(tc.a).Pow(rat(tc.b))
			assert.Error(err, "%s ^ %s", rat(tc.a), rat(tc.b))
		case "sqrt":
			_, err := rat(tc.a).Sqrt()
			assert.Error(err, "sqrt(%s)", rat(tc.a))
		}
	}
}

func TestRationalMisc(t *testing.T) {
	assert := assert.New(t)

	assert.True(rat("-1/2").Negative())
	assert.True(rat("1/-2").Negative())
	assert.False(rat("0").Negative())
	assert.False(rat("2/1").Negative())

	assert.True(rat("1/3").Less(rat("1/2")))
	assert.False(rat("1/3").Less(rat("1/-2")))
	assert.False(rat("-1/3").Less(rat("1/-2")))
	assert.False(rat("1/3").Less(rat("-2")))

	assert.True(rat("9").IsInteger())
	assert.True(rat("-9/3").IsInteger())
	assert.False(rat("1/2").IsInteger())

	assert.True(rat("0/100").Zero())
	assert.False(rat("1/200").Zero())

	assert.True(rat("-2").Even())
	assert.True(rat("2").Even())
	assert.True(rat("0").Even())
	assert.False(rat("-1").Even())
	assert.False(rat("1").Even())
	assert.False(rat("4/8").Even())
	assert.False(rat("-4/8").Even())

	assert.True(rat("3/9").Equal(rat("1/3")))
	assert.True(rat("-3/6").Equal(rat("1/-2")))
	assert.False(rat("-3/6").Equal(rat("-2")))
}
