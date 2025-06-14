package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"strconv"
)

type Stack struct{ s []map[int]bool }

func NewStack() Stack { return Stack{make([]map[int]bool, 0)} }
func (s *Stack) Push(i map[int]bool) { s.s = append(s.s, i) }
func (s *Stack) IsEmpty() bool { return len(s.s) == 0 }

func (s *Stack) Pop() map[int]bool {
	i := s.s[len(s.s)-1]
	s.s = s.s[:len(s.s)-1]
	return i
}

func closure(delta map[int]map[string]map[int]bool, z map[int]bool) map[int]bool {
	reached := make(map[int]bool)
	for q := range z {
		dfs(delta, q, reached)
	}
	return reached
}

func dfs(delta map[int]map[string]map[int]bool, q int, reached map[int]bool) {
	if !reached[q] {
		reached[q] = true
		for w := range delta[q]["lambda"] {
			dfs(delta, w, reached)
		}
	}
}

func det(nfa NFA) FA {
	hash := func(set map[int]bool) int {
		if len(set) == 0 {
			return -1
		}
		ordered := make([]int, 0, len(set))
		for i := range set {
			ordered = append(ordered, i)
		}
		sort.Ints(ordered)
		dec := 0
		for i := len(ordered) - 1; i > 0; i-- {
			dec = (dec + ordered[i]) * len(nfa.delta)
		}
		return dec + ordered[0]
	}
	q0 := closure(nfa.delta, map[int]bool{nfa.q0:true})
	states := map[int]bool{hash(q0):true}
	final := make(map[int]bool)
	delta := make(map[int]map[string]int)
	stack := NewStack()
	stack.Push(q0)
	for !stack.IsEmpty() {
		z := stack.Pop()
		for u := range z {
			if nfa.final[u] {
				final[hash(z)] = true
				break
			}
		}
		for a := range nfa.alphabet {
			if a == "lambda" {
				continue
			}
			set := make(map[int]bool)
			for u := range z {
				for q := range nfa.delta[u][a] {
					set[q] = true
				}
			}
			za := closure(nfa.delta, set)
			if !states[hash(za)] {
				states[hash(za)] = true
				stack.Push(za)
			}
			if _, ok := delta[hash(z)]; !ok {
				delta[hash(z)] = make(map[string]int)
			}
			delta[hash(z)][a] = hash(za)
		}
	}
	return FA{states, delta, nfa.alphabet, final, hash(q0)}
}

func scanInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	return n
}

type FA struct {
	states   map[int]bool
	delta    map[int]map[string]int
	alphabet map[string]int
	final    map[int]bool
	q0       int
}

type NFA struct {
	delta    map[int]map[string]map[int]bool
	alphabet map[string]int
	final    []bool
	q0       int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	n, m := scanInt(scanner), scanInt(scanner)
	var nfa NFA
	nfa.delta = make(map[int]map[string]map[int]bool)
	nfa.alphabet = make(map[string]int)
	for i := 0; i < n; i++ {
		nfa.delta[i] = make(map[string]map[int]bool)
	}
	for i, counter := 0, 0; i < m; i++ {
		q1, q2 := scanInt(scanner), scanInt(scanner)
		scanner.Scan()
		symbol := scanner.Text()
		if _, ok := nfa.delta[q1][symbol]; !ok {
			nfa.delta[q1][symbol] = make(map[int]bool)
		}
		nfa.delta[q1][symbol][q2] = true
		if _, ok := nfa.alphabet[symbol]; !ok {
			nfa.alphabet[symbol] = counter
			counter++
		}
	}
	nfa.final = make([]bool, n)
	for i := 0; i < n; i++ {
		nfa.final[i] = scanInt(scanner) == 1
	}
	nfa.q0 = scanInt(scanner)

	fa := det(nfa)

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, "digraph {")
	fmt.Fprintf(writer, "\trankdir = LR\n")
	counter := 0
	dict := make(map[int]int)
	for hash := range fa.states {
		fmt.Fprintf(writer, "\t%d ", counter)
		dict[hash] = counter
		fmt.Fprintf(writer, "[label = ")
		fmt.Fprintf(writer, "\"[")
		if hash != -1 {
			i := hash
			for i > n {
				fmt.Fprintf(writer, "%d ", i % n)
				i /= n
			}
			fmt.Fprintf(writer, "%d", i)
		}
		fmt.Fprintf(writer, "]\"")
		fmt.Fprintf(writer, ", shape = ")
		if fa.final[hash] {
			fmt.Fprintf(writer, "doublecircle")
		} else {
			fmt.Fprintf(writer, "circle")
		}
		fmt.Fprintf(writer, "]\n")
		counter++
	}
	for hash := range fa.delta {
		arrows := make(map[int][]string)
		for symbol := range fa.delta[hash] {
			arrows[fa.delta[hash][symbol]] = append(arrows[fa.delta[hash][symbol]], symbol)
		}
		for q := range arrows {
			sort.Slice(arrows[q],
			           func(a, b int) bool {
                           return nfa.alphabet[arrows[q][a]] < nfa.alphabet[arrows[q][b]]
					   })
		}
		for q := range arrows {
			fmt.Fprintf(writer, "\t%d -> %d ", dict[hash], dict[q])
			fmt.Fprintf(writer, "[label = \"%s\"]", strings.Join(arrows[q], ", "))
			fmt.Fprintf(writer, "\n")
		}
	}
	fmt.Fprintln(writer, "}")
	writer.Flush()
}
