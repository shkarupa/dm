package main

import (
	"bufio"
	"os"
	"fmt"
)

func main() {
	output := bufio.NewWriter(os.Stdout)
	var n, m, q0 int
	fmt.Scanf("%d\n%d\n%d\n", &n, &m, &q0)
	delta := make([][]int, n)
	phi := make([][]string, n)
	for i := 0; i < n; i++ {
		delta[i] = make([]int, m)
		for j := 0; j < m; j++ { fmt.Scanf("%d", &delta[i][j]) }
	}
	for i := 0; i < n; i++ {
		phi[i] = make([]string, m)
		for j := 0; j < m; j++ { fmt.Scanf("%s", &phi[i][j]) }
	}
	fmt.Fprintf(output, "digraph {\n")
	fmt.Fprintf(output, "\trankdir = LR\n")
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			fmt.Fprintf(output,
						"\t%d -> %d [label = \"%c(%s)\"]\n",
						i, delta[i][j], 97 + j, phi[i][j])
		}
	}
	fmt.Fprintf(output, "}\n")
	output.Flush()
}
