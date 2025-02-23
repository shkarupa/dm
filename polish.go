package main

import (
	"os"
	"bufio"
	"fmt"
)

type IntSlice []int

type Stack struct {
	IntSlice;
	top int;
}

func (s *Stack) Push(x int) {
	s.IntSlice = append(s.IntSlice, x);
	s.top++;
}

func (s *Stack) Pop() (p int) {
	s.top--;
	p = s.IntSlice[s.top];
	s.IntSlice = s.IntSlice[:s.top];
	return;
}

func (s *Stack) Add() {
	var a, b int = s.Pop(), s.Pop();
	s.Push(a + b);
}

func (s *Stack) Sub() {
	var a, b int = s.Pop(), s.Pop();
	s.Push(a - b);
}

func (s *Stack) Mul() {
	var a, b int = s.Pop(), s.Pop();
	s.Push(a * b);
}

func eval(expression string) int {
	var s Stack;
	s.IntSlice = make(IntSlice, 0);
	for i := len(expression) - 1; i >= 0; i-- {
		if '0' <= expression[i]  && expression[i] <= '9' {
			(&s).Push(int(expression[i] - '0'));
		} else if expression[i] == '+' {
			(&s).Add();
		} else if expression[i] == '-' {
			(&s).Sub();
		} else if expression[i] == '*' {
			(&s).Mul();
		}
	}
	return s.Pop();
}

func main() {
	input := bufio.NewReader(os.Stdin);
	s, _ := input.ReadString('\n');
	fmt.Printf("%d", eval(s));
}
