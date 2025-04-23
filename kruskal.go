package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

type DSU struct {
	i,
	depth int
	parent *DSU
}

func makeSet(i int) *DSU {
	t := &DSU{}
	t.i, t.depth, t.parent = i, 0, t
	return t
}

func find(x *DSU) *DSU {
	if x.parent == x {
		return x
	}
	x.parent = find(x.parent)
	return x.parent
}

func union(x, y *DSU) {
	rx, ry := find(x), find(y)
	if rx.depth < ry.depth {
		rx.parent = ry
	} else {
		ry.parent = rx
		if rx.depth == ry.depth && rx != ry {
			rx.depth++
		}
	}
}

type Edge struct {
	u, v int
	dist float64
}
type ByDist []Edge

func (a ByDist) Len() int           { return len(a) }
func (a ByDist) Less(i, j int) bool { return a[i].dist < a[j].dist }
func (a ByDist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func spanningTree(n int, edges []Edge) []Edge {
	sets := make([]*DSU, n)
	for i := 0; i < n; i++ {
		sets[i] = makeSet(i)
	}
	tree := make([]Edge, 0, n-1)
	for i := 0; i < len(edges) && len(tree) < n-1; i++ {
		if find(sets[edges[i].u]) != find(sets[edges[i].v]) {
			tree = append(tree, edges[i])
			union(sets[edges[i].u], sets[edges[i].v])
		}
	}
	return tree
}

func scanInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	return n
}

type Attraction struct {
	x, y int
}

func dist(a, b Attraction) float64 {
	return math.Sqrt(math.Pow(float64(a.x - b.x), 2) + math.Pow(float64(a.y - b.y), 2))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	n := scanInt(scanner)
	edges := make([]Edge, n*n)
	park := make([]Attraction, n)
	for i := 0; i < n; i++ {
		park[i] = Attraction{scanInt(scanner), scanInt(scanner)}
		for j := 0; j < i; j++ {
			edges[i*n+j] = Edge{j, i, dist(park[i], park[j])}
		}
	}

	sort.Sort(ByDist(edges))
	st := spanningTree(n, edges)
	var length float64
	for _, edge := range st {
		length += edge.dist
	}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(writer, "%.2f\n", length)
	writer.Flush()
}
