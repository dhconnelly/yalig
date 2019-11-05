package main

import (
	"fmt"
)

var builtIns = map[string]BuiltInFuncVal{
	"print": {1, func(args ...Value) Value {
		fmt.Println(args[0])
		return Null
	}},
	"<": {2, func(args ...Value) Value {
		left := args[0].(NumVal)
		right := args[1].(NumVal)
		return BoolVal(left.Value() < right.Value())
	}},
	"+": {2, func(args ...Value) Value {
		left := args[0].(NumVal)
		right := args[1].(NumVal)
		return NumVal(left.Value() + right.Value())
	}},
	"=": {2, func(args ...Value) Value {
		left := args[0].(NumVal)
		right := args[1].(NumVal)
		return BoolVal(left.Value() == right.Value())
	}},
}
