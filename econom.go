// <start> ::= LETTER | <expression>
// <expression> ::= OPEN <operation> CLOSE | LETTER
// <operation> ::= OPERATOR <expression> <expression> 

package main

import (
	"fmt"
	"unicode/utf8"
	"unicode"
	"bufio"
	"os"
)

type RPN struct {
	s string
	current int
	dict map[string]int
}

func (e *RPN) Peek() byte {
	return e.s[e.current]
}

func (e *RPN) HasNext() bool {
	return e.current < utf8.RuneCountInString(e.s)
}

func (e *RPN) Parse() int {
	if (! e.HasNext()) {
		return len(e.dict)
	} else {
		if (e.Peek() == '(') {
			e.Expression()
			return e.Parse()
		} else {
			e.Letter()
			return e.Parse()
		}
	}
}

func (e *RPN) Next() {
	e.current++
}

func (e *RPN) Expression() {
	if (unicode.IsLetter(rune(e.Peek()))) {
		e.Letter()
	} else {
		start := e.current
		e.Open()
		e.Operator()
		e.Expression()
		e.Expression()
		end := e.current
		expr := e.s[start:end]
		_, ok := e.dict[expr]
		if (! ok) {
			e.dict[expr]++
		}
		e.Close()
	}
}

func (e *RPN) Letter() {
	e.Next()
}

func (e *RPN) Open() {
	e.Next()
}

func (e *RPN) Close() {
	e.Next()
}

func (e *RPN) Operator() {
	e.Next()
}

func main() {
	var n RPN
	reader := bufio.NewReader(os.Stdin)
	n.current = 0
	n.dict = make(map[string]int)
	n.s, _ = reader.ReadString('\n')
	fmt.Println(n.Parse())
}
