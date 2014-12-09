package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
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
			return "<UNDEFINED>" // Unreachable
		}
	}
}

// String2 is for debugging purposes
func (n *Node) String2() string {
	left, right := "NIL", "NIL"
	if n.left != nil {
		left = n.left.String2()
	}
	if n.right != nil {
		right = n.right.String2()
	}
	return fmt.Sprintf("<Node %d %s, left = %s, right = %s>", n.val, opNames[n.op], left, right)
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
	if n.op == OpMinus && n.left.op == OpMinus { // - (- a) => a
		n1 = n.left.left.Simplify()
	} else if n.op == OpPow && n.left.op == OpMinus { // (-a) ^ b => a^b if b.Eval is even
		e, err := n.right.Eval()
		if err == nil && e.Even() {
			n1 = &Node{op: OpPow, left: n.left.left.Simplify(), right: n.right.Simplify()}
		} else {
			n1 = n
		}
	} else {
		n1 = n
		for _, t := range [][5]Op{
			{OpNull, OpAdd, OpMinus, OpNull, OpSub},  // a + (- b)  =>  a - b
			{OpNull, OpAdd, OpMinus, OpNull, OpSub},  // a + (- b)  =>  a - b
			{OpNull, OpSub, OpMinus, OpNull, OpAdd},  // a - (- b) => a + b
			{OpMinus, OpSub, OpNull, OpMinus, OpAdd}, // -a - b => - (a + b)
			{OpMinus, OpAdd, OpNull, OpMinus, OpSub}, // -a + b => - (a - b)
			{OpMinus, OpMul, OpMinus, OpNull, OpMul}, // (-a) * (-b) => a * b
			{OpMinus, OpDiv, OpMinus, OpNull, OpDiv}, // (-a) / (-b) => a / b
			{OpMinus, OpMul, OpNull, OpMinus, OpMul}, // (-a) * b => -(a * b)
			{OpMinus, OpDiv, OpNull, OpMinus, OpDiv}, // (-a) / b => -(a / b)
			{OpNull, OpMul, OpMinus, OpMinus, OpMul}, // a * (-b) => -(a * b)
			{OpNull, OpDiv, OpMinus, OpMinus, OpDiv}, // a * (-b) => -(a * b)
			{OpSqrt, OpMul, OpSqrt, OpSqrt, OpMul},   // sqrt(a) * sqrt(b) => sqrt(a * b)
			{OpSqrt, OpDiv, OpSqrt, OpSqrt, OpDiv},   // sqrt(a) / sqrt(b) => sqrt(a / b)
		} {
			n1 = n1.transformDuo(t[0], t[1], t[2], t[3], t[4])
		}
		for _, t := range [][4]Op{
			{OpAdd, OpAdd, OpAdd, OpAdd}, // a + (b + c) => (a + b) + c
			{OpSub, OpSub, OpSub, OpAdd}, // a - (b - c) => (a - b) + c
			{OpMul, OpMul, OpMul, OpMul}, // a * (b * c) => (a * b) * c
			// to avoid overflow {OpMul, OpDiv, OpMul, OpDiv}, // a * (b / c) => (a * b) / c
			{OpDiv, OpDiv, OpDiv, OpMul}, // a / (b / c) => (a / b) * c
		} {
			n1 = n1.transformTrio(t[0], t[1], t[2], t[3])
		}
		if n1 == n { // last attempt - simplify children
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
	if n.op < OpFact { // Binomial operator
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
	} else { // Single operator
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
	// should not be reached
	return BadRat, fmt.Errorf("Unreachable state")
}

// Solution for the formula.
// start and end indicates that we used digits[start:end] for this formula,
// where digits is the original digits string. We cannot use binary operator for s1, s2
// if s1.end != s2.start
type Solution struct {
	val        Rat
	start, end int
}

var NoSolution Solution

func init() {
	NoSolution = Solution{val: BadRat, start: -1, end: -1}
}

// Global variables are bad for you health
var solutions map[Solution][]*Node // solutions found so far
var maxDepth int64                 // If positive, only search for formulas of up to this level. If zero, only stores the first solution.

func init() {
	solutions = make(map[Solution][]*Node) // I love go's ability to have a struct as a key
}

// Add adds a new formula for s, but only if it's unique and has reasonable depth.
// If v == nil, seed solutions with initial digits.
func (s Solution) Add(v *Node) {
	if maxDepth == 0 && solutions[s] != nil {
		// Do nothing, we only need the first solution
		return
	}
	if v == nil {
		v = &Node{val: s.val}
	} else {
		v = v.Simplify()
	}
	if maxDepth != 0 && v.Depth() > maxDepth && solutions[s] != nil {
		return
	}
	for _, v1 := range solutions[s] {
		if v.Equal(v1) {
			return
		}
	}
	solutions[s] = append(solutions[s], v)
}

// Apply an unary operator to this solution, if possible, and add to all solutions
// found so far.
func (s Solution) Unary(op Op) Solution {
	if s.val.Zero() || (s.val.Equal(Rat{1, 1}) && op != OpMinus) {
		return s
	}
	var s1 Solution
	switch op {
	case OpMinus:
		s1 = Solution{val: s.val.Minus(), start: s.start, end: s.end}
	case OpFact:
		if s.val.Less(Rat{3, 1}) {
			return NoSolution
		}
		f, err := s.val.Fact()
		if err == nil {
			s1 = Solution{val: f, start: s.start, end: s.end}
		} else {
			return NoSolution
		}
	case OpSqrt:
		sq, err := s.val.Sqrt()
		if err == nil {
			s1 = Solution{val: sq, start: s.start, end: s.end}
		} else {
			return NoSolution
		}
	default:
		return NoSolution
	}
	for _, n := range solutions[s] {
		if n.op == OpMinus && op == OpMinus { // Do not apply to unary minus in a row
			continue
		}
		s1.Add(&Node{op: op, left: n})
	}
	return s1
}

// Apply a binary operator to this two solutions, if possible, and add to all solutions
// found so far. Returns a list of solutions that can be received this way.
func (s1 Solution) Binary(op Op, s2 Solution) Solution {
	if s1.end != s2.start {
		return NoSolution
	}
	var s3 Solution
	switch op {
	case OpAdd:
		s3 = Solution{
			val:   s1.val.Add(s2.val),
			start: s1.start, end: s2.end,
		}
	case OpSub:
		s3 = Solution{
			val:   s1.val.Sub(s2.val),
			start: s1.start, end: s2.end,
		}
	case OpMul:
		s3 = Solution{
			val:   s1.val.Mul(s2.val),
			start: s1.start, end: s2.end,
		}
	case OpDiv:
		if s2.val.Zero() {
			return NoSolution
		}
		s3 = Solution{
			val:   s1.val.Div(s2.val),
			start: s1.start, end: s2.end,
		}
	case OpPow:
		if p, err := s1.val.Pow(s2.val); err == nil {
			s3 = Solution{
				val:   p,
				start: s1.start, end: s2.end,
			}
		} else {
			return NoSolution
		}
	default:
		return NoSolution
	}
	for _, n1 := range solutions[s1] {
		for _, n2 := range solutions[s2] {
			if op == OpMinus && n2.op == OpMinus { // do not generate x - (-y)
				continue
			}
			s3.Add(&Node{
				op:    op,
				left:  n1,
				right: n2,
			})
		}
	}
	return s3
}

// SolutionSlice implements sort.Sortable interface
type SolutionSlice []Solution

func (p SolutionSlice) Len() int {
	return len(p)
}

func (p SolutionSlice) Less(i, j int) bool {
	return p[i].val.Less(p[j].val)
}

func (p SolutionSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p SolutionSlice) Sort() {
	sort.Sort(p)
}

// AllUnary generates all possible solutions we can get from s using unary operations,
// including itself (= no operation was applied).
func (s Solution) AllUnary() SolutionSlice {
	if s.val.Zero() {
		return SolutionSlice{s}
	}
	if s.val.Equal(Rat{1, 1}) || s.val.Equal(Rat{-1, 1}) {
		return SolutionSlice{s, Solution{s.val.Minus(), s.start, s.end}}
	}
	result := SolutionSlice{s}
	s1 := s.Unary(OpMinus)
	if s1 != NoSolution {
		result = append(result, s1)
	}
	if s.val.Negative() {
		if s1 == NoSolution {
			return result // we don't want to apply fact, sqrt to a negative number
		} else {
			s = s1
		}
	}
	// Factorial - it is never a perfect square
	for f := s.Unary(OpFact); f != NoSolution; f = f.Unary(OpFact) {
		result = append(result, f)
		result = append(result, f.Unary(OpMinus)) // And -s! was added to solution pool at this point
	}
	// Square root
	for sq := s.Unary(OpSqrt); sq != NoSolution; sq = sq.Unary(OpSqrt) {
		result = append(result, sq)
		for f := sq.Unary(OpFact); f != NoSolution; f = f.Unary(OpFact) {
			result = append(result, f)
		}
	}
	return uniq(result)
}

// AllBinary generates all possible binary solutions for s1 and s2.
func (s1 Solution) AllBinary(s2 Solution) SolutionSlice {
	result := SolutionSlice{}
	for op := OpAdd; op <= OpPow; op++ {
		for _, s3 := range s1.AllUnary() {
			for _, s4 := range s2.AllUnary() {
				if s5 := s3.Binary(op, s4); s5 != NoSolution {
					result = append(result, s5)
				}
			}
		}
	}
	return uniq(result)
}

// uniq returns only unique solutions from the list
func uniq(l SolutionSlice) SolutionSlice {
	m := make(map[Solution]bool)
	for _, n := range l {
		m[n] = true
	}
	var result SolutionSlice
	for k := range m {
		result = append(result, k)
	}
	return result
}

func atos(a string, start, end int) Solution {
	n, err := strconv.Atoi(a)
	if err != nil {
		log.Fatalf("Cannot convert %s to number\n", a)
	}
	s := Solution{val: Rat{int64(n), 1}, start: start, end: start + len(a)}
	s.Add(nil)
	return s
}

func atoi(s string) int64 {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Cannot convert %s to number\n", s)
	}
	return int64(n)
}

// min > max is a special case - to print all numbers
func (p SolutionSlice) Print(all bool, min, max int64) {
	p.Sort()
	for _, f := range p {
		if min <= max && !f.val.Integer() || (f.val.Less(Rat{min, 1}) || Rat{max, 1}.Less(f.val)) {
			continue
		}
		if all {
			fmt.Printf("---\nAll formulas for number %d up to depth = %d:\n", f.val, maxDepth)
		} else {
			fmt.Printf("%d\t= ", f.val)
		}
		answer := []string{}
		for _, n := range solutions[f] {
			answer = append(answer, fmt.Sprintf("[%2d] %s", n.Depth(), n))
		}
		sort.Strings(answer)
		fmt.Printf("%s\n", strings.Join(answer, "\n"))
	}
}

func FindAllSolutions(digits string, start int) SolutionSlice {
	if len(digits) == 0 {
		return nil
	}
	r := atos(digits, start, start+len(digits)).AllUnary()
	for i := 1; i < len(digits); i++ {
		for _, s1 := range FindAllSolutions(digits[:i], start) {
			for _, s2 := range FindAllSolutions(digits[i:], start+i) {
				r = append(r, s1.AllBinary(s2)...)
			}
		}
	}
	return uniq(r)
}

func main() {
	digits := os.Args[1]
	min := atoi(os.Args[2])
	max := atoi(os.Args[3])
	maxDepth = atoi(os.Args[4])
	var p SolutionSlice
	for _, s := range FindAllSolutions(digits, 0) {
		p = append(p, s.AllUnary()...)
	}
	p = uniq(p)
	DEBUG = len(os.Args) > 5
	p.Print(maxDepth > 0, min, max)
}

// func main() {
// 	var printAll bool
// 	if len(os.Args) > 4 {
// 		maxDepth = atoi(os.Args[4])
// 		printAll = true
// 	}
// 	formulas := generateFormulas(os.Args[1], atoi(os.Args[2]), atoi(os.Args[3]))
// 	SolutionSlice(formulas).Print(printAll)
// }
