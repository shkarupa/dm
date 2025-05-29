package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Node struct {
	value,
	color 	int
	adj 	[]*Node
}

type ByLex [][]int

func (a ByLex) Len() int      { return len(a) }
func (a ByLex) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByLex) Less(i, j int) bool {
	m := 0
	if len(a[i]) < len(a[j]) {
		m = len(a[i])
	} else {
		m = len(a[j])
	}
	for k := 0; k < m; k++ {
		if a[i][k] < a[j][k] {
			return true
		}
		if a[i][k] != a[j][k] {
			return false
		}
	}
	if len(a[i]) < len(a[j]) {
		return true
	}
	return false
}

func dfs(graph []*Node, component *([]*Node), v *Node, color int) bool {
	v.color = color
	ok := true
	*component = append(*component, v)
	for i := 0; i < len(v.adj) && ok; i++ {
		u := v.adj[i]
		if u.color == 0 {
			ok = dfs(graph, component, u, -color)
		} else if u.color != -color {
			ok = false
		}
	}
	return ok
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	graph := make([]*Node, n)
	for i := 0; i < n; i++ {
		graph[i] = &Node{i + 1, 0, make([]*Node, 0, 8)}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			scanner.Scan()
			if scanner.Text()[0] == '+' {
				graph[i].adj = append(graph[i].adj, graph[j])
			}
		}
	}

	writer := bufio.NewWriter(os.Stdout)

	ok := true
	components := make([][]*Node, 0)
	for i := 0; i < n && ok; i++ {
		if graph[i].color == 0 {
			components = append(components, make([]*Node, 0))
			ok = dfs(graph, &(components[len(components)-1]), graph[i], 1)
		}
	}

	if !ok {
		fmt.Fprintf(writer, "No solution\n")
		writer.Flush()
		return
	}

	crews := make([][]int, 1<<len(components))
	for i := 0; i < 1<<len(components); i++ {
		crews[i] = make([]int, 0)
		for j := 0; j < len(components); j++ {
			for k := 0; k < len(components[j]); k++ {
				if (components[j][k].color + 1) / 2 == i>>j%2 {
					crews[i] = append(crews[i], components[j][k].value)
				}
			}
		}
		sort.Ints(crews[i])
	}
	sort.Sort(ByLex(crews))
	var answer []int
	for _, crew := range crews {
		if len(crew) > len(answer) && len(crew) <= n / 2 {
			answer = crew
		}
	}
	for _, v := range answer {
		fmt.Fprintf(writer, "%d ", v)
	}
	writer.Flush()
}
