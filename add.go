package main

import "fmt"
import "math/rand"
import "math"

func minmax(a, b int) (min, max int) {
	if a < b {
		min, max = a, b;
	} else {
		max, min = a, b;
	}
	return;
}

func add(a, b []int32, p int) []int32 {
	var (
		minlen, maxlen int = minmax(len(a), len(b));
		radix, carryout, sum int32 = (int32)(p), 0, 0;
		maxarray *[]int32;
		result []int32;
	)
	result = make([]int32, 0, maxlen);
	for i := 0; i < minlen; i++ {
		sum = a[i] + b[i] + carryout;
		result = append(result, sum  % radix);
		carryout = sum / radix;
	}
	if maxlen == len(a) {
		maxarray = &a;
	} else {
		maxarray = &b;
	}
	for i := minlen; i < maxlen; i++ {
		sum = (*maxarray)[i] + carryout
		result = append(result, sum % radix);
		carryout = sum / radix;
	}
	if carryout == 1 {
		result = append(result, carryout);
	}
	return result;
}

func digits(p, radix int) []int32 {
	var slice []int32;
	for i := 0; i < p; i++ { slice = append(slice, rand.Int31n((int32)(radix))) };
	return slice;
}

func main() {
	var (
		p int = rand.Intn((int)(rand.Intn((int)(math.Pow(2,30)))));
		a, b []int32;
	)
	if p < 2 {
		p = 2;
	}
	a, b = digits(rand.Intn(100), p), digits(rand.Intn(100), p);
	fmt.Println(p);
	fmt.Println(a);
	fmt.Println(b);
	fmt.Println(add(a, b, p))
}
