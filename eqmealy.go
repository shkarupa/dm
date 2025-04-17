package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func scanInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	return n
}

type Machine struct {
	n, m, q int
	delta   [][]int
	phi     [][]string
}

func (a Machine) Peek(input int) string { return a.phi[a.q][input] }
func (a *Machine) SetState(q int)       { a.q = q }

func (a *Machine) Step(input int) string {
	output := a.Peek(input)
	a.SetState(a.delta[a.q][input])
	return output
}

func run(machine1, machine2 Machine, dict1, dict2 map[int]map[int]bool) bool {
	eq := true
	for i := 0; i < machine1.m && eq; i++ {
		if machine1.Peek(i) != machine2.Peek(i) {
			eq = false
		}
	}
	if eq {
		currentState1, currentState2 := machine1.q, machine2.q
		for i := 0; i < machine1.m; i++ {
			machine1.SetState(currentState1)
			machine2.SetState(currentState2)
			if !dict1[machine1.q][i] || !dict2[machine2.q][i] {
				dict1[machine1.q][i] = true
				dict2[machine2.q][i] = true
				machine1.Step(i)
				machine2.Step(i)
				eq = eq && run(machine1, machine2, dict1, dict2)
			}
		}
	} 
	return eq
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	machines := make([]Machine, 2)
	for i := 0; i < 2; i++ {
		n, m, q0 := scanInt(scanner), scanInt(scanner), scanInt(scanner)
		delta, phi := make([][]int, n), make([][]string, n)
		for i := 0; i < n; i++ {
			delta[i] = make([]int, m)
			for j := 0; j < m; j++ {
				delta[i][j] = scanInt(scanner)
			}
		}
		for i := 0; i < n; i++ {
			phi[i] = make([]string, m)
			for j := 0; j < m; j++ {
				scanner.Scan()
				phi[i][j] = scanner.Text()
			}
		}
		machines[i] = Machine{ n, m, q0, delta, phi }
	}

	writer := bufio.NewWriter(os.Stdout)
	if machines[0].m != machines[1].m {
		fmt.Fprintf(writer, "NOT EQUAL\n")
		writer.Flush()
		return
	}
	dict1 := make(map[int]map[int]bool, machines[0].n)
	dict2 := make(map[int]map[int]bool, machines[1].n)
	for i := 0; i < machines[0].n; i++ {
		dict1[i] = make(map[int]bool)
	}
	for i := 0; i < machines[1].n; i++ {
		dict2[i] = make(map[int]bool)
	}
	if run(machines[0], machines[1], dict1, dict2) {
		fmt.Fprintf(writer, "EQUAL\n")
	} else {
		fmt.Fprintf(writer, "NOT EQUAL\n")
	}
	writer.Flush()
}
