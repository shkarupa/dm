package main 

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stack struct{ s []int }

func NewStack() Stack          { return Stack{make([]int, 0)} }
func (s *Stack) Push(i int)    { s.s = append(s.s, i) }
func (s *Stack) IsEmpty() bool { return len(s.s) == 0 }

func (s *Stack) Pop() int {
	i := s.s[len(s.s)-1]
	s.s = s.s[:len(s.s)-1]
	return i
}

type Machine struct {
	n, m, q int
	delta	map[int]map[int]int
	phi	map[int]map[int]string
}

func (a Machine) String() string {
	var b strings.Builder
	b.WriteString("digraph{\n")
	b.WriteString("\trankdir = LR\n")
	for i, _ := range a.delta {
		for j := 0; j < a.m; j++ {
			b.WriteString(
				fmt.Sprintf("\t%d -> %d [label = \"%c(%s)\"]\n",
					i, a.delta[i][j], 'a'+j, a.phi[i][j]))
		}
	}
	b.WriteString("}")
	return b.String()
}

func getCanonic(a Machine) Machine {
	processed := make(map[int]bool)
	canonic := make(map[int]int)
	stack := NewStack()
	stack.Push(a.q)
	for counter := 0; !stack.IsEmpty(); {
		q := stack.Pop()
		if !processed[q] {
			canonic[q] = counter
			counter++
			for i := a.m - 1; i >= 0; i-- {
				stack.Push(a.delta[q][i])
			}
			processed[q] = true
		}
	}
	delta, phi := make(map[int]map[int]int, a.n), make(map[int]map[int]string, a.n)
	for i := range a.delta {
		delta[canonic[i]], phi[canonic[i]] = make(map[int]int, a.m), make(map[int]string, a.m)
		for j := 0; j < a.m; j++ {
			delta[canonic[i]][j] = canonic[a.delta[i][j]]
			phi[canonic[i]][j] = a.phi[i][j]
		}
	}
	return Machine{a.n, a.m, canonic[a.q], delta, phi}
}

type Node struct {
	i,
	depth 	int
	parent	*Node
}

func makeSet(i int) *Node {
	t := &Node{}
	t.i, t.depth, t.parent = i, 0, t
	return t
}

func find(x *Node) *Node {
	if x.parent == x {
		return x
	}
	x.parent = find(x.parent)
	return x.parent
}

func union(x, y *Node) {
	rx, ry := find(x), find(y)
	if rx.depth < ry.depth {
		rx.parent = ry
	} else {
		ry.parent = rx
		if rx.depth == ry.depth && rx != ry {
			rx.depth++
		}
	}
}

func split1(a Machine) ([]int, int) {
	n := a.n
	sets := make([]*Node, a.n)
	for i := 0; i < a.n; i++ {
		sets[i] = makeSet(i)
	}
	for i := 0; i < a.n; i++ {
		for j := i + 1; j < a.n; j++ {
			if find(sets[i]) != find(sets[j]) {
				eq := true
				for x := 0; x < a.m; x++ {
					if a.phi[i][x] != a.phi[j][x] {
						eq = false
						break
					}
				}
				if eq {
					union(sets[i], sets[j])
					n--
				}
			}
		}
	}
	pi := make([]int, a.n)
	for i := 0; i < a.n; i++ {
		pi[i] = (find(sets[i])).i
	}
	return pi, n
}

func split(a Machine, pi []int) int {
	n := a.n
	sets := make([]*Node, a.n)
	for i := 0; i < a.n; i++ {
		sets[i] = makeSet(i)
	}
	for i := 0; i < a.n; i++ {
		for j := i + 1; j < a.n; j++ {
			if pi[i] == pi[j] && find(sets[i]) != find(sets[j]) {
				eq := true
				for x := 0; x < a.m; x++ {
					w1, w2 := a.delta[i][x], a.delta[j][x]
					if pi[w1] != pi[w2] {
						eq = false
						break
					}
				}
				if eq {
					union(sets[i], sets[j])
					n--
				}
			}
		}
	}
	for i := 0; i < a.n; i++ {
		pi[i] = (find(sets[i])).i
	}
	return n
}

func getMinimized(a Machine) Machine {
	pi, n := split1(a)
	nextN := 0
	for {
		nextN = split(a, pi)
		if n == nextN {
			break
		}
		n = nextN
	}
	processed := make(map[int]bool)
	delta, phi := make(map[int]map[int]int, a.n), make(map[int]map[int]string, a.n)
	for i := 0; i < a.n; i++ {
		q := pi[i]
		if !processed[q] {
			delta[q] = make(map[int]int, a.m)
			phi[q] = make(map[int]string, a.m)
			for x := 0; x < a.m; x++ {
				delta[q][x] = pi[a.delta[i][x]]
				phi[q][x] = a.phi[i][x]
			}
			processed[q] = true
		}
	}
	return Machine{nextN, a.m, pi[a.q], delta, phi}
}

func scanInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	return n
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	n, m, q0 := scanInt(scanner), scanInt(scanner), scanInt(scanner)
	delta, phi := make(map[int]map[int]int, n), make(map[int]map[int]string, n)
	for i := 0; i < n; i++ {
		delta[i] = make(map[int]int, m)
		for j := 0; j < m; j++ {
			delta[i][j] = scanInt(scanner)
		}
	}
	for i := 0; i < n; i++ {
		phi[i] = make(map[int]string, m)
		for j := 0; j < m; j++ {
			scanner.Scan()
			phi[i][j] = scanner.Text()
		}
	}
	a := Machine{n, m, q0, delta, phi}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, getCanonic(getMinimized(a)))
	writer.Flush()
}
