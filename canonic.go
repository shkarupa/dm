package main

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
)

type Stack struct { s []int }

func (s *Stack) Push(i int) {
	s.s = append(s.s, i)
}

func (s *Stack) IsEmpty() bool {
	return len(s.s) == 0
}

func (s *Stack) Pop() int {
	i := s.s[len(s.s) - 1]
	s.s = s.s[:len(s.s) - 1]
	return i
}

func NewStack() Stack {
	return Stack{ make([]int, 0) }
}


func scanInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	return n
}

func makeCanonic(delta [][]int, q0 int) []int {
	processed := make([]bool, len(delta))
	canonic := make([]int, len(delta))
	counter := 0
	stack := NewStack()
	stack.Push(q0)
	for (! stack.IsEmpty()) {
		q := stack.Pop()
		if (! processed[q]) {
			canonic[q] = counter
			counter++
			for i := len(delta[q]) - 1; i >= 0; i-- {
				stack.Push(delta[q][i])
			}
			processed[q] = true
		}
	}
	return canonic
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	n, m, q0 := scanInt(scanner), scanInt(scanner), scanInt(scanner)
	output := bufio.NewWriter(os.Stdout)

	delta := make([][]int, n)
	for i := 0; i < n; i++ {
		delta[i] = make([]int, m)
		for j := 0; j < m; j++ { delta[i][j] = scanInt(scanner) }
	}

	phi := make([][]string, n)
	for i := 0; i < n; i++ {
		phi[i] = make([]string, m)
		for j := 0; j < m; j++ {
			scanner.Scan()
			phi[i][j] = scanner.Text()
		}
	}

	canonic := makeCanonic(delta, q0)
	canonicDelta := make([][]int, n)
	canonicPhi := make([][]string, n)
	for i := 0; i < n; i++ {
		canonicDelta[canonic[i]] = make([]int, m)
		copy(canonicDelta[canonic[i]], delta[i])
		for j := 0; j < m; j++ {
			canonicDelta[canonic[i]][j] = canonic[delta[i][j]]
		}
		canonicPhi[canonic[i]] = make([]string, m)
		copy(canonicPhi[canonic[i]], phi[i])
	}

	fmt.Fprintf(output, "%d\n%d\n%d\n", n, m, canonic[q0])
	for i := 0; i < n; i++ {
		for j := 0; j < m - 1; j++ {
			fmt.Fprintf(output, "%d ", canonicDelta[i][j])
		}
		fmt.Fprintf(output, "%d\n", canonicDelta[i][m - 1])
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m - 1; j++ {
			fmt.Fprintf(output, "%s ", canonicPhi[i][j])
		}
		fmt.Fprintf(output, "%s\n", canonicPhi[i][m - 1])
	}
	output.Flush()
}
