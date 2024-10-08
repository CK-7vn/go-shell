package parser

import (
	"fmt"
	"go-shell/internal/lexer"
)

type CommandType int

const (
	CommandTypeSimple CommandType = iota
	CommandTypePipe
	CommandTypeBackground
	CommandTypeList
	CommandTypeAnd
)

type Redirection struct {
	Type lexer.TokenType
	File string
}

type Command struct {
	Type         CommandType
	Cmd          string
	Args         []string
	Redirections []Redirection
	Left         *Command
	Right        *Command
}

type Parser struct {
	lexer *lexer.Lexer
	curr  *lexer.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{lexer: l}
}

func (p *Parser) nextToken() error {
	token, err := p.lexer.NextToken()
	if err != nil {
		return err
	}
	p.curr = token
	return nil
}

func (p *Parser) Parse() (*Command, error) {
	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	cmd, err := p.parseCommandList()
	if err != nil {
		return nil, err
	}

	if p.curr.Type != lexer.TokenEOF {
		return nil, fmt.Errorf("unexpected token: %v", p.curr)
	}

	return cmd, nil
}

func (p *Parser) parseCommandList() (*Command, error) {
	cmd, err := p.parseCommand()
	if err != nil {
		return nil, err
	}

	if p.curr.Type == lexer.TokenAnd {
		andCmd := &Command{Type: CommandTypeAnd, Left: cmd}
		err := p.nextToken()
		if err != nil {
			return nil, err
		}
		rightCmd, err := p.parseCommandList()
		if err != nil {
			return nil, err
		}
		andCmd.Right = rightCmd
		return andCmd, nil
	}
	if p.curr.Type == lexer.TokenBackground {
		listCmd := &Command{Type: CommandTypeList, Left: cmd}
		err := p.nextToken()
		if err != nil {
			return nil, err
		}

		rightCmd, err := p.parseCommandList()
		if err != nil {
			return nil, err
		}
		listCmd.Right = rightCmd
		return listCmd, nil
	}

	return cmd, nil
}

func (p *Parser) parseCommand() (*Command, error) {
	cmd := &Command{Type: CommandTypeSimple}

	if p.curr.Type != lexer.TokenWord {
		return nil, fmt.Errorf("Expected command, got %v", p.curr)
	}

	cmd.Cmd = p.curr.Value
	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	for p.curr.Type == lexer.TokenWord {
		cmd.Args = append(cmd.Args, p.curr.Value)
		err := p.nextToken()
		if err != nil {
			return nil, err
		}
	}

	for p.curr.Type == lexer.TokenRedirectIn || p.curr.Type == lexer.TokenRedirectOut {
		redirType := p.curr.Type
		err := p.nextToken()
		if err != nil {
			return nil, err
		}

		if p.curr.Type != lexer.TokenWord {
			return nil, fmt.Errorf("expected filename after redirection, got %v", p.curr)
		}

		cmd.Redirections = append(cmd.Redirections, Redirection{
			Type: redirType,
			File: p.curr.Value,
		})

		err = p.nextToken()
		if err != nil {
			return nil, err
		}
	}

	if p.curr.Type == lexer.TokenPipe {
		pipeCmd := &Command{Type: CommandTypePipe, Left: cmd}
		err := p.nextToken()
		if err != nil {
			return nil, err
		}

		rightCmd, err := p.parseCommand()
		if err != nil {
			return nil, err
		}
		pipeCmd.Right = rightCmd
		cmd = pipeCmd
	}

	return cmd, nil
}
