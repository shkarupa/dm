package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	BLACK = iota
	RED
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
	name    string 
	color,       
	pas, dcu int 
	imports  []string
}

func visit(graph map[string]*Node, source string) {
	unit := graph[source]
	if (unit.dcu == -1) {
		unit.color = RED
	}
	if (unit.dcu < unit.pas) {
		unit.color = RED
	}
	for _, dependency := range unit.imports {
		if dependency == unit.name {
			continue
		}
		if graph[dependency].color == RED || graph[dependency].dcu > unit.dcu {
			unit.color = RED
		}
		visit(graph, dependency)
	}
}

func findCycles(graph map[string]*Node) bool {
	color := make(map[string]int)
	var ok bool
	var visit func(source string)
	visit = func(source string) {
		color[source] = 1
		for _, dependency := range graph[source].imports {
			if color[dependency] == 1 {
				ok = true
			}
			if color[dependency] == 0 {
				visit(dependency)
			}
		}
		color[source] = 2
	}
	visit("main")
	return ok
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	n := scanInt(scanner)
	graph := make(map[string]*Node, n)
	for i := 0; i < n; i++ {
		source := scanString(scanner)
		k := scanInt(scanner)
		imports := make([]string, 0, k)
		for j := 0; j < k; j++ {
			dependency := scanString(scanner)
			imports = append(imports, dependency)
		}
		graph[source] = &Node{ source, BLACK, -1, -1, imports }
	}
	m := scanInt(scanner)
	for i := 0; i < m; i++ {
		fileName := scanString(scanner)
		sourceNextension := strings.Split(fileName, ".")
		source, extension := sourceNextension[0], sourceNextension[1]
		timestamp := scanInt(scanner)
		if (extension == "pas") {
			graph[source].pas = timestamp
		} else {
			graph[source].dcu = timestamp
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	if findCycles(graph) {
		fmt.Fprintln(writer, "!CYCLE")
		writer.Flush()
		return
	}
	for i := 0; i < n; i++ {
		visit(graph, "main")
	}
	answer := make([]string, 0)
	for unit := range graph {
		if graph[unit].color == RED {
			answer = append(answer, unit)
		}
	}
	sort.Strings(answer)
	for i := range answer {
		fmt.Fprintf(writer, "%s.pas\n", answer[i])
	}
	writer.Flush()
}
