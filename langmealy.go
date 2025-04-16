package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func dfs(delta [][]int, phi[][]string, q, m, nel int, dict map[int]map[string]bool, word string) {
	if dict[q][word] { return }
	dict[q][word] = true
	if nel == m { return }

	for i, suffix := 0, ""; i < 2; i++ {
		if phi[q][i] != "-" {
			suffix = phi[q][i]
			dfs(delta, phi, delta[q][i], m, nel + 1, dict, word + suffix)
		} else {
			dfs(delta, phi, delta[q][i], m, nel, dict, word)
		}
	}
}

func scanInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	return n
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	n := scanInt(scanner)
	delta := make([][]int, n)
	for i := 0; i < n; i++ {
		delta[i] = make([]int, 2)
		delta[i][0] = scanInt(scanner)
		delta[i][1] = scanInt(scanner)
	}
	phi := make([][]string, n)
	for i := 0; i < n; i++ {
		phi[i] = make([]string, 2)
		scanner.Scan()
		phi[i][0] = scanner.Text()
		scanner.Scan()
		phi[i][1] = scanner.Text()
	}
	q0, m := scanInt(scanner), scanInt(scanner)

	dict := make(map[int]map[string]bool, n)
	for i := 0; i < n; i++ {
		dict[i] = make(map[string]bool)
	}
	dfs(delta, phi, q0, m, 0, dict, "")

	words := make(map[string]bool)
	for _, q := range dict {
		for word, _ := range q {
			words[word] = true
		}
	}
	answer := make([]string, 0, len(words))
	for word, _ := range words {
		answer = append(answer, word)
	}
	sort.Strings(answer)

	writer := bufio.NewWriter(os.Stdout)
	for _, word := range answer {
		fmt.Fprintf(writer, "%s ", word)
	}
	writer.Flush()
}
