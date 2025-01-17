package main

import (
	"go-shell/internal/shell"
)

// Red if commands not going to work
// Trie for autocomplete
// TODO handle xdg environment variables
// TODO implement globbing

func main() {
	sh := shell.NewShell()
	sh.InitCommandMap()
	sh.Run()
	// for {
	// 	fmt.Print(sh.GetPrompt())
	// 	err := sh.ReadInput()
	// 	if err != nil {
	// 		fmt.Println("Error reading input: %v", err)
	// 		continue
	// 	}
	//
	// 	if sh.Cmd == "" {
	// 		continue
	// 	}
	//
	// 	err = sh.ParseAndExecute()
	// 	if err != nil {
	// 		fmt.Printf("Error: %v", err)
	// 	}
	// }
}
