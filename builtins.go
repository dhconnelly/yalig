package main

import (
	"fmt"
)

var builtIns = map[string]BuiltInFuncVal{
	"print": {1, func(args ...Value) (Value, error) {
		fmt.Println(args[0])
		return Null, nil
	}},
	"<": {2, func(args ...Value) (Value, error) {
		left := args[0].(NumVal)
		right := args[1].(NumVal)
		return BoolVal(left.Value() < right.Value()), nil
	}},
	"+": {2, func(args ...Value) (Value, error) {
		left := args[0].(NumVal)
		right := args[1].(NumVal)
		return NumVal(left.Value() + right.Value()), nil
	}},
	"=": {2, func(args ...Value) (Value, error) {
		switch left := args[0].(type) {
		case NumVal:
			right, ok := args[1].(NumVal)
			if !ok {
				return nil, fmt.Errorf("type mismatch: %s and %s", args[0], args[1])
			}
			return BoolVal(left.Value() == right.Value()), nil
		case NullVal:
			if _, ok := args[1].(NullVal); !ok {
				return nil, fmt.Errorf("type mismatch: %s and %s", args[0], args[1])
			}
			return BoolVal(true), nil
		}
		return nil, fmt.Errorf("bad type for '=': %s", args[0])
	}},
	"first": {1, func(args ...Value) (Value, error) {
		list := args[0].(ListVal)
		return list.Value()[0], nil
	}},
	"rest": {1, func(args ...Value) (Value, error) {
		list := args[0].(ListVal)
		return ListVal(list[1:]), nil
	}},
	"empty": {1, func(args ...Value) (Value, error) {
		list, ok := args[0].(ListVal)
		if !ok {
			return nil, fmt.Errorf("'empty' requires a list, got: %s", args[0])
		}
		return BoolVal(len(list.Value()) == 0), nil
	}},
}
