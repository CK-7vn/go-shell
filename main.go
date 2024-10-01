package main

import (
	"fmt"
	"go-shell/internal/lexer"
	"go-shell/internal/parser"
	"strings"
)

func printCommand(cmd *parser.Command, indent string) {
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
		printCommand(cmd.Left, indent+"  ")
	}
	if cmd.Right != nil {
		fmt.Printf("%sRight:\n", indent)
		printCommand(cmd.Right, indent+"  ")
	}
}

// TODO handle xdg environment variables
// TODO immplement CD
// TODO implement globbing
func main() {

	input := "ls -l | grep foo > output.txt < input.txt & echo done"
	lex := lexer.NewLexer(strings.NewReader(input))
	parse := parser.NewParser(lex)

	cmd, err := parse.Parse()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Parsed Command Structure:")
	printCommand(cmd, "")
	// shell := shell.NewShell()
	// shell.InitCommandMap()
	// for {
	// 	fmt.Print(shell.GetPrompt())
	// 	err := shell.ReadInput()
	// 	if err != nil {
	// 		fmt.Println("Error", err)
	// 		continue
	// 	}
	// 	err = shell.Evaluate()
	// 	if err != nil {
	// 		fmt.Printf("Error during evaluate or execution: %v \t\n", err)
	// 	}
	// 	continue
	// }
	//
}

// type Shell struct {
// 	envVars    map[string]string
// 	cmd        string
// 	scanr      *bufio.Scanner
// 	commandMap map[string]func(arg ...string) error
// }
//
// func defaultVar() map[string]string {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		cwd = os.Getenv("PWD")
// 	}
//
// 	return map[string]string{
// 		"PWD":  cwd,
// 		"HOME": os.Getenv("HOME"),
// 		"PATH": os.Getenv("PATH"),
// 		"PS1":  "",
// 	}
//
// }
//
// func NewShell() *Shell {
// 	return &Shell{
// 		envVars:    defaultVar(),
// 		cmd:        "",
// 		scanr:      bufio.NewScanner(os.Stdin),
// 		commandMap: make(map[string]func(arg ...string) error),
// 	}
// }
//
// func (s *Shell) initCommandMap() {
// 	exit := func(args ...string) error { os.Exit(0); return nil }
//
// 	s.commandMap = map[string]func(arg ...string) error{
// 		"help": func(args ...string) error {
// 			fmt.Println("Available commands: ")
// 			fmt.Println("help  - Show help information")
// 			fmt.Println("export - Export local environment variables")
// 			fmt.Println("exit - Exit the shell")
// 			return nil
// 		},
// 		"export": func(args ...string) error {
// 			if len(args) == 0 {
// 				return fmt.Errorf("usage: export VAR=VALUE \n")
// 			}
// 			key := args[0]
// 			val := args[1]
// 			s.handleSetVar(key, val)
// 			return nil
// 		},
// 		"exit": exit,
// 	}
// }
//
// func (s *Shell) GetPrompt() string {
// 	if prompt, exists := s.envVars["PS1"]; exists {
// 		return fmt.Sprintf("%s %s >", prompt, s.envVars["PWD"])
// 	}
// 	//if git repo display git branch etc, BUILD a PS1 Basically
//
// 	return fmt.Sprintf("%s >", s.envVars["PWD"])
// }
//
// func (s *Shell) handleSetVar(key, val string) error {
// 	err := os.Setenv(key, val)
// 	if err != nil {
// 		return err
// 	}
// 	s.envVars[key] = val
// 	return nil
// }
//
// func (s *Shell) readInput() error {
// 	if !s.scanr.Scan() {
// 		return errors.New("Issue in the scanner")
// 	}
// 	//in future cmd will be a slice of runes and will have to handle one at a time, for now
// 	s.cmd = s.scanr.Text()
// 	return nil
// }
//
// func (s *Shell) evaluate() error {
// 	args := strings.Fields(s.cmd)
// 	if len(args) == 0 {
// 		return errors.New("No command")
// 	}
//
// 	cmdName := args[0]
//
// 	var cmdArgs []string
// 	if len(args) > 1 {
// 		cmdArgs = args[1:]
// 	}
// 	if cmdFunc, exists := s.commandMap[cmdName]; exists {
// 		return cmdFunc(cmdArgs...)
// 	}
// 	return s.execute(cmdName, cmdArgs)
//
// }
//
// func (s *Shell) execute(program string, args []string) error {
// 	localPath := s.envVars["PATH"]
// 	splitPath := strings.Split(localPath, ":")
//
// 	var Joined string
// 	for _, dir := range splitPath {
// 		Joined = path.Join(dir, program)
// 		_, err := os.Stat(Joined)
// 		if err != nil {
// 			continue
// 		} else {
// 			break
// 		}
// 	}
// 	cmd := exec.Command(Joined, args...)
// 	cmd.Stdin = os.Stdin
// 	cmd.Stderr = os.Stderr
// 	cmd.Stdout = os.Stdout
// 	return cmd.Run()
// }
