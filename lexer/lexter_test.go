package lexer

import (
	"dczombera/monkey_language_interpreter/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := "=+(){},;"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	lexer := New(input)

	for index, test := range tests {
		token := lexer.NextToken()

		if token.Type != test.expectedType {
			t.Fatalf("tests[%d] - type of token is wrong. Expected=%q, got=%q", index, test.expectedType, token.Type)
		}

		if token.Literal != test.expectedLiteral {
			t.Fatalf("tests[%d] - literal of token is wrong. Exepcted=%q, got=%q", index, test.expectedLiteral, token.Literal)
		}
	}
}
