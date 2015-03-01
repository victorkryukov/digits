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

func parsedEqual(s string, n *Node) bool {
	n1, err := ParsePolish(s)
	return err == nil && n1.Equal(n)
}

// func TestNodeParse1(t *testing.T) {
// 	const s = "/ sqrt 2 + ! 3/4 -- -5/6"
// 	n, err := ParsePolish(s)
// 	fmt.Printf("s = %s\nn = %s\nerr = %s\n", s, n, err)
// }

func TestNodeParsePolish(t *testing.T) {
	assert := assert.New(t)
	assert.True(parsedEqual("* + 1/2 -3/4 - 5/6 7/8",
		newNode(
			newNode(newValNode(rational{1, 2}), OpAdd, newValNode(rational{-3, 4})),
			OpMul,
			newNode(newValNode(rational{5, 6}), OpSub, newValNode(rational{7, 8})))))
	assert.True(parsedEqual("/ sqrt 2 ^ ! 3/4 -- -5/6",
		newNode(
			newNode(newIntNode(2), OpSqrt, nil),
			OpDiv,
			newNode(
				newNode(newValNode(rational{3, 4}), OpFact, nil),
				OpPow,
				newNode(newValNode(rational{-5, 6}), OpMinus, nil)))))
	for _, s := range []string{
		"",
		"+-",
		"+",
		"+ 1",
		"+ 1 --",
		"x 1 2",
	} {
		_, err := ParsePolish(s)
		assert.Error(err, "Parsing '%s'", s)
	}
}

func TestNodeEval(t *testing.T) {
	assert := assert.New(t)
	for _, tc := range []struct{ s, v string }{
		{"+ 1/2 1/3", "5/6"},
		{"- 1/2 1/3", "1/6"},
		{"- 1/2 -1/3", "5/6"},
		{"- 1/2 -- 1/3", "5/6"},
		{"* -1/2 -1/3", "1/6"},
		{"sqrt 4", "2"},
		{"! 5", "120"},
		{"/ 1 100", "1/100"},
		{"^ 2 4", "16"},
		{"^ 4 1/2", "2"},
	} {
		n, err := ParsePolish(tc.s)
		assert.NoError(err)
		v, err := n.Eval()
		assert.NoError(err)
		assert.Equal(v, newRational(tc.v))
	}

	n1 := &Node{left: nil, op: OpAdd, right: newIntNode(5)}
	n2 := &Node{left: n1, op: OpFact, right: nil}
	n3 := &Node{left: n1, op: OpAdd, right: newIntNode(6)}
	n4 := &Node{left: newIntNode(6), op: OpAdd, right: n2}
	n5, _ := ParsePolish("/ 1 0")
	n6, _ := ParsePolish("^ 4 1/3")
	n7, _ := ParsePolish("! -2")
	n8, _ := ParsePolish("sqrt -2")
	for _, n := range []*Node{n1, n2, n3, n4, n5, n6, n7, n8} {
		_, err := n.Eval()
		assert.Error(err)
	}
}
