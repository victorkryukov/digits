package main

import (
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func rat(s string) rational {
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

type testCase struct {
	a, op, b, r string
}

func TestrationalOps(t *testing.T) {
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
			r := rat(tc.a).Add(rat(tc.b))
			assert.Equal(rat(tc.r), r, "%s %s %s <> %s (got %s)", rat(tc.a), tc.op, rat(tc.b), rat(tc.r), r)
		case "-":
			r := rat(tc.a).Sub(rat(tc.b))
			assert.Equal(rat(tc.r), r, "%s %s %s <> %s (got %s) ", rat(tc.a), tc.op, rat(tc.b), rat(tc.r), r)
		case "/":
			r := rat(tc.a).Div(rat(tc.b))
			assert.Equal(rat(tc.r), r, "%s %s %s <> %s (got %s) ", rat(tc.a), tc.op, rat(tc.b), rat(tc.r), r)
		case "*":
			r := rat(tc.a).Mul(rat(tc.b))
			assert.Equal(rat(tc.r), r, "%s %s %s <> %s (got %s) ", rat(tc.a), tc.op, rat(tc.b), rat(tc.r), r)
		case "^":
			r, err := rat(tc.a).Pow(rat(tc.b))
			assert.NoError(err)
			assert.Equal(rat(tc.r), r, "%s %s %s <> %s (got %s)", rat(tc.a), tc.op, rat(tc.b), rat(tc.r), r)
		case "!":
			r, err := rat(tc.a).Fact()
			assert.NoError(err)
			assert.Equal(rat(tc.r), r, "%s! <> %s (got %s)", rat(tc.a), rat(tc.r), r)
		case "--":
			r := rat(tc.a).Minus()
			assert.Equal(rat(tc.r), r, "-%s <> %s (got %s)", rat(tc.a), rat(tc.r), r)
		case "sqrt":
			r, err := rat(tc.a).Sqrt()
			assert.NoError(err)
			assert.Equal(rat(tc.r), r, "sqrt(%s) != %s (got %s)", rat(tc.a), rat(tc.r), r)
		}
	}
}

func TestrationalErrors(t *testing.T) {
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

func TestrationalMisc(t *testing.T) {
	assert := assert.New(t)

	assert.True(rat("-1/2").Negative())
	assert.True(rat("1/-2").Negative())
	assert.False(rat("0").Negative())
	assert.False(rat("2/1").Negative())

	assert.True(rat("1/3").Less(rat("1/2")))
	assert.False(rat("1/3").Less(rat("1/-2")))
	assert.False(rat("-1/3").Less(rat("1/-2")))
	assert.False(rat("1/3").Less(rat("-2")))

	assert.True(rat("9").Integer())
	assert.True(rat("-9/3").Integer())
	assert.False(rat("1/2").Integer())

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
