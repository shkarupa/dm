package main

import "fmt"

type Node struct {
	val,		
	mark,	
	comp	int
	parent	*Node
	adj		[]*Node
}

type Edge struct {
	u, v	int
}

const (
	WHITE = iota
	GRAY
	BLACK
)

type Queue []*Node

func NewNode(i int) *Node {
	return &Node{ i, WHITE, -1, nil,  make([]*Node, 0) }
}

func dfs1(graph []*Node, queue *([]*Node)) {
	for _, v := range graph {
		if v.mark == WHITE { visit1(v, queue) }
	}
}

func visit1(v *Node, queue *([]*Node)) {
	v.mark = GRAY
	*queue = append(*queue, v)
	for _, u := range(v.adj) {
		if u.mark == WHITE {
			u.parent = v
			visit1(u, queue)
		}
	}
	v.mark = BLACK
}

func dfs2(graph []*Node, queue *([]*Node)) {
	component := 0
	for len(*queue) != 0 {
		v := (*queue)[0]
		*queue = (*queue)[1:]
		if v.comp == -1 {
			visit2(v, component)
			component++
		}
	}
}

func visit2(v *Node, component int) {
	v.comp = component
	for _, u := range(v.adj) {
		if u.comp == -1 && u.parent != v { visit2(u, component) }
	}
}

func main() {
	var n, m int
	fmt.Scanf("%d", &n)
	fmt.Scanf("%d", &m)
	var u, v int
	graph := make([]*Node, n)
	edges := make([]Edge, m)
	for i := 0; i < n; i++ { graph[i] = NewNode(i) }
	for i := 0; i < m; i++ {
		fmt.Scanf("%d %d\n", &u, &v)
		graph[u].adj = append(graph[u].adj, graph[v])
		graph[v].adj = append(graph[v].adj, graph[u])
		edges = append(edges, Edge{u, v})
	}
	queue := make([]*Node, 0)
	dfs1(graph, &queue)
	dfs2(graph, &queue)
	var bridgenum int
	for _, e := range edges {
		if graph[e.u].comp != graph[e.v].comp { bridgenum++ }
	}
	fmt.Println(bridgenum)
}
