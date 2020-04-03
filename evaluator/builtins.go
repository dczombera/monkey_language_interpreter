package evaluator

import "dczombera/monkey_language_interpreter/object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{Fn: builtinLen},
}

func builtinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments, got=%d, expected=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError("argument to 'len' not supported, got=%s, expected=%s",
			args[0].Type(), object.STRING_OBJ)
	}
}
