// https://go.dev/talks/2011/lex.slide#1
// https://www.youtube.com/watch?v=HxaD_trXwRE
// https://ruslanspivak.com/lsbasi-part7/

package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
	"strings"
	"strconv"
	"os"
)

type Tag int

type Lexem struct {
	Tag
	Image string
}

const (
    ERROR Tag = 1 << iota  // Неправильная лексема
    NUMBER                 // Целое число
    VAR                    // Имя переменной
    PLUS                   // Знак +
    MINUS                  // Знак -
    MUL                    // Знак *
    DIV                    // Знак /
    LPAREN                 // Левая круглая скобка
    RPAREN                 // Правая круглая скобка
	EOF
)

type StateFn func(l *Lexer) StateFn

type Lexer struct {
	expr 	string
	start 	int
	pos 	int
	width 	int
	lexems 	chan Lexem
}

func (l *Lexer) run() {
	for state:= lexText; state != nil; {
		state= state(l)
	}
	close(l.lexems)
}

// <tokens> ::= <token> <tokens> | <spaces> <tokens> | <empty>
// <spaces> ::= SPACE <spaces> | <empty>
// <token> ::= <id> | <number> | OPERATION | LPAREN | RPAREN
// <id> ::= LETTER <id-tail>
// <id-tail> ::= LETTER <id-tail> | DIGIT <id-tail> | <empty>
// <number> ::= DIGIT <number-tail>
// <number-tail> ::= DIGIT <number-tail> | <empty>

func lexText(l *Lexer) StateFn {
	for {
		switch r := l.Next(); {
		case r == -1:
			return nil
		case unicode.IsSpace(r):
			l.Ignore()
		case r == '(':
			l.Emit(LPAREN)
		case r == ')':
			l.Emit(RPAREN)
		case isOperation(r):
			if r == '+' {
				l.Emit(PLUS)
			} else if r == '-' {
				l.Emit(MINUS)
			} else if r == '*' {
				l.Emit(MUL)
			} else if r == '/' {
				l.Emit(DIV)
			}
		case unicode.IsDigit(r):
			l.Backup()
			return lexNumber
		case isAlphaNumeric(r):
			l.Backup()
			return lexIdentifier
		default:
			l.Emit(ERROR)
		}
	}
}

func lexNumber(l *Lexer) StateFn {
	digits := "0123456789"
	l.AcceptRun(digits)
	l.Emit(NUMBER)
	return lexText
}

func lexIdentifier(l *Lexer) StateFn {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	l.Accept(letters)
	digits := "0123456789"
	l.AcceptRun(letters + digits)
	l.Emit(VAR)
	return lexText
}

func (l *Lexer) Next() (r rune) {
	if l.pos >= len(l.expr) {
		l.width = 0
		return -1
	}
	r, l.width = utf8.DecodeRuneInString(l.expr[l.pos:])
	l.pos += l.width
	return r
}

func (l *Lexer) Ignore() { l.start = l.pos 	}
func (l *Lexer) Backup() { l.pos -= l.width }

func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

func (l *Lexer) Emit(t Tag) {
	l.lexems <- Lexem{t, l.expr[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) Accept(valid string) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		return true
	}
	l.Backup()
	return false
}

func (l *Lexer) AcceptRun(valid string) {
	for strings.IndexRune(valid, l.Next()) >= 0 { }
	l.Backup()
}

func lexer(expr string, lexems chan Lexem) {
	l := &Lexer{
		expr: 	expr,
		lexems: lexems,
	}
	l.run()
}

func isOperation(r rune) bool {
	return 	r == '+' ||
			r == '-' ||
			r == '*' ||
			r == '/'
}

func isAlphaNumeric(r rune) bool {
    return 	(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9')
}

type Parser struct {
	lexems 	[]Lexem
	pos		int
}

func (p *Parser) Next() (n Lexem) {
	if p.pos == len(p.lexems) {
		return Lexem{ EOF, "" }
	} else {
		n = p.lexems[p.pos]
		p.pos++
	}
	return n
}

func (p *Parser) Backup() { p.pos-- }

func (p *Parser) Peek() (n Lexem) {
	n = p.Next()
	p.Backup()
	return
}

type Node struct {
	Lexem
	left, right	*Node
}

// <E>  ::= <T> <E'>.
// <E'> ::= + <T> <E'> | - <T> <E'> | .
// <T>  ::= <F> <T'>.
// <T'> ::= * <F> <T'> | / <F> <T'> | .
// <F>  ::= <number> | <var> | ( <E> ) | - <F>.

func (p *Parser) E() (n *Node) {
	defer func() {
		x := recover(); if x != nil {
			n = nil
			return
		}
	}()
	n = p.T()
	n = p.Et(n)
	return
}

func (p *Parser) T() (n *Node) {
	n = p.F()
	n = p.Tt(n)
	return
}

func (p *Parser) F() (n *Node) {
	if p.Peek().Tag & ERROR != 0 { panic("Error lexem") }
	if lx := p.Peek(); lx.Tag & (NUMBER | VAR | LPAREN | MINUS) != 0 {
		lx = p.Next()
		if lx.Tag & (NUMBER | VAR) != 0 {
			n = &Node{ lx, nil, nil }
		} else if lx.Tag & LPAREN != 0 {
			n = p.E()
			if (p.Peek()).Tag & RPAREN != 0 {
				p.Next()
			} else {
				panic("No closing parentesis")
			}
		} else if lx.Tag & MINUS != 0 {
			n = &Node{ lx,  nil, p.F() }
		}
	} else {
		panic("Not a factor")
	}
	return
}

func (p *Parser) Et(n *Node) *Node {
	if p.Peek().Tag & ERROR != 0 { panic("Error lexem") }
	if lx := p.Peek(); lx.Tag & (PLUS | MINUS) != 0 {
		lx = p.Next()
		m := &Node{ lx, n, p.T() }
		return p.Et(m)
	} else {
		return n
	}
}

func (p *Parser) Tt(n *Node) *Node {
	if p.Peek().Tag & ERROR != 0 { panic("Error lexem") }
	if lx := p.Peek(); lx.Tag & (MUL | DIV) != 0 {
		lx = p.Next()
		m := &Node{ lx, n, p.F() }
		return p.Tt(m)
	} else {
		return n
	}
}

func parser(lexems []Lexem) (*Node) {
	p := &Parser{
		lexems: lexems,
	}
	n := p.E()
	if p.Next().Tag & EOF != 0 { return n }
	return nil
}

func eval(ast *Node, dict map[string]int) (m int) {
	if ast == nil { return 0 }
	switch ast.Lexem.Tag {
	case PLUS:
		m = eval(ast.left, dict) + eval(ast.right, dict)
	case MINUS:
		m = eval(ast.left, dict) - eval(ast.right, dict)
	case MUL:
		m = eval(ast.left, dict) * eval(ast.right, dict)
	case DIV:
		m = eval(ast.left, dict) / eval(ast.right, dict)
	case NUMBER:
		n, _ := strconv.Atoi(ast.Lexem.Image)
		m = n
	case VAR:
		v, ok := dict[ast.Image]
		if ok {
			m = v
		} else {
			fmt.Scanf("%d\n", &v)
			dict[ast.Image] = v
			m = v
		}
	}
	return m
}

func main() {
	lexems := make(chan Lexem)
	stream := make([]Lexem, 0)
	go lexer(os.Args[1], lexems)
	for lexem := range lexems { stream = append(stream, lexem) }
	stream = append(stream, Lexem{ EOF, "" })
	ast := parser(stream)
	if ast == nil {
		fmt.Println("error")
	} else {
		dict := make(map[string]int, 10)
		fmt.Println(eval(ast, dict))
	}
}
