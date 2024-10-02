package main

import (
	"fmt"
	"go-shell/internal/shell"
)

// TODO handle xdg environment variables
// TODO immplement CD
// TODO implement globbing
func main() {
	sh := shell.NewShell()
	sh.InitCommandMap()
	for {
		fmt.Print(sh.GetPrompt())
		err := sh.ReadInput()
		if err != nil {
			fmt.Println("Error reading input: %v", err)
			continue
		}

		if sh.Cmd == "" {
			continue
		}

		err = sh.ParseAndExecute()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}
}
