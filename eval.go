package main

import (
	"fmt"
)

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
		ctx.Set(name, fn)
	}
	return Evaluator{ctx: ctx}
}

func (ev *Evaluator) callLambda(fn LambdaVal, args []Value) (Value, error) {
	// Check arity
	if len(args) != len(fn.fn.Names) {
		return nil, fmt.Errorf("bad arity: got %d, expected %d", len(fn.fn.Names), len(args))
	}

	// Set the context from the captured env
	evalContext := newContext()
	evalContext.up = fn.ctx

	// Bind names to values
	for i := 0; i < len(args); i++ {
		name := fn.fn.Names[i].Ident
		val := args[i]
		evalContext.Set(name, val)
	}

	lambdaEval := &Evaluator{ctx: evalContext}
	return lambdaEval.Eval(fn.fn.Body)
}

func (ev *Evaluator) call(fnVal Value, args []Value) error {
	switch fn := fnVal.(type) {
	case NullVal:
		return fmt.Errorf("can't call null as function")
	case BuiltInFuncVal:
		ev.stack.push(fn.f(args[:fn.arity]...))
		return nil
	case LambdaVal:
		val, err := ev.callLambda(fn, args)
		if err != nil {
			return err
		}
		ev.stack.push(val)
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
	var fn LambdaVal
	fn.ctx = ev.ctx.freeze()
	fn.fn = e
	ev.stack.push(fn)
	return nil
}

func (ev *Evaluator) VisitDef(e *DefExpr) error {
	val, err := ev.Eval(e.Binding)
	if err != nil {
		return err
	}
	ev.ctx.Set(e.Name, val)
	ev.stack.push(Null)
	return nil
}

func (ev *Evaluator) VisitIf(e *IfExpr) error {
	antVal, err := ev.Eval(e.Antecedent)
	if err != nil {
		return err
	}
	if antVal.(BoolVal).Value() {
		conVal, err := ev.Eval(e.Consequent)
		if err != nil {
			return err
		}
		ev.stack.push(conVal)
	} else {
		altVal, err := ev.Eval(e.Alternate)
		if err != nil {
			return err
		}
		ev.stack.push(altVal)
	}
	return nil
}

func (ev *Evaluator) VisitSeq(e *SeqExpr) error {
	var result Value
	for _, expr := range e.Body {
		val, err := ev.Eval(expr)
		if err != nil {
			return err
		}
		result = val
	}
	ev.stack.push(result)
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
	val, err := ev.ctx.Get(e.Ident)
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
