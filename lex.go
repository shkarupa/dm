// https://go.dev/talks/2011/lex.slide#1
// https://www.youtube.com/watch?v=HxaD_trXwRE

package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
	"strings"
	"math/rand"
)

type StateFn func(l *Lexer) StateFn

type Lexer struct {
	sentence 	string
	start 		int
	pos 		int
	width 		int
	seq 		[]int
	array		AssocArray
	num			int
}

func (l *Lexer) run() {
	for state:= lexText; state != nil; {
		state= state(l)
	}
}

// <tokens> ::= <token> <tokens> | <spaces> <tokens> | <empty>
// <spaces> ::= SPACE <spaces> | <empty>
// <token> ::= <id>
// <id> ::= LETTER <id-tail>
// <id-tail> ::= LETTER <id-tail> | DIGIT <id-tail> | <empty>

func lexText(l *Lexer) StateFn {
	for {
		switch r := l.Next(); {
		case r == '$':
			return nil
		case unicode.IsSpace(r):
			l.Ignore()
		case isAlphaNumeric(r):
			l.Backup()
			return lexIdentifier
		}
	}
}

func lexIdentifier(l *Lexer) StateFn {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	l.Accept(letters)
	digits := "0123456789"
	l.AcceptRun(letters + digits)
	l.Emit()
	return lexText
}

func (l *Lexer) Next() (r rune) {
	if l.pos >= len(l.sentence) {
		l.width = 0
		return '$'
	}
	r, l.width = utf8.DecodeRuneInString(l.sentence[l.pos:])
	l.pos += l.width
	return r
}

func (l *Lexer) Ignore() { l.start = l.pos }
func (l *Lexer) Backup() { l.pos -= l.width }

func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

func (l *Lexer) Emit() {
	x, exists := l.array.Lookup(l.sentence[l.start:l.pos])
	if exists {
		l.seq = append(l.seq, x)
	} else {
		l.num++
		l.seq = append(l.seq, l.num)
		l.array.Assign(l.sentence[l.start:l.pos], l.num)
	}
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
	for strings.IndexRune(valid, l.Next()) >= 0 {
	}
	l.Backup()
}

func lex(sentence string, array AssocArray) []int {
	l := &Lexer{
		sentence: 	sentence,
		seq:		make([]int, 0),
		array:		array,
	}
	l.run()
	return l.seq
}

func isAlphaNumeric(r rune) bool {
    return 	(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9')
}

type AVL struct {
	root **AVLNode
}

func (t AVL) Assign(s string, x int) { InsertAVL(t.root, s, x) }

func (t AVL) Lookup(s string) (x int, exists bool) {
	x, exists = Lookup(*(t.root), s)
	return
}

type AVLNode struct {
	k 			string
	v, balance 	int
	parent,
	left, right *AVLNode
}

func InitBinarySearchTree() (t *AVLNode) {
	t = nil
	return
}

func Lookup(t *AVLNode, k string) (x int, exists bool) {
	var n *AVLNode;
	for n = t; n != nil && n.k != k; {
		if k < n.k {
			n = n.left
			exists = true
		} else {
			n = n.right
		}
	}
	if n == nil {
		x = -1
		exists = false
	} else {
		x = n.v
		exists = true
	}
	return
}

func Insert(t **AVLNode, k string, v int) (y *AVLNode) {
	y = &AVLNode{ k: k, v: v }
	if *t == nil {
		*t = y
	} else {
		var x *AVLNode = *t
		for {
			if k == x.k {
				x.v = v
				break
			}
			if k < x.k {
				if x.left == nil {
					x.left = y
					y.parent = x
					break
				}
				x = x.left
			} else {
				if x.right == nil {
					x.right = y
					y.parent = x
					break
				}
				x = x.right
			}
		}
	}
	return
}

func ReplaceAVLNode(t **AVLNode, x, y *AVLNode) {
	if x == *t {
		*t = y
		if y != nil { y.parent = nil }
	} else {
		var p *AVLNode = x.parent
		if y != nil { y.parent = p }
		if p.left == x {
			p.left = y
		} else {
			p.right = y
		}
	}
}

func RotateLeft(t **AVLNode, x *AVLNode) {
	var y *AVLNode = x.right
	ReplaceAVLNode(t, x, y)
	var b *AVLNode = y.left
	if b != nil { b.parent = x }
	x.right = b
	x.parent = y
	y.left = x
	x.balance--
	if y.balance > 0 { x.balance -= y.balance }
	y.balance++
	if x.balance < 0 { y.balance += x.balance }
}

func RotateRight(t **AVLNode, x *AVLNode) {
	var y *AVLNode = x.left
	ReplaceAVLNode(t, x, y)
	var b *AVLNode = y.right
	if b != nil { b.parent = x }
	x.left = b
	x.parent = y
	y.right = x
	x.balance++
	if y.balance < 0 { x.balance -= y.balance }
	y.balance++
	if x.balance > 0 { y.balance += x.balance }
}

func InsertAVL(t **AVLNode, k string, v int) {
	var a *AVLNode = Insert(t, k, v)
	a.balance = 0
	for {
		var x *AVLNode = a.parent
		if x == nil { break }
		if a == x.left {
			x.balance--
			if x.balance == 0 { break }
			if x.balance == -2 {
				if a.balance == 1 { RotateLeft(t, a) }
				RotateRight(t, x)
				break
			}
		} else {
			x.balance++
			if x.balance == 0 { break }
			if x.balance == 2 {
				if a.balance == -1 { RotateRight(t, a) }
				RotateLeft(t, x)
				break
			}
		}
		a = x
	}
}

type SkipListNode struct {
	k		string
	v		int
	m		int
	next 	[]*SkipListNode
}

func InitSkipList(m int) (l *SkipListNode) {
	l = &SkipListNode{ m: m }
	l.next = make([]*SkipListNode, m)
	for i := 0; i < m; i++ { l.next[i] = nil }
	return
}

func Succ(x *SkipListNode) (y *SkipListNode) {
	y = x.next[0]
	return
}

func Skip(l *SkipListNode, k string) (p []*SkipListNode) {
	var x *SkipListNode = l
	p = make([]*SkipListNode, l.m)
	for i := l.m - 1; i >= 0; i-- {
		for ; x.next[i] != nil && x.next[i].k < k; x = x.next[i] {
		}
		p[i] = x
	}
	return
}

func (l *SkipListNode) Lookup(s string) (x int, exists bool) {
	p := Skip(l, s)
	y := Succ(p[0])
	if y != nil && y.k == s {
		x, exists = y.v, true
	} else {
		x, exists = -1, false
	}
	return
}

func (l *SkipListNode) Assign(s string, x int) {
	p := Skip(l, s)
	if p[0].next[0] != nil && p[0].next[0].k == s {
		p[0].next[0].v = x
		return
	}
	y := &SkipListNode{ next: make([]*SkipListNode, l.m), m: l.m, k: s, v: x }
	r := rand.Int() * 2
	i := 0
	for i < l.m && r % 2 == 0 {
		y.next[i] = p[i].next[i]
		p[i].next[i] = y
		i++
		r /= 2
	}
	for i < l.m {
		y.next[i] = nil
		i++
	}
}

type AssocArray interface {
    Assign(s string, x int)
    Lookup(s string) (x int, exists bool)
}

func makeSkipList() AssocArray {
	var p *SkipListNode
	p = InitSkipList(5)
	return AssocArray(p)
}

func makeAVL() AssocArray {
	var p AVL
	n := InitBinarySearchTree()
	p.root = &n
	return AssocArray(p)
}

func main() {
	arrayAVL := makeAVL()
	arraySL := makeSkipList()
	fmt.Println(lex("alpha x1 beta alpha x1 y", arrayAVL))
	fmt.Println(lex("alpha x1 beta alpha x1 y", arraySL))
}
