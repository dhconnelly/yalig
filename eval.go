package main

import (
	"fmt"
)

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

func (ctx *context) get(name string) (Value, error) {
	for cur := ctx; cur != nil; cur = ctx.up {
		if val, ok := cur.scope[name]; ok {
			return val, nil
		}
	}
	return nil, fmt.Errorf("undefined: [%s]", name)
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
	ctx := newContext()
	for name, fn := range builtIns {
		ctx.set(name, fn)
	}
	return Evaluator{ctx: ctx}
}

func (ev *Evaluator) call(fnVal Value, args []Value) error {
	switch fn := fnVal.(type) {
	case NullVal:
		return fmt.Errorf("can't call null as function")
	case BuiltInFuncVal:
		ev.stack.push(fn.f(args[:fn.arity]...))
		return nil
	default:
		panic(fmt.Errorf("unknown function type: %s", fnVal))
	}
}

func (ev *Evaluator) VisitCall(e *CallExpr) error {
	val, err := ev.Eval(e.Fn)
	if err != nil {
		return err
	}
	var args []Value
	for _, arg := range e.Args {
		val, err := ev.Eval(arg)
		if err != nil {
			return err
		}
		args = append(args, val)
	}
	return ev.call(val, args)
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
	var elems []Value
	for _, elem := range e.Elems {
		val, err := ev.Eval(elem)
		if err != nil {
			return err
		}
		elems = append(elems, val)
	}
	ev.stack.push(ListVal(elems))
	return nil
}

func (ev *Evaluator) VisitIdent(e *IdentExpr) error {
	val, err := ev.ctx.get(e.Ident)
	if err != nil {
		return err
	}
	ev.stack.push(val)
	return nil
}

func (ev *Evaluator) VisitNum(e *NumExpr) error {
	ev.stack.push(NumVal(e.Num))
	return nil
}

func (ev *Evaluator) VisitStr(e *StrExpr) error {
	ev.stack.push(StrVal(e.Str))
	return nil
}

func (ev *Evaluator) Eval(e Expr) (Value, error) {
	if err := e.visit(ev); err != nil {
		return nil, err
	}
	return ev.stack.pop(), nil
}
