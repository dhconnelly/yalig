package main

import (
	"fmt"
	"strings"
)

type ValType int

const (
	NumT ValType = iota + 1
	StrT
	FuncT
	ListT
	BoolT
	NullT
)

type Value interface {
	Type() ValType
	fmt.Stringer
}

type NullVal struct{}

func (NullVal) Type() ValType {
	return NullT
}

func (NullVal) String() string {
	return "null"
}

var Null = NullVal{}

type NumVal int

func (NumVal) Type() ValType {
	return NumT
}

func (v NumVal) Value() int {
	return int(v)
}

func (v NumVal) String() string {
	return fmt.Sprintf("%d", v.Value())
}

type StrVal string

func (StrVal) Type() ValType {
	return StrT
}

func (s StrVal) Value() string {
	return string(s)
}

func (s StrVal) String() string {
	return s.Value()
}

type ListVal []Value

func (ListVal) Type() ValType {
	return ListT
}

func (l ListVal) Value() []Value {
	return []Value(l)
}

func (l ListVal) String() string {
	var elems []string
	for _, elem := range l.Value() {
		elems = append(elems, elem.String())
	}
	return "[" + strings.Join(elems, ", ") + "]"
}

type BuiltInFuncVal struct {
	arity int
	f     func(params ...Value) Value
}

func (BuiltInFuncVal) Type() ValType {
	return FuncT
}

func (f BuiltInFuncVal) Value() func(...Value) Value {
	return f.f
}

func (f BuiltInFuncVal) String() string {
	return fmt.Sprintf("func([%d]Value)", f.arity)
}

type BoolVal bool

func (BoolVal) Type() ValType {
	return BoolT
}

func (b BoolVal) Value() bool {
	return bool(b)
}

func (b BoolVal) String() string {
	return fmt.Sprintf("%t", b)
}
