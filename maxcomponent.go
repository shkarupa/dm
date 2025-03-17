package main

import "fmt"

type Node struct {
	val,		
	mark	int
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

func NewNode(i int) *Node {
	return &Node{ i, WHITE,  make([]*Node, 0) }
}

func findMaxComp(graph []*Node) (maxComp *Node) {
	var nodes, edges, maxCompNodes, maxCompEdges int
	for _, v := range(graph) {
		nodes, edges = 0, 0
		if v.mark == WHITE {
			visit(v, &nodes, &edges)
		}
		if nodes > maxCompNodes || nodes == maxCompNodes && edges > maxCompEdges {
			maxCompNodes = nodes
			maxCompEdges = edges
			maxComp = v
		}
	}
	return
}

func visit(v *Node, nodes, edges *int) {
	v.mark = GRAY
	for _, u := range(v.adj) {
		if u.mark == WHITE {
			visit(u, nodes, edges)
		}
	}
	*nodes  += 1
	*edges += len(v.adj)
	v.mark = BLACK
}

func markMaxComp(graph []*Node, maxCompRoot *Node) {
	for _, v := range(graph) { v.mark = WHITE }
	mark(maxCompRoot)
}

func mark(v *Node) {
	v.mark = GRAY
	for _, u := range(v.adj) {
		if u.mark == WHITE { mark(u) }
	}
	v.mark = BLACK
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
		edges[i] = Edge{u, v}
	}
	markMaxComp(graph, findMaxComp(graph))
	fmt.Printf("graph {\n")
	for i, v := range graph {
		if v.mark == BLACK {
			fmt.Printf("\t%d [color=red]\n", i)
		} else {
			fmt.Printf("\t%d\n", i)
		}
	}
	for _, edge := range(edges) {
		if graph[edge.v].mark == BLACK {
			fmt.Printf("\t%d--%d [color=red]\n", edge.v, edge.u)
		} else {
			fmt.Printf("\t%d--%d\n", edge.v, edge.u)
		}
	}
	fmt.Printf("}\n")
}
