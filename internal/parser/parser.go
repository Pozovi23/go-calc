package parser

import (
	"errors"
	"fmt"

	"github.com/gleb/go-calc/internal/tokenizer"
)

// Node is the base interface for all AST nodes.
type Node interface{}

// NumberNode represents a number literal.
type NumberNode struct {
	Value float64
}

// UnaryNode represents an expression with a unary operator (like -5 or -(2+3)).
type UnaryNode struct {
	Op   tokenizer.TokenType
	Expr Node
}

// BinaryNode represents a binary operation (like 2 + 2 or 2 ^ 3).
type BinaryNode struct {
	Op    tokenizer.TokenType
	Left  Node
	Right Node
}

// Parser builds an AST from a slice of tokens.
type Parser struct {
	tokens []tokenizer.Token
	pos    int
}

// Parse takes tokens and returns the root AST node.
func Parse(tokens []tokenizer.Token) (Node, error) {
	if len(tokens) == 0 || (len(tokens) == 1 && tokens[0].Type == tokenizer.TokenEOF) {
		return nil, errors.New("empty expression")
	}

	p := &Parser{
		tokens: tokens,
		pos:    0,
	}

	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.isAtEnd() {
		tok := p.current()
		return nil, fmt.Errorf("unexpected token %q", tok.Raw)
	}

	return node, nil
}

func (p *Parser) parseExpression() (Node, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}

	for p.current().Type == tokenizer.TokenPlus || p.current().Type == tokenizer.TokenMinus {
		op := p.current().Type
		p.next()

		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}

		left = &BinaryNode{
			Op:    op,
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

func (p *Parser) parseMultiplicative() (Node, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.current().Type == tokenizer.TokenStar ||
		p.current().Type == tokenizer.TokenSlash ||
		p.current().Type == tokenizer.TokenPercent {
		op := p.current().Type
		p.next()

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		left = &BinaryNode{
			Op:    op,
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

func (p *Parser) parseUnary() (Node, error) {
	if p.current().Type == tokenizer.TokenPlus || p.current().Type == tokenizer.TokenMinus {
		op := p.current().Type
		p.next()

		expr, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		return &UnaryNode{
			Op:   op,
			Expr: expr,
		}, nil
	}

	return p.parsePower()
}

// Exponentiation is right-associative: 2 ^ 3 ^ 2 == 2 ^ (3 ^ 2).
func (p *Parser) parsePower() (Node, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if p.current().Type == tokenizer.TokenCaret {
		op := p.current().Type
		p.next()

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		return &BinaryNode{
			Op:    op,
			Left:  left,
			Right: right,
		}, nil
	}

	return left, nil
}

func (p *Parser) parsePrimary() (Node, error) {
	tok := p.current()

	if tok.Type == tokenizer.TokenNumber {
		p.next()
		return &NumberNode{Value: tok.Value}, nil
	}

	if tok.Type == tokenizer.TokenLParen {
		p.next()

		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if p.current().Type != tokenizer.TokenRParen {
			return nil, errors.New("missing closing parenthesis ')'")
		}
		p.next()

		return expr, nil
	}

	if p.isAtEnd() {
		return nil, errors.New("unexpected end of expression")
	}

	return nil, fmt.Errorf("unexpected token %q", tok.Raw)
}

func (p *Parser) current() tokenizer.Token {
	if p.pos >= len(p.tokens) {
		return tokenizer.Token{Type: tokenizer.TokenEOF, Raw: "EOF"}
	}
	return p.tokens[p.pos]
}

func (p *Parser) next() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) isAtEnd() bool {
	return p.current().Type == tokenizer.TokenEOF
}
