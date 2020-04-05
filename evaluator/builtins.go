package evaluator

import (
	"dczombera/monkey_language_interpreter/object"
	"fmt"
)

var builtins = map[string]*object.Builtin{
	"len":   &object.Builtin{Fn: builtinLen},
	"first": &object.Builtin{Fn: first},
	"last":  &object.Builtin{Fn: last},
	"rest":  &object.Builtin{Fn: rest},
	"push":  &object.Builtin{Fn: push},
	"puts":  &object.Builtin{Fn: puts},
}

func builtinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments, got=%d, expected=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to 'len' not supported, got=%s, expected=%s",
			args[0].Type(), object.STRING_OBJ)
	}
}

func first(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments, got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to 'first' must be of type %s, got %s",
			object.ARRAY_OBJ, args[0].Type())
	}

	array := args[0].(*object.Array)
	if len(array.Elements) > 0 {
		return array.Elements[0]
	}

	return NULL
}

func last(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments, got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to 'last' must be of type %s, got %s",
			object.ARRAY_OBJ, args[0].Type())
	}

	array := args[0].(*object.Array)
	length := len(array.Elements)
	if length > 0 {
		return array.Elements[length-1]
	}

	return NULL
}

func rest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments, got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to 'rest' must be of type %s, got %s",
			object.ARRAY_OBJ, args[0].Type())
	}

	array := args[0].(*object.Array)
	length := len(array.Elements)
	if length > 0 {
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, array.Elements[1:length])
		return &object.Array{Elements: newElements}
	}

	return NULL
}

func push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments, got=%d, expected=2", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to 'push' must be of type %s, got %s",
			object.ARRAY_OBJ, args[0].Type())
	}

	array := args[0].(*object.Array)
	length := len(array.Elements)

	newElements := make([]object.Object, length+1, length+1)
	copy(newElements, array.Elements[0:length])
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}
}

func puts(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return NULL
}
