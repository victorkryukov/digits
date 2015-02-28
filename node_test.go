package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeValid(t *testing.T) {
	assert := assert.New(t)
	assert.True(newIntNode(1).valid())
	assert.True(newNode(newValNode(Rat{3, 4}), OpSub, newIntNode(2)).valid())
	assert.False(newNode(newIntNode(1), OpMinus, newIntNode(2)).valid())
	assert.False(newNode(nil, OpFact, newIntNode(2)).valid())
}

func TestNodeDepth(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(0, newIntNode(1).Depth())
	n := newNode(newIntNode(1), OpSub, newIntNode(2))
	assert.Equal(1, n.Depth())
	assert.Equal(2, newNode(n, OpMinus, nil).Depth())
	assert.Equal(2, newNode(newIntNode(3), OpAdd, n).Depth())
}
