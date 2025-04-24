package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
)

type Edge struct {
	a int
	u *Node
}


type Node struct {
	index,
	key int
	value *Node
	adj   []*Edge
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { return pq[i].key < pq[j].key }

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	node := x.(*Node)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -2
	*pq = old[:n-1]
	return node
}

func (pq *PriorityQueue) update(node *Node, value *Node, key int) {
	node.value = value
	node.key = key
	heap.Fix(pq, node.index)
}

func prim(graph []*Node) int {
	pq := make(PriorityQueue, 1)
	v := graph[0]
	pq[0] = v
	heap.Init(&pq)
	dist := 0
	for {
		for _, edge := range v.adj {
			if edge.u.index == -1 {
				edge.u.key = edge.a
				edge.u.value = v
				heap.Push(&pq, edge.u)
			} else if edge.u.index != -2 && edge.a < edge.u.key {
				pq.update(edge.u, v, edge.a)
			}
		}
		if pq.Len() == 0 {
			break
		}
		v = heap.Pop(&pq).(*Node)
		dist += v.key
	}
	return dist
}

func scanInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	return n
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	n, m := scanInt(scanner), scanInt(scanner)
	graph := make([]*Node, n)
	for i := 0; i < n; i++ {
		graph[i] = &Node{-1, 0, nil, make([]*Edge, 0)}
	}
	for i := 0; i < m; i++ {
		u, v, nel := scanInt(scanner), scanInt(scanner), scanInt(scanner)
		graph[u].adj = append(graph[u].adj, &Edge{nel, graph[v]})
		graph[v].adj = append(graph[v].adj, &Edge{nel, graph[u]})
	}
	dist := prim(graph)

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, dist)
	writer.Flush()
}
