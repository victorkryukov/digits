// This file contains code for pretty-printing nodes.
package main

import "fmt"

// return true if we need parenthesis around n.String() in places like _ */-^ n
func (n *Node) needParenthesis() bool {
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
