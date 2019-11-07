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

func (ctx *context) Set(name string, val Value) {
	ctx.scope[name] = val
}

func (ctx *context) Get(name string) (Value, error) {
	for cur := ctx; cur != nil; cur = cur.up {
		if val, ok := cur.scope[name]; ok {
			return val, nil
		}
	}
	return nil, fmt.Errorf("undefined: [%s]", name)
}

func (ctx *context) freeze() *context {
	frozen := newContext()
	for key, val := range ctx.scope {
		frozen.Set(key, val)
	}
	if ctx.up != nil {
		frozen.up = ctx.up.freeze()
	}
	return &frozen
}
