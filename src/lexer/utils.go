package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func (l *Lexer) readChar() {
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
	l.column++
	l.position = l.readPosition
	if l.readPosition >= len(l.input) {
		l.ch = 0
		return
	}
	ch, width := utf8.DecodeRuneInString(l.input[l.readPosition:])
	l.ch = ch
	l.readPosition += width
}

func (l *Lexer) readNumber() string {
	position := l.position
	usedDot := false
	for isDigit(l.ch) || (l.ch == '.' && !usedDot) {
		if l.ch == '.' {
			usedDot = true
		}
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() (string, error) {
	l.readChar()
	start := l.position
	for l.position < len(l.input) {
		if l.ch == '"' {
			value := l.input[start:l.position]
			l.readChar()
			return value, nil
		}
		l.readChar()
	}
	return "", fmt.Errorf("unterminated string")
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readComment() string {
	position := l.position
	for l.ch != '\n' && l.position < len(l.input) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readMultiLineComment() (string, error) {
	position := l.position
	for l.position < len(l.input) {
		if l.ch == '*' && l.peekChar() == '/' {
			end := l.position
			l.readChar()
			l.readChar()
			return l.input[position:end], nil
		}
		l.readChar()
	}
	return "", fmt.Errorf("unterminated comment")
}
