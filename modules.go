package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Tag int

type Lexem struct {
	Tag
	Image string
}

const (
	ERROR Tag = 1 << iota  // Неправильная лексема
	COMMA                  // Запятая
	NUMBER                 // Целое число
	IDENT                  // Имя переменной
	PLUS                   // Знак +
	MINUS                  // Знак -
	MUL                    // Знак *
	DIV                    // Знак /
	EQUAL                  // Знак =
	LT                     // <
	GT                     // >
	LE                     // <=
	GE                     // >=
	NE                     // <>
	LPAREN                 // Левая круглая скобка
	RPAREN                 // Правая круглая скобка
	SEMICOLON              // Точка с запятой
	QUESTION               // Вопросительный знак
	COLON                  // Двоеточие
	WALRUS                 // Моржовый оператор
	EOF
)

var Lexems = map[Tag]string {
	ERROR : "ERROR",
	COMMA : "COMMA",
	NUMBER : "NUMBER",
	IDENT : "IDENT",
	PLUS : "PLUS",
	MINUS : "MINUS",
	MUL : "MUL",
	DIV : "DIV",
	EQUAL : "EQUAL",
	LT : "LT",
	GT : "GT",
	GE : "GE",
	NE : "NE",
	LPAREN : "LPAREN",
	RPAREN : "RPAREN",
	SEMICOLON : "SEMICOLON",
	QUESTION : "QUESTION",
	COLON : "COLON",
	WALRUS : "WALRUS",
	EOF: "EOF",
}

func (l Lexem) String() string {
	return fmt.Sprintf("{%s %s}", Lexems[l.Tag], l.Image)
}


type Lexer struct {
	start,
	pos,
	width   int
	lexems  chan Lexem
	state   StateFn
 	input   string
}

type StateFn func(l *Lexer) StateFn

func (l *Lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.lexems)
}

func lex(input string) *Lexer {
	l := &Lexer{
		input: input,
		state: lexText,
		lexems: make(chan Lexem),
	}
	go l.run()
	return l
}

func (l *Lexer) NextLexem() (Lexem, bool) {
	lexem, ok := <-l.lexems
	return lexem, ok
}

func (l *Lexer) Emit(t Tag) {
	l.lexems <- Lexem{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) Next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return -1
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
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

func (l *Lexer) Errorf(format string, args ...interface{}) StateFn {
	l.lexems <- Lexem{
		ERROR,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func lexNumber(l *Lexer) StateFn {
	for isNumeric(l.Next()) {
	}
	l.Backup()
	if isAlpha(l.Peek()) {
		l.Next()
		return l.Errorf("Bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.Emit(NUMBER)
	return lexText
}

func lexIdentifier(l *Lexer) StateFn {
	for isAlphaNumeric(l.Next()) {
	}
	l.Backup()
	l.Emit(IDENT)
	return lexText
}

func lexText(l *Lexer) StateFn {
	for {
		switch r := l.Next(); {
		case r == -1:
			l.Emit(EOF)
			return nil
		case unicode.IsSpace(r):
			l.Ignore()
		case r == '(':
			l.Emit(LPAREN)
		case r == ',':
			l.Emit(COMMA)
		case r == ')':
			l.Emit(RPAREN)
		case r == '=':
			l.Emit(EQUAL)
		case r == '<':
			if l.Peek() == '>' {
				l.Next()
				l.Emit(NE)
			} else if l.Peek() == '=' {
				l.Next()
				l.Emit(LE)
			} else {
				l.Emit(LT)
			}
		case r == '>':
			if l.Peek() == '=' {
				l.Next()
				l.Emit(GE)
			} else {
				l.Emit(GT)
			}
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
		case r == ':':
			if l.Peek() == '=' {
				l.Next()
				l.Emit(WALRUS)
			} else {
				l.Emit(COLON)
			}
		case r == '?': 
			l.Emit(QUESTION)
		case r == ';':
			l.Emit(SEMICOLON)
		case unicode.IsDigit(r):
			l.Backup()
			return lexNumber
		case isAlphaNumeric(r):
			l.Backup()
			return lexIdentifier
		default:
			return l.Errorf("Unexpexced character: %q", l.input[l.start:l.pos])
		}
	}
}

func isOperation(r rune) bool {
	return 	r == '+' ||
			r == '-' ||
			r == '*' ||
			r == '/'
}

func isAlpha(r rune) bool {
    return 	r >= 'a' && r <= 'z' ||
			r >= 'A' && r <= 'Z'
}

func isNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isNumeric(r)
}

type Function struct {
	ident           string
	dependencies    []string
	formalArgs      map[string]bool
	actualArgsCount []*int
	vars            map[string]bool
}

type Parser struct {
	Lexem
	l          *Lexer
	ident      string
	defined    map[string]*Function
	definition *Function
	functions  []*Function
}

func parser(l *Lexer) *Parser {
	p := &Parser{
		l: l,
		defined: make(map[string]*Function),
		functions: make([]*Function, 0),
	}
	return p
}

func (p *Parser) run() ([]*Function, map[string]*Function, bool) {
	defer func() {
		x := recover(); if x != nil {
			p.functions, p.defined = nil, nil
		}
	}()
	p.Next()
	p.Program()
	return p.functions, p.defined, p.Lexem.Tag & EOF != 0
}

func (p *Parser) Next() {
	lexem, ok := p.l.NextLexem()
	if ok {
		p.Lexem = lexem
	} else {
		p.Lexem = Lexem{EOF, ""}
	}
}

// <program> ::= <function> <program> | <empty>
func (p *Parser) Program() {
	if p.Lexem.Tag & IDENT != 0 {
		p.definition = &Function{p.Lexem.Image, make([]string, 0),
		                         make(map[string]bool),
								 make([]*int, 0), make(map[string]bool)}
		p.defined[p.Lexem.Image] = p.definition
		p.functions = append(p.functions, p.definition)
		p.Function()
		p.Program()
	}
}

// <function> ::= <ident> LPAREN <formal-args-list> RPAREN := <expr> SEMICOLON
func (p *Parser) Function() {
	p.Ident()
	p.Lparen()
	p.FormalArgsList()
	p.Rparen()
	p.Walrus()
	p.Expr()
	p.Semicolon()
}

func (p *Parser) Ident() {
	if p.Lexem.Tag & IDENT != 0 {
		p.ident = p.Lexem.Image
		p.definition.vars[p.Lexem.Image] = true
		p.Next()
	} else {
		panic(fmt.Sprintf("Expected IDENT but got %v", p.Lexem))
	}
}

func (p *Parser) Lparen() {
	if p.Lexem.Tag & LPAREN != 0 {
		p.Next()
	} else {
		panic(fmt.Sprintf("Expected LPAREN but got %v", p.Lexem))
	}
}

func (p *Parser) Rparen() {
	if p.Lexem.Tag & RPAREN != 0 {
		p.Next()
	} else {
		panic(fmt.Sprintf("Expected RPAREN but got %v", p.Lexem))
	}
}

func (p *Parser) Walrus() {
	if p.Lexem.Tag & WALRUS != 0 {
		p.Next()
	} else {
		panic(fmt.Sprintf("Expected WALRUS but got %v", p.Lexem))
	}
}

// <expr> ::= <comparison-expr> <expr-tail>
func (p *Parser) Expr() {
	if p.Lexem.Tag & (NUMBER | IDENT | LPAREN | MINUS) != 0 {
		p.ComparisonExpr()
		p.ExprTail()
	} else {
		panic(fmt.Sprintf("Expected NUMBER, IDENT, LPAREN or MINUS but got %v", p.Lexem))
	}
}

// <expr-tail> ::= QUESTION <comparison-expr> COLON <expr> | <empty>
func (p *Parser) ExprTail() {
	if p.Lexem.Tag & QUESTION != 0 {
		p.Next()
		p.ComparisonExpr()
		p.Colon()
		p.Expr()
	}
}

// <comparison-expr> ::= <arith-expr> <comparison-expr-tail>
func (p *Parser) ComparisonExpr() {
	if p.Lexem.Tag & (NUMBER | IDENT | LPAREN | MINUS) != 0 {
		p.ArithExpr()
		p.ComparisonExprTail()
	} else {
		panic(fmt.Sprintf("Expected NUMBER, IDENT, LPAREN or MINUS but got %v", p.Lexem))
	}
}

// <comparison-expr-tail> ::= <comparison-op> <arith-expr> | <empty>
func (p *Parser) ComparisonExprTail() {
	if p.Lexem.Tag & (EQUAL | NE | LT | GT | LE | GE) != 0 {
		p.Next()
		p.ArithExpr()
	}
}

// <arith-expr> ::= <term> <arith-expr-tail>
func (p *Parser) ArithExpr() {
	if p.Lexem.Tag & (NUMBER | IDENT | LPAREN | MINUS) != 0 {
		p.Term()
		p.ArithExprTail()
	} else {
		panic(fmt.Sprintf("Expected NUMBER, IDENT, LPAREN or MINUS but got %v", p.Lexem))
	}
}

// <arith-expr-tail> ::= PLUS <term> <arith-expr-tail> | MINUS <term> <arith-expr-tail> | <empty>
func (p *Parser) ArithExprTail() {
	if p.Lexem.Tag & (PLUS | MINUS) != 0 {
		p.Next()
		p.Term()
		p.ArithExprTail()
	}
}

// <term> ::= <factor> <term-tail>
func (p *Parser) Term() {
	if p.Lexem.Tag & (NUMBER | IDENT | LPAREN | MINUS) != 0 {
		p.Factor()
		p.TermTail()
	} else {
		panic(fmt.Sprintf("Expected NUMBER, IDENT, LPAREN or MINUS but got %v", p.Lexem))
	}
}

// <term-tail> ::= MUL <factor> <term-tail> | DIV <factor> <term-tail> | <empty>
func (p *Parser) TermTail() {
	if p.Lexem.Tag & (MUL | DIV) != 0 {
		p.Next()
		p.Factor()
		p.TermTail()
	}
}

// <factor> ::= <number> | <ident> <factor-tail> | LPAREN <expr> RPAREN | MINUS <factor>
func (p *Parser) Factor() {
	if p.Lexem.Tag & NUMBER != 0 {
		p.Next()
	} else if p.Lexem.Tag & IDENT != 0 {
		p.Ident()
		p.FactorTail()
	} else if p.Lexem.Tag & LPAREN != 0 {
		p.Lparen()
		p.Expr()
		p.Rparen()
	} else if p.Lexem.Tag & MINUS != 0 {
		p.Next()
		p.Factor()
	} else {
		panic(fmt.Sprintf("Expected NUMBER, IDENT, LPAREN or MINUS but got %v", p.Lexem))
	}
}

// <factor-tail> ::= LPAREN <actual-args-list> RPAREN | <empty>
func (p *Parser) FactorTail() {
	if p.Lexem.Tag & LPAREN != 0 {
		dependency := p.ident
		p.definition.dependencies = append(p.definition.dependencies, dependency)
		p.Lparen()
		i := 0
		count := &i
		p.definition.actualArgsCount = append(p.definition.actualArgsCount, count)
		p.ActualArgsList(count)
		p.Rparen()
	}
}

// <actual-args-tail> ::= <expr-list> | <empty>
func (p *Parser) ActualArgsList(count *int) {
	if p.Lexem.Tag & (NUMBER | IDENT | LPAREN | MINUS) != 0 {
		p.ExprList(count)
	}
}

// <expr-list> ::= <expr> <expr-list-tail>
func (p *Parser) ExprList(count *int) {
	if p.Lexem.Tag & (NUMBER | IDENT | LPAREN | MINUS) != 0 {
		p.Expr()
		*count++
		p.ExprListTail(count)
	}
}

// <expr-list-tail> ::= COMMA <expr> <expr-list-tail> | <empty>
func (p *Parser) ExprListTail(count *int) {
	if p.Lexem.Tag & COMMA != 0 {
		p.Next()
		p.Expr()
		*count++
		p.ExprListTail(count)
	}
}

func (p *Parser) Semicolon() {
	if p.Lexem.Tag & SEMICOLON != 0 {
		p.Next()
	} else {
		panic(fmt.Sprintf("Expected SEMICOLON but got %v", p.Lexem))
	}
}

// <formal-args-list> ::= <ident-list> | <empty>
func (p *Parser) FormalArgsList() {
	if p.Lexem.Tag & IDENT != 0 {
		p.IdentList()
	}
}

// <ident-list> ::= <ident> <ident-list-tail>
func (p *Parser) IdentList() {
	p.definition.formalArgs[p.Lexem.Image] = true
	p.Ident()
	p.IdentListTail()
}

// <ident-list-tail> ::= COMMA <ident> <ident-list-tail> | <empty>
func (p *Parser) IdentListTail() {
	if p.Lexem.Tag & COMMA != 0 {
		p.Next()
		p.definition.formalArgs[p.Lexem.Image] = true
		p.Ident()
		p.IdentListTail()
	}
}

func (p *Parser) Colon() {
	if p.Lexem.Tag & COLON != 0 {
		p.Next()
	} else {
		panic(fmt.Sprintf("Expected COLON but got %v", p.Lexem))
	}
}

type Stack []*Function

func (s *Stack) Push(f *Function) {
	*s = append(*s, f)
}

func (s *Stack) Pop() *Function {
	l := len(*s)
	f := (*s)[l-1]
	*s = (*s)[:l-1]
	return f
}

func tarjan(functions []*Function, defined map[string]*Function) int {
	in := make(map[string]int)
	low := make(map[string]int)
	comp := make(map[string]int)
	for _, f := range functions {
		in[f.ident], low[f.ident], comp[f.ident] = 0, 0, 0
	}
	stack := Stack(make([]*Function, 0))
	time, count := 1, 1

	var visit func(f *Function)
	visit = func(f *Function) {
		in[f.ident], low[f.ident] = time, time
		time++
		stack.Push(f)
		for _, dependency := range f.dependencies {
			if in[dependency] == 0 {
				visit(defined[dependency])
			}
			if comp[dependency] == 0 && low[f.ident] > low[dependency] {
				low[f.ident] = low[dependency]
			}
		}
		if in[f.ident] == low[f.ident] {
			for {
				u := stack.Pop()
				comp[u.ident] = count
				if u.ident == f.ident {
					break
				}
			}
			count++
		}
	}

	for _, f := range functions {
		if in[f.ident] == 0 {
			visit(f)
		}
	}
	return count - 1
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var sb strings.Builder
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
	}
	l := lex(sb.String())
	p := parser(l)
	functions, defined, ok := p.run()
	if !ok {
		fmt.Println("error")
		return
	}
	for _, f := range functions {
		for i, dependency := range f.dependencies {
			if _, isDefined := defined[dependency]; !isDefined {
				fmt.Println("error")
				return
			}
			if len(defined[dependency].formalArgs) != *f.actualArgsCount[i] {
				fmt.Println("error")
				return
			}
		}
		for v, _ := range f.vars {
			if _, isGlobal := defined[v]; !(isGlobal || f.formalArgs[v]) {
				fmt.Println("error")
				return
			}
		}
	}
	fmt.Println(tarjan(functions, defined))
}
