package main

import (
	"fmt"
)

type ValType int

const (
	NumT ValType = iota + 1
	StrT
	FuncT
	ListT
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

type context struct {
	scope map[string]Value
	up    *context
}

func newContext() context {
	return context{scope: make(map[string]Value)}
}

func (ctx *context) set(name string, val Value) {
	ctx.scope[name] = val
}

func (ctx *context) get(name string) Value {
	for cur := ctx; cur != nil; cur = ctx.up {
		if val, ok := cur.scope[name]; ok {
			return val
		}
	}
	return Null
}

type valueStack struct {
	stack []Value
}

func (stack *valueStack) pop() Value {
	if len(stack.stack) == 0 {
		panic(fmt.Errorf("stack empty"))
	}
	end := len(stack.stack) - 1
	val := stack.stack[end]
	stack.stack = stack.stack[:end]
	return val
}

func (stack *valueStack) push(val Value) {
	stack.stack = append(stack.stack, val)
}

type Evaluator struct {
	ctx   context
	stack valueStack
}

func NewEvaluator() Evaluator {
	return Evaluator{ctx: newContext()}
}

func (ev *Evaluator) VisitCall(e *CallExpr) error {
	return nil
}

func (ev *Evaluator) VisitFunc(e *FuncExpr) error {
	return nil
}

func (ev *Evaluator) VisitDef(e *DefExpr) error {
	val, err := ev.Eval(e.Binding)
	if err != nil {
		return err
	}
	ev.ctx.set(e.Name, val)
	ev.stack.push(Null)
	return nil
}

func (ev *Evaluator) VisitIf(e *IfExpr) error {
	return nil
}

func (ev *Evaluator) VisitSeq(e *SeqExpr) error {
	return nil
}

func (ev *Evaluator) VisitList(e *ListExpr) error {
	return nil
}

func (ev *Evaluator) VisitIdent(e *IdentExpr) error {
	return nil
}

func (ev *Evaluator) VisitNum(e *NumExpr) error {
	ev.stack.push(NumVal(e.Num))
	return nil
}

func (ev *Evaluator) VisitStr(e *StrExpr) error {
	return nil
}

func (ev *Evaluator) Eval(e Expr) (Value, error) {
	if err := e.visit(ev); err != nil {
		return nil, err
	}
	return ev.stack.pop(), nil
}
