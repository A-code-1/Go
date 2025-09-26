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
func (p *Parser) parseNumber() int {
	p.skipWhitespace()
	start := p.pos
	for unicode.IsDigit(p.peek()) {
		p.next()
	}
	num, _ := strconv.Atoi(p.input[start:p.pos])
	return num
}

// parseFactor — числа или скобки
func (p *Parser) parseFactor() int {
	p.skipWhitespace()
	if p.peek() == '(' {
		p.next()
		val := p.parseExpression()
		p.skipWhitespace()
		if p.peek() == ')' {
			p.next()
		}
		return val
	}
	return p.parseNumber()
}

// parseTerm — умножение и деление
func (p *Parser) parseTerm() int {
	val := p.parseFactor()
	for {
		p.skipWhitespace()
		switch p.peek() {
		case '*':
			p.next()
			val *= p.parseFactor()
		case '/':
			p.next()
			val /= p.parseFactor()
		default:
			return val
		}
	}
}

// parseExpression — сложение и вычитание
func (p *Parser) parseExpression() int {
	val := p.parseTerm()
	for {
		p.skipWhitespace()
		switch p.peek() {
		case '+':
			p.next()
			val += p.parseTerm()
		case '-':
			p.next()
			val -= p.parseTerm()
		default:
			return val
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
		result := parser.parseExpression()
		fmt.Println(result)
	}
}
