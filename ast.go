package main

import (
	"fmt"
)

type ExprType int

const (
	Call ExprType = iota + 1
	Func
	Def
	If
	Seq
	List
	Ident
	Num
	Str
)

// Expr := Call | Func | Def | If | Seq | List | IDENT | NUM | STR
type Expr interface {
	fmt.Stringer
}

// Call := "(" Expr Expr* ")"
type CallExpr struct {
	Fn   Expr
	Args []Expr
}

func (e *CallExpr) String() string {
	return fmt.Sprintf("CallExpr(Fn=%s, Args=%s)", e.Fn, e.Args)
}

// Func := "(" "fn" "(" ident* ")" Expr ")"
type FuncExpr struct {
	Names []*IdentExpr
	Body  Expr
}

func (e *FuncExpr) String() string {
	return fmt.Sprintf("FuncExpr(Names=%s, Body=%s", e.Names, e.Body)
}

// Def := "(" "def" ident Expr ")"
type DefExpr struct {
	Name    *IdentExpr
	Binding Expr
}

func (e *DefExpr) String() string {
	return fmt.Sprintf("DefExpr(Name=%s, Binding=%s", e.Name, e.Binding)
}

// If := "(" "if" Expr Expr Expr ")"
type IfExpr struct {
	Antecedent Expr
	Consequent Expr
	Alternate  Expr
}

func (e *IfExpr) String() string {
	return fmt.Sprintf("IfExpr(Antecedent=%s, Consequent=%s, Alternate=%s)", e.Antecedent, e.Consequent, e.Alternate)
}

// Seq := "(" "seq" Expr* ")"
type SeqExpr struct {
	Body []Expr
}

func (e *SeqExpr) String() string {
	return fmt.Sprintf("SeqExpr(Body=%s)", e.Body)
}

// List := QUOTE "(" Expr* ")"
type ListExpr struct {
	Elems []Expr
}

func (e *ListExpr) String() string {
	return fmt.Sprintf("ListExpr(Elems=%s)", e.Elems)
}

type IdentExpr struct {
	Ident string
}

func (e *IdentExpr) String() string {
	return fmt.Sprintf("IdentExpr(\"%s\")", e.Ident)
}

type NumExpr struct {
	Num int
}

func (e *NumExpr) String() string {
	return fmt.Sprintf("NumExpr(%d)", e.Num)
}

type StrExpr struct {
	Str string
}

func (e *StrExpr) String() string {
	return fmt.Sprintf("StrExpr(%s)", e.Str)
}
