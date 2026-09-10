package menu

import (
	"bufio"
	"io"
	"os"

	"golang.org/x/term"
)

// choose shows a menu of entries on stdin/stdout: an arrow-key
// navigable, highlighted menu when stdin is an interactive terminal, or
// a numbered list read one line at a time otherwise — e.g. when input
// is piped, such as in scripts or CI. It returns the selected entry's
// index.
//
// input must be the single *bufio.Reader shared across every menu shown
// during a run (wrapping stdin exactly once, at the top of main), so a
// sequence of prompts doesn't drop bytes buffered-but-unread by an
// earlier one.
func choose(stdin *os.File, input *bufio.Reader, stdout io.Writer, title string, entries []string) (int, error) {
	fd := int(stdin.Fd())
	if !term.IsTerminal(fd) {
		return chooseLine(stdout, input, title, entries)
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return chooseLine(stdout, input, title, entries)
	}
	defer term.Restore(fd, oldState)

	return chooseInteractive(stdout, input, title, entries)
}
