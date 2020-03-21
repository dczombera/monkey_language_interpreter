# Interpreter for the simplistic monkey language

This interpreter implements the monkey language, which is a very simple language that has C-like syntax and a couple of features such as:

- Variable binding
- Arithmetic expressions
- Integers
- Booleans
- Strings
- Arrays
- Hashes
- Closures
- Built-in functions
- First-class functions
- Higher-order functions

## How does the interpreter work?

The interpreter tokenizes and parses Monkey source code, builds an abstract syntac tree (AST) and evaluates it. The major parts of this interpreter are:

- Lexer
- Parser
- Abstract Syntac Tree (AST)
- Internal Object System
- Evaluate

The interpreter takes source code as input, tokenizes it (i.e. creates tokens) in a process called "lexical analysis" and pipes these tokens to the parser, which then turns them into an AST.

```
---------------      ---------    -------
| Source Code |  ->  | Token | -> | AST |
---------------      ---------    -------
```

## Example

Given the code snippet

```javascript
let square = fn(x) {
    x * x;
};
```

The lexer generates the token sequence

```javascript
{token.LET, "let"},
{token.IDENT, "square"},
{token.ASSIGN, "="},
{token.FUNCTION, "fn"},
{token.LPAREN, "("},
{token.IDENT, "x"},
{token.RPAREN, ")"},
{token.LBRACE, "{"},
{token.IDENT, "x"},
{token.ASTERISK, "*"},
{token.IDENT, "x"},
{token.SEMICOLON, ";"},
{token.RBRACE, "}"},
{token.SEMICOLON, ";"},
```

Which is then used by the parser to produce the AST. When the interpreter evaluates the program, it "walks" the AST and interprets it.

```javascript
#TODO: Add AST output of parser
```
