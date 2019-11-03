package main

import (
	"fmt"
	"io"
	"strconv"
)

type Parser struct {
	l Lexer
}

func NewParser(l Lexer) Parser {
	return Parser{l: l}
}

func (p *Parser) eatLitOrDie(lit string) {
	tok, err := p.l.Next()
	if err != nil {
		panic(err)
	}
	if tok.Lit != lit {
		panic(fmt.Errorf("expected %s, got %s", lit, tok.Lit))
	}
}

func (p *Parser) eat(typ TokType) (Token, error) {
	tok, err := p.l.Next()
	if err != nil {
		return Empty, err
	}
	if tok.Typ != typ {
		return Empty, fmt.Errorf("expected %s, got %s", typ, tok)
	}
	return tok, nil
}

func (p *Parser) callExpr() (*CallExpr, error) {
	fn, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse call: %w", err)
	}
	var args []Expr
	for {
		if tok, _ := p.l.Peek(); tok.Typ == RPAREN {
			break
		}
		arg, err := p.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse call: %w", err)
		}
		args = append(args, arg)
	}
	_, err = p.eat(RPAREN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse call: %w", err)
	}
	return &CallExpr{fn, args}, nil
}

func (p *Parser) funcExpr() (*FuncExpr, error) {
	p.eatLitOrDie("fn")
	_, err := p.eat(LPAREN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fn: %w", err)
	}
	var names []*IdentExpr
	for {
		if tok, _ := p.l.Peek(); tok.Typ == RPAREN {
			break
		}
		name, err := p.identExpr()
		if err != nil {
			return nil, fmt.Errorf("failed to parse fn: %w", err)
		}
		names = append(names, name)
	}
	_, err = p.eat(RPAREN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fn: %w", err)
	}
	body, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse fn: %w", err)
	}
	_, err = p.eat(RPAREN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fn: %w", err)
	}
	return &FuncExpr{names, body}, nil
}

func (p *Parser) defExpr() (*DefExpr, error) {
	p.eatLitOrDie("def")
	name, err := p.identExpr()
	if err != nil {
		return nil, fmt.Errorf("failed to parse def: %w", err)
	}
	binding, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse def: %w", err)
	}
	_, err = p.eat(RPAREN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse def: %w", err)
	}
	return &DefExpr{name, binding}, nil
}

func (p *Parser) ifExpr() (*IfExpr, error) {
	p.eatLitOrDie("if")
	ant, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse if: %w", err)
	}
	con, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse if: %w", err)
	}
	alt, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse if: %w", err)
	}
	_, err = p.eat(RPAREN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse if: %w", err)
	}
	return &IfExpr{ant, con, alt}, nil
}

func (p *Parser) seqExpr() (*SeqExpr, error) {
	p.eatLitOrDie("seq")
	var body []Expr
	for {
		if tok, _ := p.l.Peek(); tok.Typ == RPAREN {
			break
		}
		expr, err := p.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse seq: %w", err)
		}
		body = append(body, expr)
	}
	_, err := p.eat(RPAREN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse seq: %w", err)
	}
	return &SeqExpr{body}, nil
}

func (p *Parser) listExpr() (*ListExpr, error) {
	p.eatLitOrDie("'")
	_, err := p.eat(LPAREN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse list: %w", err)
	}
	var elems []Expr
	for {
		if tok, _ := p.l.Peek(); tok.Typ == RPAREN {
			break
		}
		elem, err := p.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse list: %w", err)
		}
		elems = append(elems, elem)
	}
	_, err = p.eat(RPAREN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse list: %w", err)
	}
	return &ListExpr{elems}, nil
}

func (p *Parser) identExpr() (*IdentExpr, error) {
	tok, err := p.eat(IDENT)
	if err != nil {
		return nil, err
	}
	return &IdentExpr{tok.Lit}, nil
}

func (p *Parser) numExpr() (*NumExpr, error) {
	tok, err := p.eat(NUM)
	if err != nil {
		return nil, err
	}
	num, err := strconv.Atoi(tok.Lit)
	if err != nil {
		return nil, err
	}
	return &NumExpr{num}, nil
}

func (p *Parser) strExpr() (*StrExpr, error) {
	tok, err := p.eat(STR)
	if err != nil {
		return nil, err
	}
	return &StrExpr{tok.Lit}, nil
}

func (p *Parser) sExpr() (Expr, error) {
	_, err := p.eat(LPAREN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expr: %w", err)
	}
	tok, err := p.l.Peek()
	if err != nil {
		return nil, fmt.Errorf("failed to parse expr: %w", err)
	}
	switch tok.Typ {
	case IDENT:
		return p.callExpr()
	case KEYWORD:
		switch tok.Lit {
		case "fn":
			return p.funcExpr()
		case "def":
			return p.defExpr()
		case "if":
			return p.ifExpr()
		case "seq":
			return p.seqExpr()
		}
	}
	return nil, fmt.Errorf("failed to parse expr: bad token %s", tok)
}

func (p *Parser) Parse() (Expr, error) {
	tok, err := p.l.Peek()
	if err != nil {
		return nil, err
	}
	switch tok.Typ {
	case LPAREN:
		return p.sExpr()
	case QUOTE:
		return p.listExpr()
	case IDENT:
		return p.identExpr()
	case NUM:
		return p.numExpr()
	case STR:
		return p.strExpr()
	case EOF:
		return nil, io.EOF
	}
	return nil, fmt.Errorf("failed to parse: unknown token %s", tok)
}
