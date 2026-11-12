package tokenizer_test

import (
	"testing"

	"github.com/gleb/go-calc/internal/tokenizer"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTypes []tokenizer.TokenType
	}{
		{
			name:      "simple addition",
			input:     "2 + 2",
			wantTypes: []tokenizer.TokenType{tokenizer.TokenNumber, tokenizer.TokenPlus, tokenizer.TokenNumber, tokenizer.TokenEOF},
		},
		{
			name:  "all operators",
			input: "(1 + 2 - 3 * 4 / 5 % 6 ^ 7)",
			wantTypes: []tokenizer.TokenType{
				tokenizer.TokenLParen,
				tokenizer.TokenNumber,
				tokenizer.TokenPlus,
				tokenizer.TokenNumber,
				tokenizer.TokenMinus,
				tokenizer.TokenNumber,
				tokenizer.TokenStar,
				tokenizer.TokenNumber,
				tokenizer.TokenSlash,
				tokenizer.TokenNumber,
				tokenizer.TokenPercent,
				tokenizer.TokenNumber,
				tokenizer.TokenCaret,
				tokenizer.TokenNumber,
				tokenizer.TokenRParen,
				tokenizer.TokenEOF,
			},
		},
		{
			name:  "floats and scientific notation",
			input: ".5 + 5. * 1e-3",
			wantTypes: []tokenizer.TokenType{
				tokenizer.TokenNumber,
				tokenizer.TokenPlus,
				tokenizer.TokenNumber,
				tokenizer.TokenStar,
				tokenizer.TokenNumber,
				tokenizer.TokenEOF,
			},
		},
		{
			name:  "consecutive minuses",
			input: "2--2",
			wantTypes: []tokenizer.TokenType{
				tokenizer.TokenNumber,
				tokenizer.TokenMinus,
				tokenizer.TokenMinus,
				tokenizer.TokenNumber,
				tokenizer.TokenEOF,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenizer.Tokenize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tokens) != len(tt.wantTypes) {
				t.Fatalf("got %d tokens, want %d", len(tokens), len(tt.wantTypes))
			}
			for i, exp := range tt.wantTypes {
				if tokens[i].Type != exp {
					t.Errorf("token %d: got type %v, want %v", i, tokens[i].Type, exp)
				}
			}
		})
	}
}

func TestTokenize_Errors(t *testing.T) {
	badInputs := []string{
		"2 + @",
		"x + 5",
		".",
		"1e",
		"1e+",
		"1e-",
	}

	for _, input := range badInputs {
		_, err := tokenizer.Tokenize(input)
		if err == nil {
			t.Errorf("expected error for input %q, got nil", input)
		}
	}
}
