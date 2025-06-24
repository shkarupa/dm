package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func scanInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	return n
}

func scanString(scanner *bufio.Scanner) string {
	scanner.Scan()
	s := scanner.Text()
	return s
}

type Node struct {
	index,
	dist  int
	adj   map[int]*Node
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

func tarjan(graph []*Node, q0 int) (map[int]int, map[int]int) {
	comp := make(map[int]int)
	compSize := make(map[int]int)
	in := make(map[int]int)
	low := make(map[int]int)
	stack := Stack(make([]*Node, 0))
	time, count := 1, 1
	var visit func(v *Node)
	visit = func(v *Node) {
		in[v.index], low[v.index] = time, time
		time++
		stack.Push(v)
		for _, u := range v.adj {
			if in[u.index] == 0 {
				visit(u)
			}
			if comp[u.index] == 0 && low[v.index] > low[u.index] {
				low[v.index] = low[u.index]
			}
		}
		if in[v.index] == low[v.index] {
			for {
				u := stack.Pop()
				comp[u.index] = count
				compSize[count]++
				if u.index == v.index {
					break
				}
			}
			count++
		}
	}
	visit(graph[q0])
	return comp, compSize
}

func order(graph []*Node, q0 int) []*Node {
	visited := make(map[int]bool, len(graph))
	ordered := make([]*Node, 0)
	var visit func(*Node)
	visit = func(v *Node) {
		visited[v.index] = true
		for _, u := range v.adj {
			if !visited[u.index] { 
				visit(u)
			}
		}
		ordered = append(ordered, v)
	}

	visit(graph[q0])
	for i, j := 0, len(ordered) - 1; i < j; i, j = i + 1, j - 1 {
		ordered[i], ordered[j] = ordered[j], ordered[i]
	}
	return ordered
}

func dfs(graph []*Node, q0 int) []bool {
	visited := make([]bool, len(graph))
	var visit func(v *Node)
	visit = func(v *Node) {
		visited[v.index] = true
		for j := range v.adj {
			if !visited[v.adj[j].index] {
				visit(v.adj[j])
			}
		}
	}
	visit(graph[q0])
	return visited
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	m, n, q0 := scanInt(scanner), scanInt(scanner), scanInt(scanner)
	final := make([]bool, n)
	graph := make([]*Node, n)
	for i := 0; i < n; i++ {
		graph[i] = &Node{i, 0, make(map[int]*Node, m)}
	}
	for i := 0; i < n; i++ {
		sign := scanString(scanner)
		final[i] = sign == "+"
		graph[i].adj = make(map[int]*Node, m)
		for j := 0; j < m; j++ {
			graph[i].adj[j] = graph[scanInt(scanner)]
		}
	}

	writer := bufio.NewWriter(os.Stdout)

	comp, compSize := tarjan(graph, q0)
	writer.Flush()
	inf := false
	for i := range graph {
		buckle := false
		reaches := false
		for j := range graph[i].adj {
			if graph[i].adj[j].index == i {
				buckle = true
			}
			if final[graph[i].adj[j].index] {
				reaches = true
			}
		}
		if compSize[comp[i]] > 1 {
			if reaches {
				inf = true
			}
		}

		inf = inf || buckle && reaches
	}
	if inf {
		fmt.Fprintln(writer, "INF")
		writer.Flush()
		return
	}

	ordered := order(graph, q0)
	for i := range ordered {
		for j := range graph[i].adj {
			if ordered[i].adj[j].dist < ordered[i].dist + 1 {
				ordered[i].adj[j].dist = ordered[i].dist + 1
			}
		}
	}
	visited := dfs(graph, q0)
	empty := true
	for i := range visited {
		if visited[i] && final[i] {
			empty = false
		}
	}
	if empty {
		fmt.Fprintln(writer, "EMPTY")
		writer.Flush()
		return
	}

	max := 0
	for i := range final {
		if !final[i] {
			continue
		}
		if graph[i].dist > max {
			max = graph[i].dist
		}
	}
	fmt.Fprintln(writer, max)
	writer.Flush()
}
