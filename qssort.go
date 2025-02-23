

package main

import "fmt"

type Task struct {
	low, high int
}

type Stack struct {
	data []Task
}


func (s *Stack) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Stack) Push(t Task) {
	s.data = append(s.data, t) 
}

func (s *Stack) Pop() (t Task) {
	t = s.data[len(s.data) - 1] 
	s.data = s.data[:len(s.data) - 1]
	return t
}

func partition(low, high int, less func(i, j int) bool, swap func(i, j int)) int {
	var i, j int = low, low
	for (j < high) {
		if (less(j, high)) {
			swap(i, j)
			i++
		}
		j++
	}
	swap(i, high)
	return i
}

func qssort(n int, less func(i, j int) bool, swap func(i, j int)) {
	var (
		s Stack = Stack{data: make([]Task, 0)}
		left, right, wen Task
		q int
	)
	s.Push(Task{0, n - 1})
	for (! s.IsEmpty()) {
		wen = s.Pop()
		if (wen.low < wen.high) {
			q = partition(wen.low, wen.high, less, swap)
			left.low = wen.low
			left.high = q - 1
			right.low = q + 1
			right.high = wen.high
			s.Push(right)
			s.Push(left)
		}
	}
}

var seq []int = []int{30, 5, 8, 1, 12, 4, 5, 6, 18, 3, 3, 3}

func less(i, j int) bool {
	return seq[i] < seq[j]
}

func swap(i, j int) {
	seq[i], seq[j] = seq[j], seq[i]
}

func main() {
	fmt.Println(seq)
	qssort(len(seq), less, swap)
	fmt.Println(seq)
}
