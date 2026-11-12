package eval

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gleb/go-calc/internal/parser"
	"github.com/gleb/go-calc/internal/tokenizer"
)

// Eval takes an expression string and returns the computed value.
func Eval(expr string) (float64, error) {
	tokens, err := tokenizer.Tokenize(expr)
	if err != nil {
		return 0, fmt.Errorf("token error: %w", err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return 0, fmt.Errorf("syntax error: %w", err)
	}

	return evalNode(ast)
}

func evalNode(node parser.Node) (float64, error) {
	switch n := node.(type) {
	case *parser.NumberNode:
		return n.Value, nil

	case *parser.UnaryNode:
		val, err := evalNode(n.Expr)
		if err != nil {
			return 0, err
		}

		if n.Op == tokenizer.TokenMinus {
			return checkResult(-val)
		}
		if n.Op == tokenizer.TokenPlus {
			return checkResult(val)
		}
		return 0, fmt.Errorf("unsupported unary operator: %v", n.Op)

	case *parser.BinaryNode:
		left, err := evalNode(n.Left)
		if err != nil {
			return 0, err
		}

		right, err := evalNode(n.Right)
		if err != nil {
			return 0, err
		}

		switch n.Op {
		case tokenizer.TokenPlus:
			return checkResult(left + right)
		case tokenizer.TokenMinus:
			return checkResult(left - right)
		case tokenizer.TokenStar:
			return checkResult(left * right)
		case tokenizer.TokenSlash:
			if right == 0 {
				return 0, errors.New("division by zero")
			}
			return checkResult(left / right)
		case tokenizer.TokenPercent:
			if right == 0 {
				return 0, errors.New("modulo by zero")
			}
			return checkResult(math.Mod(left, right))
		case tokenizer.TokenCaret:
			if left == 0 && right < 0 {
				return 0, errors.New("division by zero")
			}
			return checkResult(math.Pow(left, right))
		default:
			return 0, fmt.Errorf("unsupported binary operator: %v", n.Op)
		}

	default:
		return 0, fmt.Errorf("unknown node type: %T", node)
	}
}

func checkResult(val float64) (float64, error) {
	if math.IsInf(val, 0) {
		return 0, errors.New("result overflow: infinity")
	}
	if math.IsNaN(val) {
		return 0, errors.New("invalid operation: result is NaN")
	}
	return val, nil
}

// Format formats the number with the given precision and trims trailing zeros.
func Format(val float64, precision int) string {
	if precision < 0 {
		precision = 0
	}
	if val == 0 {
		return "0"
	}

	if math.Abs(val) >= 1e16 || math.Abs(val) < 1e-6 {
		return strconv.FormatFloat(val, 'g', precision, 64)
	}

	res := strconv.FormatFloat(val, 'f', precision, 64)
	if strings.Contains(res, ".") {
		res = strings.TrimRight(res, "0")
		res = strings.TrimRight(res, ".")
	}
	if res == "-0" || res == "" {
		return "0"
	}
	return res
}
