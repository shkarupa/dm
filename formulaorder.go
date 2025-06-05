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
	VAR                    // Имя переменной
	PLUS                   // Знак +
	MINUS                  // Знак -
	MUL                    // Знак *
	DIV                    // Знак /
	ASSIGN                 // Знак =
	LPAREN                 // Левая круглая скобка
	RPAREN                 // Правая круглая скобка
	NEWLINE                // Перевод строки
	EOF
)

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
	l.Emit(VAR)
	return lexText
}

func lexText(l *Lexer) StateFn {
	for {
		switch r := l.Next(); {
		case r == -1:
			l.Emit(EOF)
			return nil
		case unicode.IsSpace(r) && r != '\n':
			l.Ignore()
		case r == '\n':
			l.Emit(NEWLINE)
		case r == '(':
			l.Emit(LPAREN)
		case r == ',':
			l.Emit(COMMA)
		case r == ')':
			l.Emit(RPAREN)
		case r == '=':
			l.Emit(ASSIGN)
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
			return l.Errorf("Unexpexced character: %q", l.input[l.start:l.pos])
		}
	}
}

type Parser struct {
	Lexem
	id        int
	l         *Lexer
	equasion  *Formula
	equasions []*Formula
	defined   map[string]*Formula
}

type Formula struct {
	id           int
	variables,
	expressions,
	dependencies []string
}

func (f Formula) String() string {
	return fmt.Sprintf("%s = %s", strings.Join(f.variables, " "), strings.Join(f.expressions, " "))
}

func parser(l *Lexer) *Parser {
	p := &Parser{
		l: l,
		equasions: make([]*Formula, 0),
		defined: make(map[string]*Formula),
	}
	return p
}

func (p *Parser) run() ([]*Formula, map[string]*Formula, bool) {
	defer func() {
		x := recover(); if x != nil {
			p.equasions = nil
			p.defined = nil
		}
	}()
	p.Next()
	p.Formulas()
	return p.equasions, p.defined, p.Lexem.Tag & EOF != 0
}

func (p *Parser) Next() {
	lexem, ok := p.l.NextLexem()
	if ok {
		p.Lexem = lexem
	} else {
		p.Lexem = Lexem{EOF, ""}
	}
}

// <formulas> ::= <formula> NEWLINE <formulas> | <empty>
func (p *Parser) Formulas() {
	if p.Lexem.Tag & VAR != 0 {
		p.Formula()
		p.NewLine()
		p.Formulas()
	}
}

// <formula> ::= <variables> = <expressions>
func (p *Parser) Formula() {
	if p.Lexem.Tag & VAR != 0 {
		p.equasion = &Formula{p.id, make([]string, 0), make([]string, 0), make([]string, 0)}
		p.VariableVariables()
		p.FormulaInner()
		p.Expression()
	} else {
		panic(fmt.Sprintf("Expected VAR but got %v", p.Lexem))
	}
}

func (p *Parser) VariableVariables() {
	if p.Lexem.Tag & VAR != 0 {
		if _, isDefined := p.defined[p.Lexem.Image]; isDefined {
			panic(fmt.Sprintf("Variable '%s' is already defined", p.Lexem.Image))
		}
		p.defined[p.Lexem.Image] = p.equasion
		p.equasion.variables = append(p.equasion.variables, p.Lexem.Image)
		p.Next()
	}
}

// <formula-inner> ::= COMMA <variable> <fomula-inner> <expression> COMMA | =
func (p *Parser) FormulaInner() {
	if p.Lexem.Tag & COMMA != 0 {
		p.CommaVariables()
		p.VariableVariables()
		p.FormulaInner()
		p.Expression()
		p.CommaExpressions()
	} else if p.Lexem.Tag & ASSIGN != 0 {
		p.Next()
	} else {
		panic(fmt.Sprintf("Expected COMMA or ASSIGN but got %v", p.Lexem))
	}
}

func (p *Parser) CommaVariables() {
	if p.Lexem.Tag & COMMA != 0 {
		p.equasion.variables = append(p.equasion.variables, ",")
		p.Next()
	} else {
		panic(fmt.Sprintf("Expected COMMA but got %v", p.Lexem))
	}
}

func (p *Parser) CommaExpressions() {
	if p.Lexem.Tag & COMMA != 0 {
		p.equasion.expressions = append(p.equasion.expressions, ",")
		p.Next()
	} else {
		panic(fmt.Sprintf("Expected COMMA but got %v", p.Lexem))
	}
}

// <expression> ::= <term> <expression-tail>
func(p *Parser) Expression() {
	if p.Lexem.Tag & (NUMBER | VAR | LPAREN | MINUS) != 0 {
		p.Term()
		p.ExpressionTail()
	} else {
		panic(fmt.Sprintf("Expected NUMBER, VAR or LPAREN but got %v", p.Lexem))
	}
}

// <term> ::= <factor> <term-tail>
func (p *Parser) Term() {
	p.Factor()
	p.TermTail()
}

// <factor> ::= <number> <expression-tail>   |
//              <variable> <expression-tail> |
//              ( <expression> ) 
//				MINUS <factor>
func (p *Parser) Factor() {
	if p.Lexem.Tag & NUMBER != 0 {
		p.Number()
	} else if p.Lexem.Tag & VAR != 0 {
		p.VariableExpressions()
	} else if p.Lexem.Tag & LPAREN != 0 {
		p.Lparen()
		p.Expression()
		p.Rparen()
	} else if p.Lexem.Tag & MINUS != 0 {
		p.equasion.expressions = append(p.equasion.expressions, p.Lexem.Image)
		p.Next()
		p.Factor()
	} else {
		panic(fmt.Sprintf("Expected NUMBER, VAR or LPAREN but got %v", p.Lexem))
	}
}

func (p *Parser) Number() {
	p.equasion.expressions = append(p.equasion.expressions, p.Lexem.Image)
	p.Next()
}

func (p *Parser) VariableExpressions() {
	p.equasion.expressions = append(p.equasion.expressions, p.Lexem.Image)
	p.equasion.dependencies = append(p.equasion.dependencies, p.Lexem.Image)
	p.Next()
}

func (p *Parser) Lparen() {
	p.equasion.expressions = append(p.equasion.expressions, p.Lexem.Image)
	p.Next()
}

// <expression-tail> ::= + <term> <expression-tail> | - <term> <expression-tail> | <empty>
func (p *Parser) ExpressionTail() {
	if p.Lexem.Tag & (PLUS | MINUS) != 0 {
		p.equasion.expressions = append(p.equasion.expressions, p.Lexem.Image)
		p.Next()
		p.Term()
		p.ExpressionTail()
	} 
}

func (p *Parser) Rparen() {
	if p.Lexem.Tag & RPAREN != 0 {
		p.equasion.expressions = append(p.equasion.expressions, p.Lexem.Image)
		p.Next()
	} else {
		panic(fmt.Sprintf("Expected RPAREN but got %v", p.Lexem))
	}
}

// <term-tail> ::= * <factor> <term-tail> | / <factor> <term-tail> | <empty>
func (p *Parser) TermTail() {
	if p.Lexem.Tag & (MUL | DIV) != 0 {
		p.equasion.expressions = append(p.equasion.expressions, p.Lexem.Image)
		p.Next()
		p.Factor()
		p.TermTail()
	}
}

func (p *Parser) NewLine() {
	if p.Lexem.Tag & NEWLINE != 0 {
		p.equasions = append(p.equasions, p.equasion)
		p.id++
		p.Next()
	} else {
		panic(fmt.Sprintf("Expected NEWLINE but got %v", p.Lexem))
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

const (
	WHITE = iota
	GRAY
	BLACK
)

func dfs(formulas []*Formula, defined map[string]*Formula) ([]int, bool) {
	ids := make([]int, 0)
	colors := make(map[int]int)
	for id := range formulas {
		colors[id] = WHITE
	}
	ok := true
	for _, formula := range formulas {
		if ids, ok = visit(formula, defined, colors, ids); !ok {
			return nil, false
		}
	}
	return ids, true
}

func visit(formula *Formula, defined map[string]*Formula, colors map[int]int, ids []int) ([]int, bool) {
	if colors[formula.id] == BLACK {
		return ids, true
	}
	if colors[formula.id] == GRAY {
		return nil, false
	}
	colors[formula.id] = GRAY
	ok := true
	for _, variable := range formula.dependencies {
		if ids, ok = visit(defined[variable], defined, colors, ids); !ok {
			return nil, false
		}
	}
	colors[formula.id] = BLACK
	return append(ids, formula.id), true
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var sb strings.Builder
	rawFormulas := make([]string, 0)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteByte('\n')
		rawFormulas = append(rawFormulas, scanner.Text())
	}
	l := lex(sb.String())
	p := parser(l)
	formulas, defined, ok := p.run()
	if !ok {
		fmt.Println("syntax error")
		return
	}
	for _, formula := range formulas {
		for _, variable := range formula.dependencies {
			if _, isDefined := defined[variable]; !isDefined {
				fmt.Println("syntax error")
				return
			}
		}
	}
	if ids, hasNoCycle := dfs(formulas, defined); hasNoCycle {
		for _, id := range ids {
			fmt.Println(rawFormulas[id])
		}
	} else {
		fmt.Println("cycle")
	}
}
