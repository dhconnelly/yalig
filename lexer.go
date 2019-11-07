package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type TokType int

const (
	LPAREN TokType = iota + 1
	RPAREN
	IDENT
	KEYWORD
	NUM
	STR
	NULL
	QUOTE
	EOF
)

func (typ TokType) String() string {
	switch typ {
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case IDENT:
		return "IDENT"
	case KEYWORD:
		return "KEYWORD"
	case NUM:
		return "NUM"
	case STR:
		return "STR"
	case NULL:
		return "NULL"
	case QUOTE:
		return "QUOTE"
	case EOF:
		return "EOF"
	}
	return ""
}

func isOp(r rune) bool {
	return strings.ContainsRune(`<=+-`, r)
}

var keywords = [...]string{
	"fn",
	"def",
	"defun",
	"if",
	"seq",
}

func isKeyword(s string) bool {
	for _, kw := range keywords {
		if kw == s {
			return true
		}
	}
	return false
}

type Token struct {
	Typ TokType
	Lit string
}

var Empty = Token{EOF, ""}

func (t Token) String() string {
	return fmt.Sprintf("{%s, %s}", t.Typ, t.Lit)
}

type Lexer struct {
	b   *bufio.Reader
	cur Token
	err error
}

func NewLexer(b *bufio.Reader) Lexer {
	l := Lexer{b: b}
	l.advance()
	return l
}

func (l *Lexer) nextChar() (rune, error) {
	inComment := false
	for {
		r, _, err := l.b.ReadRune()
		if err != nil {
			return 0, err
		}
		switch r {
		case ' ':
			continue
		case '\t':
			continue
		case '\n':
			inComment = false
			continue
		case ';':
			inComment = true
			continue
		}
		if inComment {
			continue
		}
		return r, nil
	}
}

func isSep(r rune) bool {
	return strings.ContainsRune(") \n", r)
}

func (l *Lexer) readWhile(first rune, pred func(rune) bool, typ TokType) (string, error) {
	lit := []rune{first}
	for {
		r, _, err := l.b.ReadRune()
		if isSep(r) {
			l.b.UnreadRune()
			return string(lit), nil
		}
		if err != nil && err != io.EOF {
			return "", err
		}
		if !pred(r) {
			return "", fmt.Errorf("expected %s, got %c", typ, r)
		}
		lit = append(lit, r)
	}
}

func (l *Lexer) ident(first rune) (string, error) {
	lit, err := l.readWhile(first, unicode.IsLetter, IDENT)
	if err != nil {
		return "", fmt.Errorf("failed to scan ident: %w", err)
	}
	return lit, nil
}

func (l *Lexer) num(first rune) (string, error) {
	lit, err := l.readWhile(first, unicode.IsDigit, NUM)
	if err != nil {
		return "", fmt.Errorf("failed to scan num: %w", err)
	}
	return lit, nil
}

func (l *Lexer) str(first rune) (string, error) {
	lit := []rune{first}
	for {
		r, _, err := l.b.ReadRune()
		if err != nil {
			return "", fmt.Errorf("failed to scan str: %w", err)
		}
		lit = append(lit, r)
		if r == '"' {
			return string(lit), nil
		}
	}
	return string(lit), nil
}

func (l *Lexer) advance() {
	l.cur, l.err = Empty, nil
	r, err := l.nextChar()
	if err != nil {
		l.err = err
		return
	}

	switch {
	case r == '\'':
		l.cur = Token{QUOTE, `'`}
	case r == '"':
		lit, err := l.str(r)
		if err != nil {
			l.err = err
			return
		}
		l.cur = Token{STR, lit}
	case r == '(':
		l.cur = Token{LPAREN, `(`}
	case r == ')':
		l.cur = Token{RPAREN, `)`}
	case unicode.IsLetter(r):
		lit, err := l.ident(r)
		if err != nil && err != io.EOF {
			l.err = err
			return
		}
		if isKeyword(lit) {
			l.cur = Token{KEYWORD, lit}
		} else {
			l.cur = Token{IDENT, lit}
		}
	case unicode.IsDigit(r):
		lit, err := l.num(r)
		if err != nil && err != io.EOF {
			l.err = err
			return
		}
		l.cur = Token{NUM, lit}
	case isOp(r):
		l.cur = Token{IDENT, string(r)}
	default:
		l.err = fmt.Errorf("failed to scan: unknown token: %c", r)
	}
}

func (l *Lexer) Peek() (Token, error) {
	return l.cur, l.err
}

func (l *Lexer) Next() (Token, error) {
	tok, err := l.cur, l.err
	l.advance()
	return tok, err
}
