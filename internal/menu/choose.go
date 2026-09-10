package menu

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// chooseLine writes a numbered list of entries under title to w, reads a
// selection from r, and returns its index. It returns an error if the
// input isn't a valid choice or r is exhausted without one.
//
// r must be the same *bufio.Reader across a sequence of prompts sharing
// one underlying input stream (e.g. stdin): reading a line here can pull
// in bytes belonging to a later prompt, and only survives to serve that
// prompt if everyone reads through the same buffer instead of each
// wrapping the raw stream fresh.
func chooseLine(w io.Writer, r *bufio.Reader, title string, entries []string) (int, error) {
	fmt.Fprintln(w, title)
	for i, entry := range entries {
		fmt.Fprintf(w, "  %d) %s\n", i+1, entry)
	}
	fmt.Fprint(w, "Enter choice: ")

	line, err := r.ReadString('\n')
	if line == "" && err != nil {
		return 0, fmt.Errorf("menu: reading selection: %w", err)
	}

	choice, convErr := strconv.Atoi(strings.TrimSpace(line))
	if convErr != nil {
		return 0, fmt.Errorf("menu: invalid selection %q", strings.TrimSpace(line))
	}

	if index := choice - 1; index >= 0 && index < len(entries) {
		return index, nil
	}
	return 0, fmt.Errorf("menu: selection %d out of range", choice)
}

// chooseInteractive drives an arrow-key menu against r/w, redrawing the
// framed box after every navigation key, and returns the selected
// entry's index once the user presses Enter. It returns ErrCancelled if
// the user cancels, or an I/O error if r is exhausted before either
// happens.
//
// r must be the same *bufio.Reader across a sequence of prompts sharing
// one underlying input stream, for the same reason as chooseLine.
func chooseInteractive(w io.Writer, r *bufio.Reader, title string, entries []string) (int, error) {
	linesPerDraw := len(entries) + 6 // border, title, blank, entries, blank, footer, border

	selected := 0

	writeFrame(w, Render(title, entries, selected))

	for {
		key, err := readKey(r)
		if err != nil {
			return 0, fmt.Errorf("menu: reading selection: %w", err)
		}

		switch key {
		case KeyQuit:
			return 0, ErrCancelled
		case KeyEnter:
			return selected, nil
		case KeyUp, KeyDown:
			selected = NextIndex(selected, len(entries), key)
			// Move the cursor back to the top of the box and redraw in
			// place, rather than printing a new box below the old one.
			fmt.Fprintf(w, "\x1b[%dA", linesPerDraw)
			writeFrame(w, Render(title, entries, selected))
		}
	}
}
