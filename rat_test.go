package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRat(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Rat{1, 2}, Rat{-10, -20}.Normalize())
	assert.True(Rat{5, 6}.Equal(Rat{10, 20}.Add(Rat{3, 9})))
	p, err := Rat{3, 1}.Pow(Rat{2, 1})
	assert.NoError(err)
	assert.Equal(Rat{9, 1}, p)
	s, err := p.Sqrt()
	assert.NoError(err)
	assert.Equal(Rat{3, 1}, s)
	f, err := s.Fact()
	assert.NoError(err)
	assert.Equal(Rat{6, 1}, f)
	assert.True(s.Add(s.Minus()).Zero())
	assert.True(s.Sub(s).Zero())
	assert.Equal(Rat{9, 1}, s.Mul(s))
	assert.Equal(Rat{1, 9}, Rat{1, 1}.Div(s.Mul(s)))
	p, err = Rat{1, 9}.Pow(Rat{0, 1})
	assert.NoError(err)
	assert.Equal(Rat{1, 1}, p)
}
