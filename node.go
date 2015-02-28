package main

import (
	"fmt"
	"log"
)

type Op byte // Operators

const (
	OpNull Op = iota

	// Binary ops start here
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPow

	// Unary ops start here
	OpFact
	OpSqrt
	OpMinus // unary minus
)

var opNames = map[Op]string{
	OpNull:  "NULL",
	OpAdd:   "+",
	OpSub:   "-",
	OpMul:   "*",
	OpDiv:   "/",
	OpPow:   "^",
	OpFact:  "!",
	OpSqrt:  "sqrt",
	OpMinus: "-",
}

// Node represents a formula parse tree, storing value (for a leaf) or
// operand with left and right sub-nodes. Nodes with unary operators will have their
// right sub-node nil, which is checked by Node.valid().
type Node struct {
	left, right *Node
	val         rational
	op          Op
}

// valid returns true for correct nodes. It does NOT check the subnodes recursively.
func (n *Node) valid() bool {
	if n.op == OpNull {
		return n.left == nil && n.right == nil
	} else if n.op <= OpPow {
		return n.left != nil && n.right != nil
	} else {
		return n.left != nil && n.right == nil
	}
}

// newNode creates a new formula Node. It panics if requested Node will be not valid.
func newNode(left *Node, op Op, right *Node) *Node {
	n := &Node{left: left, op: op, right: right}
	if !n.valid() {
		panic(fmt.Sprintf("Cannot create non-valid node: %v %v %v", left, op, right))
	}
	return n
}

// newValNode creates a new value Node from a rational.
func newValNode(val rational) *Node {
	return &Node{val: val}
}

// newIntNode creates a new value Node from an integer.
func newIntNode(val int64) *Node {
	return &Node{val: rational{n: val, d: 1}}
}

// Depth returns distance of the deepest leaf to the root.
func (n *Node) Depth() int64 {
	if n.op == OpNull {
		return 0
	}
	depth := n.left.Depth()
	if n.right != nil {
		if d := n.right.Depth(); d > depth {
			depth = d
		}
	}
	return depth + 1
}

// Equal returns true if two nodes have identical structure and leafs.
func (n *Node) Equal(n1 *Node) bool {
	if n1 == nil || n.op != n1.op {
		return false
	}
	if n.op != OpNull {
		return n.left.Equal(n1.left) && (n.right == nil || n.right.Equal(n1.right))
	} else {
		return n.val.Equal(n1.val)
	}
}

// transformDuo transorms all expressions of the form (op1 a) op2 (op3 b) into op4 (a op5 b),
// and leaves other expressions intact. In the form above, (OpNull x) is treated as x.
func (n *Node) transformDuo(op1, op2, op3, op4, op5 Op) *Node {
	var a, b *Node
	if n.op == op2 {
		if n.left.op == op1 && n.left.left != nil {
			a = n.left.left.Simplify()
		} else if op1 == OpNull {
			a = n.left.Simplify()
		} else {
			return n
		}
		if n.right.op == op3 && n.right.left != nil {
			b = n.right.left.Simplify()
		} else if op3 == OpNull {
			b = n.right.Simplify()
		} else {
			return n
		}
		n1 := &Node{op: op5, left: a, right: b}
		if op4 != OpNull {
			n1 = &Node{op: op4, left: n1.Simplify()}
		}
		return n1
	} else {
		return n
	}
}

// transformTrio transforms an expression of the form a op1 (b op2 c) into (a op3 b) op4 c,
// and leaves other expressions intact.
func (n *Node) transformTrio(op1, op2, op3, op4 Op) *Node {
	if n.op == op1 && n.right.op == op2 {
		n1 := &Node{op: op3, left: n.left.Simplify(), right: n.right.left.Simplify()}
		return &Node{op: op4, left: n1.Simplify(), right: n.right.right.Simplify()}
	} else {
		return n
	}
}

// Make various simplifications to convert n into a canonical form.
func (n *Node) Simplify() *Node {
	var n1 *Node
	if n.op == OpMinus && n.left.op == OpMinus {
		n1 = n.left.left.Simplify()
	} else if n.op == OpPow && n.left.op == OpMinus {
		e, err := n.right.Eval()
		if err == nil && e.Even() {
			n1 = &Node{op: OpPow, left: n.left.left.Simplify(), right: n.right.Simplify()}
		} else {
			n1 = n
		}
	} else {
		n1 = n
		for _, t := range [][5]Op{
			{OpNull, OpAdd, OpMinus, OpNull, OpSub},
			{OpNull, OpAdd, OpMinus, OpNull, OpSub},
			{OpNull, OpSub, OpMinus, OpNull, OpAdd},
			{OpMinus, OpSub, OpNull, OpMinus, OpAdd},
			{OpMinus, OpAdd, OpNull, OpMinus, OpSub},
			{OpMinus, OpMul, OpMinus, OpNull, OpMul},
			{OpMinus, OpDiv, OpMinus, OpNull, OpDiv},
			{OpMinus, OpMul, OpNull, OpMinus, OpMul},
			{OpMinus, OpDiv, OpNull, OpMinus, OpDiv},
			{OpNull, OpMul, OpMinus, OpMinus, OpMul},
			{OpNull, OpDiv, OpMinus, OpMinus, OpDiv},
			{OpSqrt, OpMul, OpSqrt, OpSqrt, OpMul},
			{OpSqrt, OpDiv, OpSqrt, OpSqrt, OpDiv},
		} {
			n1 = n1.transformDuo(t[0], t[1], t[2], t[3], t[4])
		}
		for _, t := range [][4]Op{
			{OpAdd, OpAdd, OpAdd, OpAdd},
			{OpSub, OpSub, OpSub, OpAdd},
			{OpMul, OpMul, OpMul, OpMul},

			{OpDiv, OpDiv, OpDiv, OpMul},
		} {
			n1 = n1.transformTrio(t[0], t[1], t[2], t[3])
		}
		if n1 == n {
			var l, r *Node
			if n.left != nil {
				l = n.left.Simplify()
			}
			if n.right != nil {
				r = n.right.Simplify()
			}
			if l != n.left || r != n.right {
				n1 = &Node{op: n.op, val: n.val, left: l, right: r}
			}
		}
	}
	if n1 != nil && n1 != n {
		e, err := n.Eval()
		if err != nil {
			log.Fatalf("Error evaluating %s: %s\n", n, err)
		}
		e1, err1 := n1.Eval()
		if err != nil {
			log.Fatalf("Error evaluating %s: %s\n", n1, err1)
		}
		if e != e1 {
			log.Fatalf("Before: n = %s\t[%d]\nAfter: n = %s\t[%d]\n", n, e, n1, e1)
		}
		return n1.Simplify()
	} else {
		return n
	}
}

// Eval returns the value of this node's expression
func (n *Node) Eval() (rational, error) {
	if n.op == OpNull {
		if n.left == nil && n.right == nil {
			return n.val, nil
		} else {
			return rational{}, fmt.Errorf("Undefined operation")
		}
	}
	if n.op < OpFact {
		if n.left == nil || n.right == nil {
			return rational{}, fmt.Errorf("Not enough arguments for %s", opNames[n.op])
		}
		left, err := n.left.Eval()
		if err != nil {
			return rational{}, err
		}
		right, err := n.right.Eval()
		if err != nil {
			return rational{}, err
		}
		switch n.op {
		case OpAdd:
			return left.Add(right), nil
		case OpSub:
			return left.Sub(right), nil
		case OpMul:
			return left.Mul(right), nil
		case OpDiv:
			if right.Zero() {
				return rational{}, fmt.Errorf("Division by 0")
			}
			return left.Div(right), nil
		case OpPow:
			return left.Pow(right)
		}
	} else {
		if n.right != nil {
			return rational{}, fmt.Errorf("Right operand for %s should be empty", opNames[n.op])
		}
		if n.left == nil {
			return rational{}, fmt.Errorf("Left operand for %s should NOT be empty", opNames[n.op])
		}
		left, err := n.left.Eval()
		if err != nil {
			return rational{}, err
		}
		switch n.op {
		case OpFact:
			return left.Fact()
		case OpSqrt:
			return left.Sqrt()
		case OpMinus:
			return left.Minus(), nil
		}
	}

	return rational{}, fmt.Errorf("Unreachable state")
}
