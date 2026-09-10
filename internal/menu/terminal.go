package menu

import (
	"os"

	"golang.org/x/term"
)

// RunMenu shows the scenario picker on stdin/stdout: an arrow-key
// navigable, highlighted menu when stdin is an interactive terminal, or
// the numbered-list Prompt (reads a line, no raw mode) otherwise — e.g.
// when input is piped, such as in scripts or CI. It returns the chosen
// scenario's path (empty for the default built-in track).
func RunMenu(stdin, stdout *os.File, scenarios []Option) (string, error) {
	fd := int(stdin.Fd())
	if !term.IsTerminal(fd) {
		return Prompt(stdout, stdin, scenarios)
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return Prompt(stdout, stdin, scenarios)
	}
	defer term.Restore(fd, oldState)

	return RunInteractive(stdout, stdin, scenarios)
}
