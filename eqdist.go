package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Node struct {
	val		int
	level	int
	mark	bool
	adj		map[int]bool
}

func NewNode(i int) *Node {
	return &Node{ i, -1, false, make(map[int]bool, 0) }
}

func intersect(a, b map[int]bool) map[int]bool {
	set := make(map[int]bool, 0)
	for k, _ := range a {
		if b[k] {
			set[k] = true
		}
	}
	return set
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func scanInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	return n
}

func bfs(graph []*Node, w *Node) []map[int]bool {
	for _, v := range graph {
		v.mark = false
		v.level = -1
	}
	w.mark = true
	w.level = 0
	level := 0
	set := make(map[int]bool, 0)
	sets := make([]map[int]bool, 0)
	queue := make([]int, 0)
	queue = append(queue, w.val)
	for len(queue) != 0 {
		vVal := queue[0]
		v := graph[vVal]
		queue = queue[1:]
		if v.level != level {
			level = v.level
			sets = append(sets, set)
			set = make(map[int]bool)
		}
		set[v.val] = true
		for uVal, _ := range v.adj {
			u := graph[uVal]
			if !u.mark {
				u.mark = true
				u.level = v.level + 1
				queue = append(queue, u.val)
			}
		}
	}
	sets = append(sets, set)
	return sets
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	n, m := scanInt(scanner), scanInt(scanner)
	graph := make([]*Node, n)
	var i int
	for i = 0; i < n; i++ {
		graph[i] = NewNode(i)
	}
	for i = 0; i < m; i++ {
		u, v := scanInt(scanner), scanInt(scanner)
		graph[u].adj[v] = true
		graph[v].adj[u] = true
	}

	k := scanInt(scanner)
	pivots := make([]int, k)
	for i = 0; i < k; i++ {
		pivots[i] = scanInt(scanner)
	}

	writer := bufio.NewWriter(os.Stdout)

	sets := bfs(graph, graph[pivots[0]])
	var limit int
	for i = 1; i < k; i++ {
		stes := bfs(graph, graph[pivots[i]])
		limit = min(len(sets), len(stes))
		for j := 0; j < limit; j++ {
			sets[j] = intersect(sets[j], stes[j])
		}
	}

	equidistant := make([]int, 0)
	for i := 0; i < limit; i++ {
		for k, _ := range sets[i] {
			equidistant = append(equidistant, k)
		}
	}
	if len(equidistant) == 0 {
		fmt.Fprintf(writer, "-\n")
	} else {
		sort.Ints(equidistant)
		for i = 0; i < len(equidistant) - 1; i++ {
			fmt.Fprintf(writer, "%d ", equidistant[i])
		}
		fmt.Fprintf(writer, "%d\n", equidistant[i])
	}
	writer.Flush()
}
