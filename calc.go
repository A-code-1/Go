package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrDivideByZero   = errors.New("деление на ноль")
	ErrExpectedNumber = errors.New("ожидалось число")
	ErrUnclosedParen  = errors.New("нет закрывающей скобки")
)

type Parser struct {
	input string
	pos   int
}

func (p *Parser) peek() rune {
	if p.pos >= len(p.input) {
		return 0
	}
	return rune(p.input[p.pos])
}

func (p *Parser) next() rune {
	if p.pos >= len(p.input) {
		return 0
	}
	ch := rune(p.input[p.pos])
	p.pos++
	return ch
}

func (p *Parser) skipWhitespace() {
	for unicode.IsSpace(p.peek()) {
		p.next()
	}
}

// parseNumber — парсит число
func (p *Parser) parseNumber() (int, error) {
	p.skipWhitespace()
	start := p.pos
	for unicode.IsDigit(p.peek()) {
		p.next()
	}
	if start == p.pos {
		return 0, fmt.Errorf("%w на позиции %d", ErrExpectedNumber, p.pos)
	}
	num, _ := strconv.Atoi(p.input[start:p.pos])
	return num, nil
}

// parseFactor — числа или скобки
func (p *Parser) parseFactor() (int, error) {
	p.skipWhitespace()
	if p.peek() == '+' {
		p.next()
		return p.parseFactor()
	}
	if p.peek() == '-' {
		p.next()
		val, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		return -val, nil
	}
	if p.peek() == '(' {
		p.next()
		val, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		p.skipWhitespace()
		if p.peek() != ')' {
			return 0, fmt.Errorf("%w на позиции %d", ErrUnclosedParen, p.pos)
		}
		p.next()
		return val, nil
	}
	return p.parseNumber()
}

// parseTerm — умножение и деление
func (p *Parser) parseTerm() (int, error) {
	val, err := p.parseFactor()
	if err != nil {
		return 0, err
	}
	for {
		p.skipWhitespace()
		switch p.peek() {
		case '*':
			p.next()
			rhs, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			val *= rhs
		case '/':
			p.next()
			rhs, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			if rhs == 0 {
				return 0, ErrDivideByZero
			}
			val /= rhs
		default:
			return val, nil
		}
	}
}

// parseExpression — сложение и вычитание
func (p *Parser) parseExpression() (int, error) {
	val, err := p.parseTerm()
	if err != nil {
		return 0, err
	}
	for {
		p.skipWhitespace()
		switch p.peek() {
		case '+':
			p.next()
			rhs, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			val += rhs
		case '-':
			p.next()
			rhs, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			val -= rhs
		default:
			return val, nil
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parser := &Parser{input: line}
		result, err := parser.parseExpression()
		if err != nil {
			fmt.Println("Ошибка:", err)
		}
		fmt.Println(result)
	}
}
