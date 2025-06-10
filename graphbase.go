package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Node struct {
	index int
	adj   map[int]*Node
}

type Arrow struct {
	u, v int
}

type Stack []*Node

func (s *Stack) Push(n *Node) {
	*s = append(*s, n)
}

func (s *Stack) Pop() *Node {
	l := len(*s)
	n := (*s)[l-1]
	*s = (*s)[:l-1]
	return n
}

func tarjan(graph []*Node) (map[int]int, int) {
	in := make(map[int]int)
	low := make(map[int]int)
	comp := make(map[int]int)
	for _, n := range graph {
		in[n.index], low[n.index], comp[n.index] = 0, 0, -1
	}
	stack := Stack(make([]*Node, 0))
	time, count := 1, 0

	var visit func(v *Node)
	visit = func(v *Node) {
		in[v.index], low[v.index] = time, time
		time++
		stack.Push(v)
		for _, u := range v.adj {
			if in[u.index] == 0 {
				visit(u)
			}
			if comp[u.index] == -1 && low[v.index] > low[u.index] {
				low[v.index] = low[u.index]
			}
		}
		if in[v.index] == low[v.index] {
			for {
				u := stack.Pop()
				comp[u.index] = count
				if u.index == v.index {
					break
				}
			}
			count++
		}
	}

	for _, v := range graph {
		if in[v.index] == 0 {
			visit(v)
		}
	}
	return comp, count
}

func scanInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	return n
}

func dfs(condensation []*Node) map[int]bool {
	kernel := make(map[int]bool)
	visited := make(map[int]bool)

	var visit func(v *Node)
	visit = func(v *Node) {
		visited[v.index] = true
		for _, u := range v.adj {
			if !visited[u.index] {
				visit(u)
			}
			delete(kernel, u.index)
		}
	}

	for _, scc := range condensation {
		if _, isVisited := visited[scc.index]; !isVisited {
			visit(scc)
			kernel[scc.index] = true
		}
	}
	return kernel
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	n, m := scanInt(scanner), scanInt(scanner)
	graph := make([]*Node, n)
	for i := 0; i < n; i++ {
		graph[i] = &Node{i, make(map[int]*Node, 0)}
	}
	arrows := make([]Arrow, m)
	for i := 0; i < m; i++ {
		u, v := scanInt(scanner), scanInt(scanner)
		graph[u].adj[v] = graph[v]
		arrows[i] = Arrow{u, v}
	}

	comp, count := tarjan(graph)
	condensation := make([]*Node, count)
	for i, _ := range condensation {
		condensation[i] = &Node{index:i, adj:make(map[int]*Node, 0)}
	}
	for i, _ := range arrows {
		if comp[arrows[i].u] == comp[arrows[i].v] {
			continue
		}
		condensation[comp[arrows[i].u]].adj[comp[arrows[i].v]] = condensation[comp[arrows[i].v]]
	}
	kernel := dfs(condensation)
	minCompIndex := make(map[int]int)
	for _, v := range graph {
		if _, ok := minCompIndex[comp[v.index]]; !ok {
			minCompIndex[comp[v.index]] = v.index
		}
	}
	answer := make([]int, 0, len(kernel))
	for index, _ := range kernel {
		answer = append(answer, minCompIndex[index])
	}
	sort.Ints(answer)

	writer := bufio.NewWriter(os.Stdout)
	for _, index := range answer {
		fmt.Fprintf(writer, "%d ", index)
	}
	writer.Flush()
}
