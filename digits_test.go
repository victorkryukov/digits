package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeString(t *testing.T) {
	assert := assert.New(t)
	a := &Node{
		op: OpMinus,
		left: &Node{
			op: OpFact,
			left: &Node{
				op: OpSqrt,
				left: &Node{
					val: Rat{9, 1},
				},
			},
		},
	}
	assert.Equal("-(sqrt(9)!)", a.String())
}
