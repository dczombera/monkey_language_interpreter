package evaluator

import (
	"dczombera/monkey_language_interpreter/lexer"
	"dczombera/monkey_language_interpreter/object"
	"dczombera/monkey_language_interpreter/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", int64(10)},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", int64(10)},
		{"if (1 < 2) { 10 }", int64(10)},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", int64(20)},
		{"if (1 < 2) { 10 } else { 20 }", int64(10)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int64)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input  string
		output int64
	}{
		{"return 42;", 42},
		{"return 42; 9;", 42},
		{"return 2 * 21; 9;", 42},
		{"9: return 2 * 21; 9;", 42},
		{
			`
			if (42 > 1) {
				if (42 > 2) {
					return 42;
				}

				return 1;
			}
			`,
			42,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.output)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operand: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			` if (10 > 1) {
			  	if (10 > 1) {
					return true + false;
				}
			  return 1; }
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hola" - "mundo"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"Hash key FUNCTION does not implement Hashable interface",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object found. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 42; a;", 42},
		{"let a = 5 * 5; a;", 25},
		{"let a = 42; let b = a; b;", 42},
		{"let a = 42; let b = a; let c = a + b + 42; c;", 126},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 42; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not of type Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("number of parameters is not 1. got=%d", len(fn.Parameters))
	}

	paraName := fn.Parameters[0].String()
	if paraName != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", paraName)
	}

	expectedBody := "(x + 42)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(42);", 42},
		{"let identity = fn(x) { return x; }; identity(42);", 42},
		{"let double = fn(x) { x * 2; }; double(21);", 42},
		{"let add = fn(x, y) { x + y; }; add(21, 21);", 42},
		{"let add = fn(x, y) { x + y; }; add(add(21, 21), 21);", 63},
		{"fn(x) { x; }(42)", 42},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"It's a trap!";`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not of type String. got=%T (%+v)", evaluated, evaluated)
	}

	expected := "It's a trap!"
	if str.Value != expected {
		t.Fatalf("String has wrong value. expected=%q, got=%q", expected, str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hola" + " " + "mundo!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not of type String. got=%T (%+v)", evaluated, evaluated)
	}

	expected := "Hola mundo!"
	if str.Value != expected {
		t.Fatalf("String has wrong value. expected=%q, got=%q", expected, str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hola mundo!")`, 11},
		{`len(1)`, "argument to 'len' not supported, got=INTEGER, expected=STRING"},
		{`len("one", "two")`, "wrong number of arguments, got=2, expected=1"},
		{`len([1, 2, 3, 4])`, 4},
		{`len([1 + 1, 2 * 3])`, 2},
		{`len([])`, 0},
		{`let array = [1, 2, 3 * 3]; len(array)`, 3},
		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`let array = [2 * 2, 3 + 3, 4 / 4]; first(array)`, 4},
		{`last([1, 2, 3])`, 3},
		{`let array = [2 * 2, 3 + 3, 4 / 4]; last(array)`, 1},
		{`last([])`, nil},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`push([1, 2, 3], 4)`, []int{1, 2, 3, 4}},
		{`let array = [1]; push(array, 2)`, []int{1, 2}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case []int:
			testArrayObject(t, evaluated, expected)
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not of type Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not of type Array, got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong number of elements, expected=3, got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{

			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2},
		{
			"[1, 2, 3][2]",
			3},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		}, {
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		}, {
			"[1, 2, 3][3]",
			nil},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
			"one": 10 - 9,
			two: 1 + 1,
			"thr" + "ee": 6 / 2,
			4: 4,
			true: 5,
			false: 6
	}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{(&object.String{Value: "one"}).HashKey(): 1, (&object.String{Value: "two"}).HashKey(): 2, (&object.String{Value: "three"}).HashKey(): 3, (&object.Integer{Value: 4}).HashKey(): 4, TRUE.HashKey(): 5, FALSE.HashKey(): 6}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}
func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{

			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
		{
			`let h = {"respuesta": 42}; h["respuesta"];`,
			42,
		},
		{
			`let h = {"respuesta": 42}; h["hola"];`,
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testArrayObject(t *testing.T, obj object.Object, expected []int) bool {
	array, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not of type Array. got=%T (%+v)", obj, obj)
		return false
	}

	for idx, el := range expected {
		intObj, ok := array.Elements[idx].(*object.Integer)
		if !ok {
			t.Errorf("array element is not of type Integer. got=%T (%+v)", el, el)
			return false
		}
		if intObj.Value != int64(el) {
			t.Errorf("integer has wrong value. expected=%d, got=%d",
				int64(el), intObj.Value)
			return false
		}
	}

	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not of type Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%d, got=%d",
			expected, result.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not of type Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%t, got=%t", expected, result.Value)
		return false
	}

	return true
}
