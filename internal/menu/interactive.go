package menu

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Key is a decoded terminal input event relevant to menu navigation.
type Key int

const (
	KeyOther Key = iota
	KeyUp
	KeyDown
	KeyEnter
	KeyQuit // Ctrl+C, Esc, or 'q': cancel selection
)

// ErrCancelled is returned by RunInteractive when the user cancels the
// menu (Ctrl+C, Esc, or 'q') instead of selecting an entry.
var ErrCancelled = errors.New("menu: selection cancelled")

const menuContentWidth = 76 // characters between the "| " and " |" borders

// NextIndex returns the selection index after applying key to a menu
// with count entries, wrapping around at either end. Keys other than
// KeyUp/KeyDown leave the index unchanged.
func NextIndex(current, count int, key Key) int {
	switch key {
	case KeyDown:
		if current >= count-1 {
			return 0
		}
		return current + 1
	case KeyUp:
		if current <= 0 {
			return count - 1
		}
		return current - 1
	default:
		return current
	}
}

// readKey decodes one key event from r: printable/control bytes directly,
// and "ESC [ A"/"ESC [ B" arrow-key escape sequences. A lone ESC (not
// followed by '[') is treated as a cancel, matching typical terminal
// behavior for pressing Escape.
func readKey(r *bufio.Reader) (Key, error) {
	b, err := r.ReadByte()
	if err != nil {
		return KeyOther, err
	}

	switch b {
	case '\r', '\n':
		return KeyEnter, nil
	case 3, 'q', 'Q': // Ctrl+C
		return KeyQuit, nil
	case 0x1b: // ESC
		b2, err := r.ReadByte()
		if err != nil || b2 != '[' {
			return KeyQuit, nil
		}
		b3, err := r.ReadByte()
		if err != nil {
			return KeyOther, nil
		}
		switch b3 {
		case 'A':
			return KeyUp, nil
		case 'B':
			return KeyDown, nil
		default:
			return KeyOther, nil
		}
	default:
		return KeyOther, nil
	}
}

// Render draws entries as an ASCII-framed box, with the entry at
// selected highlighted in reverse video and a marker ("> ") prefix.
// Every line is padded/truncated to the same fixed width so the frame
// stays rectangular regardless of terminal width or entry length.
func Render(title string, entries []string, selected int) string {
	var b strings.Builder

	border := "+" + strings.Repeat("-", menuContentWidth+2) + "+"
	blank := "|" + strings.Repeat(" ", menuContentWidth+2) + "|"

	b.WriteString(border + "\n")
	b.WriteString(frameLine(title) + "\n")
	b.WriteString(blank + "\n")
	for i, entry := range entries {
		marker := "  "
		if i == selected {
			marker = "> "
		}
		line := frameLine(marker + entry)
		if i == selected {
			line = highlight(line)
		}
		b.WriteString(line + "\n")
	}
	b.WriteString(blank + "\n")
	b.WriteString(frameLine("(Use up/down arrows and Enter; q to cancel)") + "\n")
	b.WriteString(border + "\n")

	return b.String()
}

// frameLine pads or truncates text to menuContentWidth and wraps it in
// the "| ... |" border.
func frameLine(text string) string {
	runes := []rune(text)
	if len(runes) > menuContentWidth {
		if menuContentWidth > 3 {
			text = string(runes[:menuContentWidth-3]) + "..."
		} else {
			text = string(runes[:menuContentWidth])
		}
		runes = []rune(text)
	}
	text = text + strings.Repeat(" ", menuContentWidth-len(runes))
	return "| " + text + " |"
}

// highlight wraps a rendered frame line's inner content in reverse-video
// so only the text between the border pipes is highlighted, leaving the
// pipes themselves unstyled.
func highlight(line string) string {
	inner := strings.TrimSuffix(strings.TrimPrefix(line, "| "), " |")
	return "| \x1b[7m" + inner + "\x1b[0m |"
}

// RunInteractive drives an arrow-key menu against r/w: a built-in
// "default track" entry followed by scenarios. It redraws the framed
// menu after every navigation key and returns the selected scenario's
// Path (empty for the default entry) once the user presses Enter. It
// returns ErrCancelled if the user cancels, or an I/O error if r is
// exhausted before either happens.
func RunInteractive(w io.Writer, r io.Reader, scenarios []Option) (string, error) {
	entries := make([]string, 0, len(scenarios)+1)
	entries = append(entries, defaultLabel)
	for _, opt := range scenarios {
		entries = append(entries, opt.Label)
	}

	const title = "Select a configuration to run:"
	linesPerDraw := len(entries) + 6 // border, title, blank, entries, blank, footer, border

	br := bufio.NewReader(r)
	selected := 0

	writeFrame(w, Render(title, entries, selected))

	for {
		key, err := readKey(br)
		if err != nil {
			return "", fmt.Errorf("menu: reading selection: %w", err)
		}

		switch key {
		case KeyQuit:
			return "", ErrCancelled
		case KeyEnter:
			if selected == 0 {
				return "", nil
			}
			return scenarios[selected-1].Path, nil
		case KeyUp, KeyDown:
			selected = NextIndex(selected, len(entries), key)
			// Move the cursor back to the top of the box and redraw in
			// place, rather than printing a new box below the old one.
			fmt.Fprintf(w, "\x1b[%dA", linesPerDraw)
			writeFrame(w, Render(title, entries, selected))
		}
	}
}

// writeFrame writes s with "\r\n" line endings. A raw terminal has
// output post-processing (OPOST) disabled, so a bare "\n" doesn't
// return the cursor to column 0 the way it does in normal cooked mode -
// without the "\r" every line after the first would drift further
// right. Render itself stays newline-only so its output is simple to
// test.
func writeFrame(w io.Writer, s string) {
	fmt.Fprint(w, strings.ReplaceAll(s, "\n", "\r\n"))
}
