package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeValid(t *testing.T) {
	assert := assert.New(t)
	assert.True(newIntNode(1).valid())
	assert.True(newNode(newValNode(rational{3, 4}), OpSub, newIntNode(2)).valid())
	n := &Node{left: newIntNode(1), right: newIntNode(2), op: OpMinus}
	assert.False(n.valid())
	n = &Node{left: nil, right: newIntNode(2), op: OpFact}
	assert.False(n.valid())
}

func TestNodeDepth(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(0, newIntNode(1).Depth())
	n := newNode(newIntNode(1), OpSub, newIntNode(2))
	assert.Equal(1, n.Depth())
	assert.Equal(2, newNode(n, OpMinus, nil).Depth())
	assert.Equal(2, newNode(newIntNode(3), OpAdd, n).Depth())
}

func TestNodeEqual(t *testing.T) {
	assert := assert.New(t)
	n1 := newNode(newIntNode(1), OpSub, newIntNode(2))
	n2 := newNode(newIntNode(1), OpMinus, nil)
	n3 := newNode(newIntNode(2), OpSub, newIntNode(1))
	n4 := newNode(newIntNode(1), OpAdd, newIntNode(2))
	nodes := []*Node{n1, n2, n3, n4}
	for _, n := range nodes {
		assert.True(n.Equal(n))
	}
	for i := 0; i <= 3; i++ {
		for j := i + 1; j <= 3; j++ {
			assert.False(nodes[i].Equal(nodes[j]))
		}
	}
	n5 := newNode(newValNode(rational{3, 4}), OpAdd, newValNode(rational{-1, 2}))
	n6 := newNode(newValNode(rational{6, 8}), OpAdd, newValNode(rational{2, -4}))
	assert.True(n5.Equal(n5))
	assert.True(n5.Equal(n6))
	assert.True(n6.Equal(n6))
}
