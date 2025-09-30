package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
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
		return 0, fmt.Errorf("на позиции %d ожидалось число", p.pos)
	}
	num, _ := strconv.Atoi(p.input[start:p.pos])
	return num, nil
}

// parseFactor — числа или скобки
func (p *Parser) parseFactor() (int, error) {
	p.skipWhitespace()
	if p.peek() == '(' {
		p.next()
		val, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		p.skipWhitespace()
		if p.peek() != ')' {
			return 0, fmt.Errorf("на позиции %d нет закрывающей скобки", p.pos)
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
				return 0, fmt.Errorf("деление на ноль")
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
