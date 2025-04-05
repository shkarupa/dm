package main

import (
	"container/heap"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const INF = 2 << 32

type Node struct {
	path,
	dist,
	index	int
	adj		[]int
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
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
	node.index = -1
	*pq = old[:n-1]
	return node
}

func relax(u, v *Node) bool {
	if u.dist + v.path < v.dist {
		v.dist = u.dist + v.path
		return true
	}
	return false
}

func dijkstra(graph []*Node) {
	pq := make(PriorityQueue, len(graph))
	copy(pq, graph)
	graph[0].dist = 0
	for pq.Len() != 0 {
		v := heap.Pop(&pq).(*Node)
		for _, neighbour := range v.adj {
			if graph[neighbour].index != -1 && relax(v, graph[neighbour]) {
				heap.Fix(&pq, graph[neighbour].index)
			}
		}
	}

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	graph := make([]*Node, n * n)
	for i := 0; i < n * n; i++ {
		scanner.Scan()
		graph[i] = &Node{ int(scanner.Text()[0] - '0'), INF,  i, make([]int, 0, 4) }
		if i % n != 0 {
			graph[i].adj = append(graph[i].adj, i - 1)
		}
		if (i + 1) % n != 0 {
			graph[i].adj = append(graph[i].adj, i + 1)
		}
		if (i - n) >= 0 {
			graph[i].adj = append(graph[i].adj, i - n)
		}
		if (i + n) < n * n {
			graph[i].adj = append(graph[i].adj, i + n)
		}
	}
	
	dijkstra(graph)
	fmt.Println(graph[0].path + graph[len(graph) - 1].dist)
}
