package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Arrow struct {
	u     *Node
	color int
}

type Node struct {
	visited bool
	wave    int
	adj     map[int]*Arrow
}

func bfs1(start *Node) {
	queue := []*Node{start}
	start.visited = true
	for len(queue) != 0 {
		v := queue[0]
		queue = queue[1:]
		for _, arrow := range v.adj {
			if u := arrow.u; !u.visited {
				u.wave = v.wave + 1
				u.visited = true
				queue = append(queue, u)
			}
		}
	}
}

func bfs2(start *Node) []int {
	wave := start.wave
	queue := []*Node{start}
	pool := make([]*Arrow, 0)
	colors := make([]int, 0, start.wave)
	minColor := 110
	for len(queue) != 0 {
		v := queue[0]
		queue = queue[1:]
		if v.wave != wave {
			wave = v.wave
			pool = make([]*Arrow, 0)
			colors = append(colors, minColor)
			minColor = 110
		}
		for _, arrow := range v.adj {
			if arrow.u.wave + 1 == wave {
				if minColor > arrow.color {
					minColor = arrow.color
					pool = make([]*Arrow, 0)
				}
				pool = append(pool, arrow)
			}
		}
		if len(queue) == 0 {
			for _, arrow := range pool {
				if arrow.color == minColor {
					queue = append(queue, arrow.u)
				}
			}
		}
	}
	return colors
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
	dict := make(map[int]map[int]int)
	for i := 0; i < n; i++ {
		graph[i] = &Node{false, 0, make(map[int]*Arrow, 0)}
		dict[i] = make(map[int]int)
	}
	for i := 0; i < m; i++ {
		a, b, c := scanInt(scanner) - 1, scanInt(scanner) - 1, scanInt(scanner)
		if oldColor, ok := dict[a][b]; !ok {
			dict[a][b], dict[b][a] = c, c
			graph[a].adj[b] = &Arrow{graph[b], c}
			graph[b].adj[a] = &Arrow{graph[a], c}
		} else if oldColor > c {
			graph[a].adj[b].color = c
			graph[b].adj[a].color = c
		}
	}

	bfs1(graph[n - 1])

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, graph[0].wave)
	answer := bfs2(graph[0])
	for _, color := range answer {
		fmt.Fprintf(writer, "%d ", color)
	}
	writer.Flush()
}
