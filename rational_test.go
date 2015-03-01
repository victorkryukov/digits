package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	a, op, b, r string
}

func TestRationalOps(t *testing.T) {
	assert := assert.New(t)
	cases := []testCase{
		{"1/2", "+", "1/2", "1"},
		{"1/3", "+", "-1/2", "-1/6"},
		{"2/4", "+", "1/3", "5/6"},
		{"3/1", "+", "-4/-1", "7"},
		{"3/1", "+", "4/1", "7"},
		{"0/1", "+", "-4/-6", "2/3"},

		{"1/2", "-", "1/2", "0"},
		{"1/2", "-", "1/3", "1/6"},

		{"1/2", "*", "2/3", "1/3"},
		{"2/5", "*", "120", "48"},

		{"1/2", "/", "2/1", "1/4"},
		{"3/5", "/", "5/3", "9/25"},
		{"3/5", "/", "3/5", "1"},

		{"1/2", "^", "2", "1/4"},
		{"2/3", "^", "3", "8/27"},
		{"0", "^", "1", "0"},
		{"1", "^", "0", "1"},
		{"-1", "^", "0", "1"},
		{"1/2", "^", "-1", "2"},
		{"1/4", "^", "-1/2", "2"},
		{"4/9", "^", "3/2", "8/27"},
		{"4/9", "^", "3/-2", "27/8"},
		{"8/27", "^", "2/3", "4/9"},
		{"8/27", "^", "-2/3", "9/4"},
		{"8/-27", "^", "1/3", "-2/3"},
		{"-1", "^", "3", "-1"},
		{"-1", "^", "0", "1"},
		{"-1", "^", "-1/3", "-1"},
		{"-1", "^", "1/3", "-1"},

		{"0", "!", "", "1"},
		{"1", "!", "", "1"},
		{"2", "!", "", "2"},
		{"3", "!", "", "6"},
		{"4", "!", "", "24"},
		{"5", "!", "", "120"},

		{"9", "sqrt", "", "3"},
		{"9/4", "sqrt", "", "3/2"},
		{"1", "sqrt", "", "1"},
		{"1000002000001", "sqrt", "", "1000001"},
	}
	for _, tc := range cases {
		switch tc.op {
		case "+":
			r := newRational(tc.a).Add(newRational(tc.b))
			assert.Equal(newRational(tc.r), r, "%s %s %s <> %s (got %s)", newRational(tc.a), tc.op, newRational(tc.b), newRational(tc.r), r)
		case "-":
			r := newRational(tc.a).Sub(newRational(tc.b))
			assert.Equal(newRational(tc.r), r, "%s %s %s <> %s (got %s) ", newRational(tc.a), tc.op, newRational(tc.b), newRational(tc.r), r)
		case "/":
			r := newRational(tc.a).Div(newRational(tc.b))
			assert.Equal(newRational(tc.r), r, "%s %s %s <> %s (got %s) ", newRational(tc.a), tc.op, newRational(tc.b), newRational(tc.r), r)
		case "*":
			r := newRational(tc.a).Mul(newRational(tc.b))
			assert.Equal(newRational(tc.r), r, "%s %s %s <> %s (got %s) ", newRational(tc.a), tc.op, newRational(tc.b), newRational(tc.r), r)
		case "^":
			r, err := newRational(tc.a).Pow(newRational(tc.b))
			assert.NoError(err)
			assert.Equal(newRational(tc.r), r, "%s %s %s <> %s (got %s)", newRational(tc.a), tc.op, newRational(tc.b), newRational(tc.r), r)
		case "!":
			r, err := newRational(tc.a).Fact()
			assert.NoError(err)
			assert.Equal(newRational(tc.r), r, "%s! <> %s (got %s)", newRational(tc.a), newRational(tc.r), r)
		case "--":
			r := newRational(tc.a).Minus()
			assert.Equal(newRational(tc.r), r, "-%s <> %s (got %s)", newRational(tc.a), newRational(tc.r), r)
		case "sqrt":
			r, err := newRational(tc.a).Sqrt()
			assert.NoError(err)
			assert.Equal(newRational(tc.r), r, "sqrt(%s) != %s (got %s)", newRational(tc.a), newRational(tc.r), r)
		}
	}
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
			_, err := newRational(tc.a).Fact()
			assert.Error(err, "%s!", newRational(tc.a))
		case "^":
			_, err := newRational(tc.a).Pow(newRational(tc.b))
			assert.Error(err, "%s ^ %s", newRational(tc.a), newRational(tc.b))
		case "sqrt":
			_, err := newRational(tc.a).Sqrt()
			assert.Error(err, "sqrt(%s)", newRational(tc.a))
		}
	}
}

func TestRationalMisc(t *testing.T) {
	assert := assert.New(t)

	assert.True(newRational("-1/2").Negative())
	assert.True(newRational("1/-2").Negative())
	assert.False(newRational("0").Negative())
	assert.False(newRational("2/1").Negative())

	assert.True(newRational("1/3").Less(newRational("1/2")))
	assert.False(newRational("1/3").Less(newRational("1/-2")))
	assert.False(newRational("-1/3").Less(newRational("1/-2")))
	assert.False(newRational("1/3").Less(newRational("-2")))

	assert.True(newRational("9").Integer())
	assert.True(newRational("-9/3").Integer())
	assert.False(newRational("1/2").Integer())

	assert.True(newRational("0/100").Zero())
	assert.False(newRational("1/200").Zero())

	assert.True(newRational("-2").Even())
	assert.True(newRational("2").Even())
	assert.True(newRational("0").Even())
	assert.False(newRational("-1").Even())
	assert.False(newRational("1").Even())
	assert.False(newRational("4/8").Even())
	assert.False(newRational("-4/8").Even())

	assert.True(newRational("3/9").Equal(newRational("1/3")))
	assert.True(newRational("-3/6").Equal(newRational("1/-2")))
	assert.False(newRational("-3/6").Equal(newRational("-2")))
}
