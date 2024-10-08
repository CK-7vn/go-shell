package lexer

import (
	"bufio"
	"io"
	"strings"
)

type TokenType int

const (
	TokenWord TokenType = iota
	TokenPipe
	TokenAnd
	TokenRedirectOut
	TokenRedirectIn
	TokenBackground
	TokenNewLine
	TokenEOF
	TokenDollar
	TokenEquals
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input *bufio.Reader
}

func NewLexer(input io.Reader) *Lexer {
	return &Lexer{input: bufio.NewReader(input)}
}

func (l *Lexer) NextToken() (*Token, error) {
	for {
		r, _, err := l.input.ReadRune()
		if err != nil {
			if err == io.EOF {
				return &Token{Type: TokenEOF}, nil
			}
			return nil, err
		}
		switch r {
		case ' ', '\t':
			continue
		case '\n':
			return &Token{Type: TokenNewLine, Value: "\n"}, nil
		case '|':
			return &Token{Type: TokenPipe, Value: "|"}, nil
		case '>':
			return &Token{Type: TokenRedirectOut, Value: ">"}, nil
		case '<':
			return &Token{Type: TokenRedirectIn, Value: "<"}, nil
		case '&':
			nextRune, _, _ := l.input.ReadRune()
			if nextRune == '&' {
				return &Token{Type: TokenAnd, Value: "&&"}, nil
			}
			l.input.UnreadRune()
			return &Token{Type: TokenBackground, Value: "&"}, nil
		case '$':
			return &Token{Type: TokenDollar, Value: "$"}, nil
		case '=':
			return &Token{Type: TokenEquals, Value: "="}, nil
		default:
			return l.lexWord(r)
		}
	}
}

func (l *Lexer) lexWord(firstChar rune) (*Token, error) {
	var word strings.Builder
	word.WriteRune(firstChar)
	for {
		r, _, err := l.input.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if r == ' ' || r == '\t' || r == '\n' || r == '|' || r == '>' || r == '<' || r == '&' || r == '=' {
			l.input.UnreadRune()
			break
		}
		word.WriteRune(r)
	}
	return &Token{Type: TokenWord, Value: word.String()}, nil
}
