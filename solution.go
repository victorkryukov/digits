package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

// start and end indicates that we used digits[start:end] for this formula,
// where digits is the original digits string. We cannot use binary operator for s1, s2
// if s1.end != s2.start
type Solution struct {
	val        Value
	start, end int
}

var NoSolution Solution

func init() {
	NoSolution = Solution{val: rational{}, start: -1, end: -1}
}

// Global variables are bad for you health
var solutions map[Solution][]*Node // solutions found so far
var maxDepth int64                 // If positive, only search for formulas of up to this level. If zero, only stores the first solution.

func init() {
	solutions = make(map[Solution][]*Node)
}

// Add adds a new formula for s, but only if it's unique and has reasonable depth.
// If v == nil, seed solutions with initial digits.
func (s Solution) Add(v *Node) {
	if maxDepth == 0 && solutions[s] != nil {

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
	if s.val.Zero() || (s.val.One() && op != OpMinus) {
		return s
	}
	v1, err := s.val.PerformUnary(op)
	if err != nil {
		return NoSolution
	}
	s1 := Solution{val: v1, start: s.start, end: s.end}
	for _, n := range solutions[s] {
		if n.op == OpMinus && op == OpMinus {
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
	v1, err := s1.val.PerformBinary(op, s2.val)
	if err != nil {
		return NoSolution
	}
	s3 := Solution{val: v1, start: s1.start, end: s2.end}
	for _, n1 := range solutions[s1] {
		for _, n2 := range solutions[s2] {
			if op == OpMinus && n2.op == OpMinus {
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
	if s.val.One() || s.val.MinusOne() {
		minusS, _ := s.val.PerformUnary(OpMinus)
		return SolutionSlice{s, Solution{minusS, s.start, s.end}}
	}
	result := SolutionSlice{s}
	s1 := s.Unary(OpMinus)
	if s1 != NoSolution {
		result = append(result, s1)
	}
	if s.val.Negative() {
		if s1 == NoSolution {
			return result
		} else {
			s = s1
		}
	}
	for f := s.Unary(OpFact); f != NoSolution; f = f.Unary(OpFact) {
		result = append(result, f)
		result = append(result, f.Unary(OpMinus))
	}
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
	s := Solution{val: rational{int64(n), 1}, start: start, end: start + len(a)}
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
		if min <= max && !f.val.IsInteger() || (f.val.Less(rational{min, 1}) || rational{max, 1}.Less(f.val)) {
			continue
		}
		if all {
			fmt.Printf("---\nAll formulas for number %s up to depth = %d:\n", f.val, maxDepth)
		} else {
			fmt.Printf("%s\t= ", f.val)
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
