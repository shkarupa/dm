

// http://e-maxx.ru/algo/linear_systems_gauss
// sicp 2.1.1

package main

import (
	"fmt"
)

type Rat struct {
	numer, denom int
}

func addRat(x, y Rat) Rat {
	return NewRat(
		x.numer * y.denom + y.numer * x.denom,
		x.denom * y.denom)
}

func subRat(x, y Rat) Rat {
	return NewRat(
		x.numer * y.denom - y.numer * x.denom,
		x.denom * y.denom)
}

func mulRat(x, y Rat) Rat {
	return NewRat( x.numer * y.numer, x.denom * y.denom )
}

func divRat(x, y Rat) Rat {
	return NewRat( x.numer * y.denom, x.denom * y.numer )
}

func NewRat(x, y int) Rat {
	if x == 0 {
		return Rat{0, 1}
	}
	g := gcd(x, y)
	return Rat{x / g, y / g}
}

func embellish(r Rat) Rat {
	if r.denom < 0 {
		return Rat{-r.numer, -r.denom}
	}
	return r
}

func (r Rat) String() string {
	return fmt.Sprintf("%d/%d", r.numer, r.denom)
}

func absRat(x Rat) Rat {
	return NewRat(abs(x.numer), abs(x.denom))
}

func compareRat(x, y Rat) int {
	return x.numer * y.denom - x.denom * y.numer
}

func abs(x int) int {
	if (x > 0) { return x }
	return -x
}

func gcd(x, y int) int {
	for ; y != 0; x, y = y, x % y {
	}
	return x
}

type Matrix [][]Rat

func (m Matrix) Add(i, j int, k Rat) {
	for index := 0; index < len(m[i]); index++ {
		m[i][index] = addRat(m[i][index], mulRat(m[j][index], k))
	}
}

func (m Matrix) Sub(i, j int, k Rat) {
	for index := 0; index < len(m[i]); index++ {
		m[i][index] = subRat(m[i][index], mulRat(m[j][index], k))
	}
}

func (m Matrix) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (a Matrix) Gauss() (ans []Rat, exists bool) {
	n, m := len(a), len(a[0]) - 1
	where := make([]int, m)
	for i := range where { where[i] = -1 }
	for col, row := 0, 0; col <= m && row < n; col, row = col + 1, row + 1 {
		sel := row
		for i := row; i < n; i++ {
			if compareRat(absRat(a[i][col]), absRat(a[sel][col])) > 0 { sel = i }
		}
		if a[sel][col].numer == 0 { continue }
		a.Swap(sel, row)
		where[col] = row
		for i := 0; i < n; i++ {
			if i != row {
				a.Sub(i, row, divRat(a[i][col], a[row][col]))
			}
		}
	}
	ans = make([]Rat, m)
	for i := 0; i < m; i++ {
		if where[i] != -1 {
			ans[i] = divRat(a[where[i]][m], a[where[i]][i]);
		}
	}
	for i := 0; i < n; i++ {
		sum := NewRat(0, 1)
		for j := 0; j < m; j++ {
			sum = addRat(sum, mulRat(ans[j], a[i][j]))
		}
		if compareRat(subRat(sum, a[i][m]), Rat{0, 1}) != 0 {
			exists = false
			return
		}
	}
	for i := 0; i < m; i++ {
		if where[i] == -1 {
			exists = false
			return
		}
	}
	exists = true
	return 
}

func main() {
	var n, i, j, coeff int;
	fmt.Scan(&n);
	var matrix Matrix = make([][]Rat, n);
	for i = 0; i < n; i++ {
		matrix[i] = make([]Rat, n + 1, n + 1);
		for j= 0; j < n + 1; j++ {
			fmt.Scan(&coeff);
			matrix[i][j] = Rat{coeff, 1}
		}
	}
	ans, exists := matrix.Gauss()
	if exists {
		for _, v := range ans {
			fmt.Println(embellish(v))
		}
	} else {
		fmt.Println("No solution")
	}
}
