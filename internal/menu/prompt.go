package menu

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const defaultLabel = "Default (single built-in track)"

// Prompt writes a numbered menu to w — a built-in "default track" choice
// followed by scenarios — reads a selection from r, and returns the
// chosen scenario's Path ("" for the default choice). It returns an
// error if the input isn't a valid choice or the reader is exhausted
// without one.
func Prompt(w io.Writer, r io.Reader, scenarios []Option) (string, error) {
	fmt.Fprintln(w, "Select a configuration to run:")
	fmt.Fprintf(w, "  1) %s\n", defaultLabel)
	for i, opt := range scenarios {
		fmt.Fprintf(w, "  %d) %s\n", i+2, opt.Label)
	}
	fmt.Fprint(w, "Enter choice: ")

	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("menu: reading selection: %w", err)
		}
		return "", fmt.Errorf("menu: no selection given")
	}

	choice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return "", fmt.Errorf("menu: invalid selection %q", scanner.Text())
	}

	if choice == 1 {
		return "", nil
	}
	if index := choice - 2; index >= 0 && index < len(scenarios) {
		return scenarios[index].Path, nil
	}
	return "", fmt.Errorf("menu: selection %d out of range", choice)
}
