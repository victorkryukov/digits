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
