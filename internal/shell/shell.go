package shell

import (
	"bufio"
	"errors"
	"fmt"
	"go-shell/internal/parser"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Shell struct {
	envVars    map[string]string
	cmd        string
	scanr      *bufio.Scanner
	commandMap map[string]func(arg ...string) error
	parser     *parser.Parser
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
		envVars:    defaultVar(),
		cmd:        "",
		scanr:      bufio.NewScanner(os.Stdin),
		commandMap: make(map[string]func(arg ...string) error),
	}
}
func (s *Shell) GetPrompt() string {
	if prompt, exists := s.envVars["PS1"]; exists {
		return fmt.Sprintf("%s %s >", prompt, s.envVars["PWD"])
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
		"exit": exit,
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

func (s *Shell) executeNode(node *parser.Node) error {
	if node.Type != parser.NTExec {
		return errors.New("unsupported node type")
	}

	if len(node.Argv) == 0 {
		return errors.New("No command provided")
	}

	cmdName := node.Argv[0]
	cmdArgs := node.Argv[1:]
	if cmdFunc, exists := s.commandMap[cmdName]; exists {
		return cmdFunc(cmdArgs...)
	}
	return s.executeExternal(cmdName, cmdArgs)
}

func (s *Shell) ReadInput() error {
	if !s.scanr.Scan() {
		return errors.New("Issue in the scanner")
	}
	//in future cmd will be a slice of runes and will have to handle one at a time, for now
	s.cmd = s.scanr.Text()
	return nil
}

func (s *Shell) Evaluate() error {
	args := strings.Fields(s.cmd)
	if len(args) == 0 {
		return errors.New("No command")
	}

	cmdName := args[0]

	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	}
	if cmdFunc, exists := s.commandMap[cmdName]; exists {
		return cmdFunc(cmdArgs...)
	}
	return s.Execute(cmdName, cmdArgs)

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
