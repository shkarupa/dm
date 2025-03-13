package main

import "fmt"

func factorize(n int) map[int]int {
	factorization := make(map[int]int, 1)
	for ; n % 2 == 0; n /= 2 { factorization[2]++ }
	for i := 3; i * i <= n; i += 2 {
		for ; n % i == 0; n /= i { factorization[i]++ }
	}
	if n > 2 { factorization[n] = 1 }
	return factorization
}

func node(factorization, indices map[int]int,
		  processed map[int]bool,
		  powers []int,
		  n int) {
	for factor, maxPower := range factorization {
		if powers[indices[factor]] < maxPower {
			nextN := n * factor
			fmt.Printf("\t%d--%d\n", n, nextN)
			if (! processed[nextN]) {
				nextNPowers := make([]int, len(powers))
				copy(nextNPowers, powers)
				nextNPowers[indices[factor]]++
				node(factorization, indices, processed, nextNPowers, nextN)
			}
		}
	}
	processed[n] = true
}

func main() {
	var n int
	fmt.Scanf("%d", &n)
	fmt.Printf("graph {\n")
	fmt.Printf("\t1\n")
	factorization := factorize(n)
	indices := make(map[int]int, len(factorization))
	i := 0
	for k, _ := range factorization {
		indices[k] = i
		i++
	}
	processed := make(map[int]bool)
	powers := make([]int, len(factorization))
	node(factorization, indices, processed, powers, 1)
	fmt.Printf("}\n")
}
