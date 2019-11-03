package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

type tokType int

const (
	LPAREN tokType = iota + 1
	RPAREN
	IDENT
	KEYWORD
	NUM
	STR
	OP
	NULL
	QUOTE
	EOF
)

func (typ tokType) String() string {
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
	case OP:
		return "OP"
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
	"def",
	"fn",
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

type token struct {
	typ tokType
	lit string
}

var empty = token{EOF, ""}

func (t token) String() string {
	return fmt.Sprintf("{%s, %s}", t.typ, t.lit)
}

type lexer struct {
	b   *bufio.Reader
	cur token
	err error
}

func newLexer(b *bufio.Reader) lexer {
	l := lexer{b: b}
	l.advance()
	return l
}

func (l *lexer) nextChar() (rune, error) {
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

func (l *lexer) readWhile(first rune, pred func(rune) bool, typ tokType) (string, error) {
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

func (l *lexer) ident(first rune) (string, error) {
	lit, err := l.readWhile(first, unicode.IsLetter, IDENT)
	if err != nil {
		return "", fmt.Errorf("failed to scan ident: %w", err)
	}
	return lit, nil
}

func (l *lexer) num(first rune) (string, error) {
	lit, err := l.readWhile(first, unicode.IsDigit, NUM)
	if err != nil {
		return "", fmt.Errorf("failed to scan num: %w", err)
	}
	return lit, nil
}

func (l *lexer) str(first rune) (string, error) {
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

func (l *lexer) advance() {
	l.cur, l.err = empty, nil
	r, err := l.nextChar()
	if err != nil {
		l.err = err
		return
	}

	switch {
	case r == '\'':
		l.cur = token{QUOTE, `'`}
	case r == '"':
		lit, err := l.str(r)
		if err != nil {
			l.err = err
			return
		}
		l.cur = token{STR, lit}
	case r == '(':
		l.cur = token{LPAREN, `(`}
	case r == ')':
		l.cur = token{RPAREN, `)`}
	case unicode.IsLetter(r):
		lit, err := l.ident(r)
		if err != nil && err != io.EOF {
			l.err = err
			return
		}
		if isKeyword(lit) {
			l.cur = token{KEYWORD, lit}
		} else {
			l.cur = token{IDENT, lit}
		}
	case unicode.IsDigit(r):
		lit, err := l.num(r)
		if err != nil && err != io.EOF {
			l.err = err
			return
		}
		l.cur = token{NUM, lit}
	case isOp(r):
		l.cur = token{OP, string(r)}
	default:
		l.err = fmt.Errorf("failed to scan: unknown token: %c", r)
	}
}

func (l *lexer) peek() (token, error) {
	return l.cur, l.err
}

func (l *lexer) next() (token, error) {
	tok, err := l.cur, l.err
	l.advance()
	return tok, err
}

func main() {
	// Get the input stream.
	var b *bufio.Reader
	switch len(os.Args) {
	case 1:
		b = bufio.NewReader(os.Stdin)
	case 2:
		in, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		b = bufio.NewReader(in)
	default:
		log.Fatal("Usage: yalig [file]")
	}
	l := newLexer(b)
	for {
		tok, err := l.next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tok)
	}
}
