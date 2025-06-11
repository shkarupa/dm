package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	BLACK = iota
	RED
	BLUE
)

type Node struct {
	ident     string
	time,
	dist,
	color     int
	parents,
	incoming,
	outcoming []*Node
}

func replaceParen(r rune) rune {
	if r == '(' || r == ')' {
		return ' '
	}
	return r
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

func tarjan(graph []*Node) {
	in := make(map[string]int)
	low := make(map[string]int)
	comp := make(map[string]int)
	stack := Stack(make([]*Node, 0))
	time, count := 1, 1

	var visit func(v *Node)
	visit = func(v *Node) {
		in[v.ident], low[v.ident] = time, time
		time++
		stack.Push(v)
		for _, u := range v.outcoming {
			if in[u.ident] == 0 {
				visit(u)
			}
			if comp[u.ident] == 0 && low[v.ident] > low[u.ident] {
				low[v.ident] = low[u.ident]
			}
		}
		if in[v.ident] == low[v.ident] {
			counter := 1
			for {
				u := stack.Pop()
				comp[u.ident] = count
				if u.ident == v.ident {
					if counter > 1 {
						u.color = BLUE
					}
					break
				}
				counter++
				u.color = BLUE
			}
			count++
		}
	}

	for _, v := range graph {
		if in[v.ident] == 0 {
			visit(v)
		}
	}
}

func order(graph []*Node) []*Node {
	visited := make(map[string]bool, len(graph))
	ordered := make([]*Node, 0)
	var visit func(*Node)
	visit = func(v *Node) {
		visited[v.ident] = true
		for _, u := range v.outcoming {
			if !visited[u.ident] && u.color != BLUE {
				visit(u)
			}
		}
		ordered = append(ordered, v)
	}

	for _, v := range graph {
		if !visited[v.ident] && v.color != BLUE {
			visit(v)
		}
	}

	for i, j := 0, len(ordered) - 1; i < j; i, j = i + 1, j - 1 {
		ordered[i], ordered[j] = ordered[j], ordered[i]
	}
	return ordered
}

func maximize(graph []*Node) {
	for _, v := range graph {
		if len(v.incoming) == 0 {
			v.dist = v.time
			continue
		}
		max := 0
		for _, u := range v.incoming {
			if max < u.dist + v.time {
				max = u.dist + v.time
			}
		}
		v.dist = max
		for _, u := range v.incoming {
			if max == u.dist + v.time {
				v.parents = append(v.parents, u)
			}
		}
	}
}

func dyeBlue(graph []*Node) {
	visited := make(map[string]bool, len(graph))
	var visit func(*Node)
	visit = func(v *Node) {
		visited[v.ident] = true
		v.color = BLUE
		for _, u := range v.outcoming {
			if !visited[u.ident] {
				visit(u)
			}
		}
	}

	for _, v := range graph {
		if !visited[v.ident] && v.color == BLUE {
			visit(v)
		}
	}
}

func dyeRed(v *Node) {
	v.color = RED
	for _, u := range v.parents {
		dyeRed(u)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	var sb strings.Builder
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
	}

	jobs := make(map[string]*Node)
	sentences := strings.Split(sb.String(), ";")
	graph := make([]*Node, 0)
	for _, sentence := range sentences {
		current := &Node{"", 0, 0, BLACK, make([]*Node, 0), make([]*Node, 0), make([]*Node, 0)}
		for _, job := range strings.Split(sentence, "<") {
			job = strings.Map(replaceParen, job)
			tokens := strings.Fields(job)
			if len(tokens) == 0 {
				continue
			}
			ident := tokens[0]
			if len(tokens) == 2 {
				time, _ := strconv.Atoi(tokens[1])
				job := &Node{ident, time, 0, BLACK, make([]*Node, 0), make([]*Node, 0), make([]*Node, 0)}
				jobs[ident] = job
				graph = append(graph, job)
			}
			if len(current.ident) != 0 {
				if current == jobs[ident] {
					current.color = BLUE
				}
				jobs[ident].incoming = append(jobs[ident].incoming, current)
				current.outcoming = append(current.outcoming, jobs[ident])
			}
			current = jobs[ident]
		}
	}

	tarjan(graph)
	dyeBlue(graph)
	ordered := order(graph)
	maximize(ordered)

	max := 0
	for _, job := range ordered {
		if max < job.dist {
			max = job.dist
		}
	}
	for _, job := range ordered {
		if job.dist == max {
			dyeRed(job)
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, "digraph {")
	for _, job := range graph {
		fmt.Fprintf(writer, "\t%s [label = \"%s(%d)\"", job.ident, job.ident, job.time)
		if job.color == RED {
			fmt.Fprintf(writer, ", color = red")
		} else if job.color == BLUE {
			fmt.Fprintf(writer, ", color = blue")
		}
		fmt.Fprintf(writer, "]\n")
	}
	for _, job := range graph {
		for _, subroutine := range job.outcoming {
			fmt.Fprintf(writer, "\t%s -> %s", job.ident, subroutine.ident)
			if job.color == subroutine.color && job.color != BLACK {
				if job.color == RED {
					if subroutine.dist == subroutine.time + job.dist {
						fmt.Fprintf(writer, "[color = red]")
					}
				} else {
					fmt.Fprintf(writer, "[color = blue]")
				}
			}
			fmt.Fprintf(writer, "\n")
		}
	}
	fmt.Fprintln(writer, "}")
	writer.Flush()
}
