package input

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

// Read from os.Stdin file descriptor, append rune to os.Stdout, when enter is hit, flush buffer, parse. If rune/input is ASCII character move appropriately
//for each rune/character read, append it to the slice of bytes, and increase pointer(char location), if arrow keys are used ++ or --
//location of pointer, and location?

type InputReader struct {
	buffer []rune
	pos    int
}

func NewInputReader() *InputReader {
	return &InputReader{
		buffer: make([]rune, 0),
		pos:    0,
	}
}
func (ir *InputReader) Read() (string, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	reader := bufio.NewReader(os.Stdin)

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == bufio.ErrBufferFull {
				return "", nil
			}
			return "", err
		}
		switch r {
		case '\n':
			input := string(ir.buffer)
			fmt.Print("\n")
			ir.buffer = make([]rune, 0)
			ir.pos = 0
			return input, nil
		case 127:
			if ir.pos > 0 {
				ir.buffer = append(ir.buffer[:ir.pos-1], ir.buffer[ir.pos:]...)
				ir.pos--
				fmt.Print("\b")
				ir.redrawFromCursor()
				fmt.Print(" \b")
			}
		case '\x1b':
			r2, _, _ := reader.ReadRune()
			r3, _, _ := reader.ReadRune()
			if r2 == '[' {
				switch r3 {
				case 'C':
					if ir.pos < len(ir.buffer) {
						fmt.Print(string(ir.buffer[ir.pos]))
						ir.pos++
					}
				case 'D':
					if ir.pos > 0 {
						fmt.Print("\b")
						ir.pos--
					}
				}
			}
		default:
			ir.buffer = append(ir.buffer[:ir.pos], append([]rune{r}, ir.buffer[ir.pos:]...)...)
			ir.pos++
			fmt.Print(string(r))
			ir.redrawFromCursor()
		}
	}
}
func (ir *InputReader) redrawFromCursor() {
	fmt.Print("\033[K")
	for _, r := range ir.buffer[ir.pos:] {
		fmt.Print(string(r))
	}
	for range ir.buffer[ir.pos:] {
		fmt.Print("\b")
	}
}

func (ir *InputReader) InsertRune(r rune) {
	ir.buffer = append(ir.buffer[:ir.pos], append([]rune{r}, ir.buffer[ir.pos:]...)...)
	ir.pos++
}

func (ir *InputReader) MoveCursorLeft() {
	if ir.pos > 0 {
		ir.pos--
	}
}
func (ir *InputReader) MoveCursorRight() {
	if ir.pos < len(ir.buffer) {
		ir.pos++
	}
}
func (ir *InputReader) DeleteRune() {
	if ir.pos > 0 {
		ir.buffer = append(ir.buffer[:ir.pos-1], ir.buffer[ir.pos:]...)
		ir.pos--
	}
}
