package shell

import (
	"bufio"
	"errors"
	"fmt"
	"go-shell/internal/lexer"
	"go-shell/internal/parser"
	"go-shell/internal/trie"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Shell struct {
	envVars      map[string]string
	Cmd          string
	scanr        *bufio.Scanner
	commandMap   map[string]func(arg ...string) error
	commandTrie  *trie.Trie
	currentInput string
}

func defaultVar() map[string]string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = os.Getenv("PWD")
	}

	return map[string]string{
		"PWD":  cwd,
		"HOME": os.Getenv("HOME"),
		"PATH": os.Getenv("PATH"),
		"PS1":  "",
	}

}
func NewShell() *Shell {
	return &Shell{
		envVars:     defaultVar(),
		Cmd:         "",
		scanr:       bufio.NewScanner(os.Stdin),
		commandMap:  make(map[string]func(arg ...string) error),
		commandTrie: trie.NewTrie(),
	}
}

func (s *Shell) GetPrompt() string {
	if prompt, exists := s.envVars["PS1"]; exists {
		return fmt.Sprintf("%s %s >\n >", prompt, s.envVars["PWD"])
	}
	//if git repo display git branch etc, BUILD a PS1 Basically
	return fmt.Sprintf("%s >", s.envVars["PWD"])

}

func (s *Shell) handleSetVar(key, val string) error {
	err := os.Setenv(key, val)
	if err != nil {
		return err
	}
	s.envVars[key] = val
	return nil
}

func (s *Shell) InitCommandMap() {
	exit := func(args ...string) error { os.Exit(0); return nil }

	s.commandMap = map[string]func(arg ...string) error{
		"help": func(args ...string) error {
			fmt.Println("Available commands: ")
			fmt.Println("help  - Show help information")
			fmt.Println("export - Export local environment variables")
			fmt.Println("exit - Exit the shell")
			fmt.Println("printCmd - format: printCmd <command> prints the AST representation of the command")
			return nil
		},
		"export": func(args ...string) error {
			if len(args) == 0 {
				return fmt.Errorf("usage: export VAR=VALUE \n")
			}
			key := args[0]
			val := args[1]
			s.handleSetVar(key, val)
			return nil
		},
		"cd": func(args ...string) error {
			if len(args) == 0 {
				err := os.Chdir(s.envVars["HOME"])
				if err != nil {
					fmt.Printf("Failed to cd to home: %v", err)
				}
				cwd, err := os.Getwd()
				if err != nil {
					cwd = os.Getenv("PWD")
				}
				s.envVars["PWD"] = cwd
			} else {
				err := os.Chdir(args[0])
				if err != nil {
					fmt.Printf("Failed to cd to %s, err: %v", args[0], err)
				}
				cwd, err := os.Getwd()
				if err != nil {
					cwd = os.Getenv("PWD")
				}
				s.envVars["PWD"] = cwd

			}

			return nil
		},
		"pwd": func(args ...string) error {
			fmt.Printf("%s \n", s.envVars["PWD"])
			return nil
		},
		"exit": exit,
		"printCmd": func(args ...string) error {
			if len(args) == 0 {
				return fmt.Errorf("usage: printCmd <command>\n")
			}

			cmdString := strings.Join(args, " ")
			fmt.Printf("cmdString is: %s\n", cmdString)

			lex := lexer.NewLexer(strings.NewReader(cmdString))
			parse := parser.NewParser(lex)
			cmd, err := parse.Parse()
			if err != nil {
				fmt.Printf("error parsing command")
			}
			fmt.Printf("AST for command: %s\n", cmdString)
			s.printCommand(cmd, " ")
			return nil
		},
	}
}

func (s *Shell) executeExternal(program string, args []string) error {
	localPath := s.envVars["PATH"]
	splitPath := strings.Split(localPath, ":")

	var Joined string
	for _, dir := range splitPath {
		Joined = path.Join(dir, program)
		_, err := os.Stat(Joined)
		if err != nil {
			continue
		} else {
			break
		}
	}
	cmd := exec.Command(Joined, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (s *Shell) ReadInput() error {
	if !s.scanr.Scan() {
		return errors.New("Issue in the scanner")
	}
	s.Cmd = s.scanr.Text()
	return nil
}

func (s *Shell) ParseAndExecute() error {
	fields := strings.Fields(s.Cmd)
	if len(fields) > 0 && fields[0] == "printCmd" {
		return s.commandMap["printCmd"](fields[1:]...)
	}
	lex := lexer.NewLexer(strings.NewReader(s.Cmd))
	parse := parser.NewParser(lex)
	cmd, err := parse.Parse()
	if err != nil {
		return fmt.Errorf("error parsing cmd: %v", err)
	}
	return s.executeCommand(cmd)

}

func (s *Shell) executeCommand(cmd *parser.Command) error {
	switch cmd.Type {
	case parser.CommandTypeSimple:
		if builtin, ok := s.commandMap[cmd.Cmd]; ok {
			return builtin(cmd.Args...)
		}
		return s.executeExternal(cmd.Cmd, cmd.Args)
	case parser.CommandTypePipe:
		return fmt.Errorf("Pipe execution not implemented yet")
	case parser.CommandTypeBackground:
		return fmt.Errorf("Background execution not implemented yet")
	case parser.CommandTypeList:
		err := s.executeCommand(cmd.Left)
		if err != nil {
			return err
		}
		return s.executeCommand(cmd.Right)
	case parser.CommandTypeAnd:
		err := s.executeCommand(cmd.Left)
		if err != nil {
			return err
		}
		return s.executeCommand(cmd.Right)
	default:
		return fmt.Errorf("Unknown command type")
	}
}

func (s *Shell) Execute(program string, args []string) error {
	localPath := s.envVars["PATH"]
	splitPath := strings.Split(localPath, ":")

	var Joined string
	for _, dir := range splitPath {
		Joined = path.Join(dir, program)
		_, err := os.Stat(Joined)
		if err != nil {
			continue
		} else {
			break
		}
	}
	cmd := exec.Command(Joined, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (s *Shell) ParseCommand() (*parser.Command, error) {
	lex := lexer.NewLexer(strings.NewReader(s.Cmd))
	parse := parser.NewParser(lex)
	return parse.Parse()
}

func (s *Shell) printCommand(cmd *parser.Command, indent string) {
	fmt.Printf("%sCommand Type: %v\n", indent, cmd.Type)
	if cmd.Cmd != "" {
		fmt.Printf("%s Cmd: %v\n", indent, cmd.Cmd)
	}
	if len(cmd.Args) > 0 {
		fmt.Printf("%sArgs: %v\n", indent, cmd.Args)
	}
	for _, redir := range cmd.Redirections {
		fmt.Printf("%sRedirection: %v -> %s\n", indent, redir.Type, redir.File)
	}
	if cmd.Left != nil {
		fmt.Printf("%sLeft:\n", indent)
		s.printCommand(cmd.Left, indent+"  ")
	}
	if cmd.Right != nil {
		fmt.Printf("%sRight:\n", indent)
		s.printCommand(cmd.Right, indent+"  ")
	}
}
