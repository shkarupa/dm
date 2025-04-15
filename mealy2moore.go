package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func scanInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	return n
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	var k int

	k = scanInt(scanner)
	X := make([]string, k)
	for i := 0; i < k; i++ {
		scanner.Scan()
		X[i] = scanner.Text()
	}

	k = scanInt(scanner)
	Y := make([]string, k)
	for i := 0; i < k; i++ {
		scanner.Scan()
		Y[i] = scanner.Text()
	}

	n := scanInt(scanner)
	delta := make([][]int, n)
	for i := 0; i < n; i++ {
		delta[i] = make([]int, len(X))
		for j := 0; j < len(X); j++ {
			delta[i][j] = scanInt(scanner)
		}
	}
	phi := make([][]int, n)
	for i := 0; i < n; i++ {
		phi[i] = make([]int, len(X))
		for j := 0; j < len(X); j++ {
			phi[i][j] = scanInt(scanner)
		}
	}

	states := make(map[int]map[int]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < len(X); j++ {
			if states[delta[i][j]] == nil {
				states[delta[i][j]] = make(map[int]int)
			}
			states[delta[i][j]][phi[i][j]] = -1
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(writer, "digraph {\n")
	fmt.Fprintf(writer, "\trankdir = LR\n")
	for i, counter := 0, 0; i < len(states); i++ {
		keys := make([]int, 0)
		for k, _ := range states[i] {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			fmt.Fprintf(writer,
			            "\t%d [label = \"(%d, %s)\"]\n",
				    counter, i, Y[k])
			states[i][k] = counter
			counter++
		}
	}
	for i := 0; i < len(states); i++ {
		for j := 0; j < len(states[i]); j++ {
			for k := 0; k < len(X); k++ {
				fmt.Fprintf(writer,
				            "\t%d -> %d [label = \"%s\"]\n",
					    states[i][j], states[delta[i][k]][phi[i][k]], X[k])
			}
		}
	}
	fmt.Fprintf(writer, "}\n")
	writer.Flush()
}
