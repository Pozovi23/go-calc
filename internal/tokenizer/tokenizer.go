package tokenizer

import (
	"fmt"
	"strconv"
)

// TokenType represents the type of a token.
type TokenType int

// Token types supported by the calculator.
const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenPercent
	TokenCaret
	TokenLParen
	TokenRParen
)

// Token holds a token type, string representation, and numeric value.
type Token struct {
	Type  TokenType
	Raw   string
	Value float64
}

// Tokenize splits an expression string into tokens.
func Tokenize(input string) ([]Token, error) {
	var tokens []Token
	n := len(input)
	i := 0

	for i < n {
		ch := input[i]

		if isSpace(ch) {
			i++
			continue
		}

		switch ch {
		case '+':
			tokens = append(tokens, Token{Type: TokenPlus, Raw: "+"})
			i++
			continue
		case '-':
			tokens = append(tokens, Token{Type: TokenMinus, Raw: "-"})
			i++
			continue
		case '*':
			tokens = append(tokens, Token{Type: TokenStar, Raw: "*"})
			i++
			continue
		case '/':
			tokens = append(tokens, Token{Type: TokenSlash, Raw: "/"})
			i++
			continue
		case '%':
			tokens = append(tokens, Token{Type: TokenPercent, Raw: "%"})
			i++
			continue
		case '^':
			tokens = append(tokens, Token{Type: TokenCaret, Raw: "^"})
			i++
			continue
		case '(':
			tokens = append(tokens, Token{Type: TokenLParen, Raw: "("})
			i++
			continue
		case ')':
			tokens = append(tokens, Token{Type: TokenRParen, Raw: ")"})
			i++
			continue
		}

		if isDigit(ch) || (ch == '.' && i+1 < n && isDigit(input[i+1])) {
			start := i
			hasDot := false
			hasExp := false

			for i < n {
				c := input[i]
				if isDigit(c) {
					i++
				} else if c == '.' && !hasDot && !hasExp {
					hasDot = true
					i++
				} else if (c == 'e' || c == 'E') && !hasExp {
					hasExp = true
					i++
					if i < n && (input[i] == '+' || input[i] == '-') {
						i++
					}
					if i >= n || !isDigit(input[i]) {
						return nil, fmt.Errorf("invalid scientific notation at index %d", start)
					}
				} else {
					break
				}
			}

			raw := input[start:i]
			val, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid number %q: %w", raw, err)
			}

			tokens = append(tokens, Token{
				Type:  TokenNumber,
				Raw:   raw,
				Value: val,
			})
			continue
		}

		return nil, fmt.Errorf("unknown character '%c'", ch)
	}

	tokens = append(tokens, Token{Type: TokenEOF, Raw: "EOF"})
	return tokens, nil
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
