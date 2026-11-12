package parser_test

import (
	"testing"

	"github.com/gleb/go-calc/internal/parser"
	"github.com/gleb/go-calc/internal/tokenizer"
)

func TestParse_Success(t *testing.T) {
	validExpressions := []string{
		"42",
		"2 + 2",
		"2 + 3 * 4",
		"(2 + 3) * 4",
		"-(2 + 3) * 4",
		"2--2",
		"2 ^ 3 ^ 2",
		"2 ^ -1",
		"10 % 3",
		"+5 + +2",
	}

	for _, expr := range validExpressions {
		tokens, err := tokenizer.Tokenize(expr)
		if err != nil {
			t.Fatalf("tokenize error for %q: %v", expr, err)
		}

		node, err := parser.Parse(tokens)
		if err != nil {
			t.Fatalf("parse error for %q: %v", expr, err)
		}
		if node == nil {
			t.Fatalf("expected node, got nil for %q", expr)
		}
	}
}

func TestParse_Errors(t *testing.T) {
	invalidExpressions := []string{
		"",
		"(2 + 3",
		"2 + 3)",
		"()",
		"2 +",
		"2 + * 3",
		"2 3",
	}

	for _, expr := range invalidExpressions {
		tokens, err := tokenizer.Tokenize(expr)
		if err != nil {
			continue // tokenizing failed, which is also an error
		}

		_, err = parser.Parse(tokens)
		if err == nil {
			t.Errorf("expected parse error for %q, got nil", expr)
		}
	}
}
