package main

import "os"

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
