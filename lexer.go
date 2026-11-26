package main

type TokenType string

const (
	IDENT    TokenType = "IDENT"
	NUMBER   TokenType = "NUMBER"
	STRING   TokenType = "STRING"
	LPAREN   TokenType = "LPAREN"
	RPAREN   TokenType = "RPAREN"
	EQUALS   TokenType = "EQUALS"
	PLUS     TokenType = "PLUS"
	MINUS    TokenType = "MINUS"
	STAR     TokenType = "STAR"
	SLASH    TokenType = "SLASH"
	LBRACE   TokenType = "LBRACE"
	RBRACE   TokenType = "RBRACE"
	COMMA    TokenType = "COMMA"
	COLON    TokenType = "COLON"
	DOT      TokenType = "DOT"
	PERCENT  TokenType = "PERCENT"
	LBRACKET TokenType = "LBRACKET"
	RBRACKET TokenType = "RBRACKET"
	EQEQ     TokenType = "EQEQ"  // ==
	NOTEQ    TokenType = "NOTEQ" // !=
	LT       TokenType = "LT"    // <
	GT       TokenType = "GT"    // >
	LTE      TokenType = "LTE"   // <=
	GTE      TokenType = "GTE"   // >=
	ILLEGAL  TokenType = "ILLEGAL"
)

// A Token is one small piece of the program. Line says which line of the file it came from.
type Token struct {
	Type  TokenType
	Value string
	Line  int
}

// unterminatedString is the Value of the ILLEGAL token made for a string with no closing quote.
const unterminatedString = "unterminated string"

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func Lex(src string) []Token {
	var tokens []Token
	line := 1 // the line we are on; goes up by one at every newline
	add := func(t TokenType, value string) {
		tokens = append(tokens, Token{t, value, line})
	}

	i := 0
	for i < len(src) {
		c := src[i]
		switch {
		case c == '\n':
			line++
			i++
		case c == ' ' || c == '\t' || c == '\r':
			i++
		case c == '(':
			add(LPAREN, "(")
			i++
		case c == ')':
			add(RPAREN, ")")
			i++
		case c == '{':
			add(LBRACE, "{")
			i++
		case c == '}':
			add(RBRACE, "}")
			i++
		case c == ',':
			add(COMMA, ",")
			i++
		case c == ':':
			add(COLON, ":")
			i++
		case c == '.':
			add(DOT, ".")
			i++
		case c == '%':
			add(PERCENT, "%")
			i++
		case c == '[':
			add(LBRACKET, "[")
			i++
		case c == ']':
			add(RBRACKET, "]")
			i++
		// Operators made of two characters: look one character ahead.
		// "=" alone is assignment, "==" is a comparison. Same idea for the others.
		case c == '=' && i+1 < len(src) && src[i+1] == '=':
			add(EQEQ, "==")
			i += 2
		case c == '!' && i+1 < len(src) && src[i+1] == '=':
			add(NOTEQ, "!=")
			i += 2
		case c == '<' && i+1 < len(src) && src[i+1] == '=':
			add(LTE, "<=")
			i += 2
		case c == '>' && i+1 < len(src) && src[i+1] == '=':
			add(GTE, ">=")
			i += 2
		case c == '<':
			add(LT, "<")
			i++
		case c == '>':
			add(GT, ">")
			i++
		case c == '=':
			add(EQUALS, "=")
			i++
		case c == '+':
			add(PLUS, "+")
			i++
		case c == '-':
			add(MINUS, "-")
			i++
		case c == '*':
			add(STAR, "*")
			i++
		case c == '/' && i+1 < len(src) && src[i+1] == '/':
			// a comment: ignore everything up to the end of the line
			// (the newline itself is left alone so line counting stays right)
			for i < len(src) && src[i] != '\n' {
				i++
			}
		case c == '/':
			add(SLASH, "/")
			i++
		case c == '"':
			startLine := line // a string is reported on the line where it starts
			start := i + 1
			i++
			for i < len(src) && src[i] != '"' {
				if src[i] == '\n' {
					line++
				}
				i++
			}
			if i >= len(src) {
				tokens = append(tokens, Token{ILLEGAL, unterminatedString, startLine})
			} else {
				tokens = append(tokens, Token{STRING, src[start:i], startLine})
				i++
			}
		case isLetter(c):
			start := i
			for i < len(src) && (isLetter(src[i]) || isDigit(src[i])) {
				i++
			}
			add(IDENT, src[start:i])
		case isDigit(c):
			start := i
			for i < len(src) && isDigit(src[i]) {
				i++
			}
			// a decimal point followed by digits is part of the number: 3.14
			if i+1 < len(src) && src[i] == '.' && isDigit(src[i+1]) {
				i++
				for i < len(src) && isDigit(src[i]) {
					i++
				}
			}
			add(NUMBER, src[start:i])
		default:
			add(ILLEGAL, string(c))
			i++
		}
	}
	return tokens
}
