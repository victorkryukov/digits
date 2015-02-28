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
	OpFact  // factorial
	OpSqrt  // square root
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

type Node struct {
	left, right *Node
	val         Rat // Number stored in this node

	// Operation to be perfomed on this node. Require left != nil && right != nil for binary
	// and left != nil && right = nil for unary operators.
	op Op
}

var DEBUG bool

// return true if we need parenthesis around n.String() in places like _ */-^ n
func (n *Node) needParenthesis() bool {
	if DEBUG {
		return true
	}
	return n.op >= OpAdd && n.op <= OpPow
}

// String returns a formula for n, sometimes (always *sigh*) with excessive paranthesis
func (n *Node) String() string {
	var left, right string
	if n.left != nil {
		left = n.left.String()
	}
	if n.right != nil {
		right = n.right.String()
		if n.right.needParenthesis() {
			right = "(" + right + ")"
		}
	}
	if left == "" && right == "" {
		return n.val.String()
	} else {
		switch n.op {
		case OpAdd:
			return fmt.Sprintf("%s + %s", left, right)
		case OpSub:
			return fmt.Sprintf("%s - %s", left, right)
		case OpMul, OpDiv:
			if n.left.needParenthesis() {
				left = "(" + left + ")"
			}
			return fmt.Sprintf("%s %s %s", left, opNames[n.op], right)
		case OpPow:
			if n.left.needParenthesis() || n.left.op == OpMinus {
				left = "(" + left + ")"
			}
			return fmt.Sprintf("%s %s %s", left, opNames[n.op], right)
		case OpFact:
			if n.left.needParenthesis() || n.left.op == OpMinus {
				left = "(" + left + ")"
			}
			return fmt.Sprintf("%s!", left)
		case OpSqrt:
			return fmt.Sprintf("sqrt(%s)", left)
		case OpMinus:
			if n.left.needParenthesis() || n.left.op == OpFact {
				left = "(" + left + ")"
			}
			return "-" + left
		default:
			return "<UNDEFINED>"
		}
	}
}

func (n *Node) Depth() int64 {
	if n.left == nil && n.right == nil {
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

func (n *Node) Equal(n1 *Node) bool {
	if n.val != n1.val || n.op != n1.op {
		return false
	}
	if (n.left == nil && n1.left != nil) ||
		(n.left != nil && n1.left == nil) ||
		(n.left != nil && n1.left != nil && !n.left.Equal(n1.left)) {
		return false
	}
	if (n.right == nil && n1.right != nil) ||
		(n.right != nil && n1.right == nil) ||
		(n.right != nil && n1.right != nil && !n.right.Equal(n1.right)) {
		return false
	}
	return true
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
func (n *Node) Eval() (Rat, error) {
	if n.op == OpNull {
		if n.left == nil && n.right == nil {
			return n.val, nil
		} else {
			return BadRat, fmt.Errorf("Undefined operation")
		}
	}
	if n.op < OpFact {
		if n.left == nil || n.right == nil {
			return BadRat, fmt.Errorf("Not enough arguments for %s", opNames[n.op])
		}
		left, err := n.left.Eval()
		if err != nil {
			return BadRat, err
		}
		right, err := n.right.Eval()
		if err != nil {
			return BadRat, err
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
				return BadRat, fmt.Errorf("Division by 0")
			}
			return left.Div(right), nil
		case OpPow:
			return left.Pow(right)
		}
	} else {
		if n.right != nil {
			return BadRat, fmt.Errorf("Right operand for %s should be empty", opNames[n.op])
		}
		if n.left == nil {
			return BadRat, fmt.Errorf("Left operand for %s should NOT be empty", opNames[n.op])
		}
		left, err := n.left.Eval()
		if err != nil {
			return BadRat, err
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

	return BadRat, fmt.Errorf("Unreachable state")
}
