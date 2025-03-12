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

func node(factorization, indices map[int]int, seen map[int]bool, powers []int, n int) {
	seen[n] = true
	for k, v := range factorization {
		if powers[indices[k]] < v {
			fmt.Printf("\t%d--%d\n", n, n * k)
			p := make([]int, len(powers))
			copy(p, powers)
			p[indices[k]]++
			if (! seen[n * k]) {
				node(factorization, indices, seen, p, n * k)
			}
		}
	}
}

func main() {
	var n int
	fmt.Scanf("%d", &n)
	fmt.Printf("graph {\n")
	factorization := factorize(n)
	indices := make(map[int]int, len(factorization))
	i := 0
	for k, _ := range factorization {
		indices[k] = i
		i++
	}
	powers := make([]int, len(factorization))
	seen := make(map[int]bool)
	node(factorization, indices, seen, powers, 1)
	fmt.Printf("}\n")
}
