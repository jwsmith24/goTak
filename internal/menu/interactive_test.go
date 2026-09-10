package menu

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func newTestReader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}

func TestNextIndex_DownWrapsAtEnd(t *testing.T) {
	if got := NextIndex(2, 3, KeyDown); got != 0 {
		t.Errorf("NextIndex(2, 3, KeyDown) = %d, want 0", got)
	}
}

func TestNextIndex_DownAdvances(t *testing.T) {
	if got := NextIndex(0, 3, KeyDown); got != 1 {
		t.Errorf("NextIndex(0, 3, KeyDown) = %d, want 1", got)
	}
}

func TestNextIndex_UpWrapsAtStart(t *testing.T) {
	if got := NextIndex(0, 3, KeyUp); got != 2 {
		t.Errorf("NextIndex(0, 3, KeyUp) = %d, want 2", got)
	}
}

func TestNextIndex_UpRetreats(t *testing.T) {
	if got := NextIndex(2, 3, KeyUp); got != 1 {
		t.Errorf("NextIndex(2, 3, KeyUp) = %d, want 1", got)
	}
}

func TestNextIndex_OtherKeysDoNotMove(t *testing.T) {
	if got := NextIndex(1, 3, KeyEnter); got != 1 {
		t.Errorf("NextIndex(1, 3, KeyEnter) = %d, want unchanged 1", got)
	}
}

func TestReadKey_ArrowUp(t *testing.T) {
	key, err := readKey(newTestReader("\x1b[A"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != KeyUp {
		t.Errorf("readKey = %v, want KeyUp", key)
	}
}

func TestReadKey_ArrowDown(t *testing.T) {
	key, err := readKey(newTestReader("\x1b[B"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != KeyDown {
		t.Errorf("readKey = %v, want KeyDown", key)
	}
}

func TestReadKey_Enter(t *testing.T) {
	for _, in := range []string{"\r", "\n"} {
		key, err := readKey(newTestReader(in))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if key != KeyEnter {
			t.Errorf("readKey(%q) = %v, want KeyEnter", in, key)
		}
	}
}

func TestReadKey_QuitKeys(t *testing.T) {
	for _, in := range []string{"q", "Q", "\x03", "\x1b"} {
		key, err := readKey(newTestReader(in))
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", in, err)
		}
		if key != KeyQuit {
			t.Errorf("readKey(%q) = %v, want KeyQuit", in, key)
		}
	}
}

func TestReadKey_UnrecognizedByte(t *testing.T) {
	key, err := readKey(newTestReader("x"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != KeyOther {
		t.Errorf("readKey = %v, want KeyOther", key)
	}
}

func TestReadKey_EOF(t *testing.T) {
	_, err := readKey(newTestReader(""))
	if err == nil {
		t.Fatal("expected error on EOF, got nil")
	}
}

func TestRender_IsRectangularASCIIBox(t *testing.T) {
	out := Render("Select a configuration to run:", []string{
		"Default (single built-in track)",
		"a.json - scenario A",
		"b.json - scenario B",
	}, 1)

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 4 {
		t.Fatalf("got %d lines, want at least 4", len(lines))
	}

	width := len([]rune(stripANSI(lines[0])))
	for i, line := range lines {
		plain := stripANSI(line)
		if len([]rune(plain)) != width {
			t.Errorf("line %d width = %d, want %d (all lines must line up): %q", i, len([]rune(plain)), width, plain)
		}
		for _, r := range plain {
			if r > 127 {
				t.Errorf("line %d contains non-ASCII rune %q: %q", i, r, plain)
			}
		}
	}

	first, last := stripANSI(lines[0]), stripANSI(lines[len(lines)-1])
	if !strings.HasPrefix(first, "+") || !strings.HasSuffix(first, "+") {
		t.Errorf("top border = %q, want it to start and end with '+'", first)
	}
	if !strings.HasPrefix(last, "+") || !strings.HasSuffix(last, "+") {
		t.Errorf("bottom border = %q, want it to start and end with '+'", last)
	}
}

func TestRender_HighlightsOnlySelectedEntry(t *testing.T) {
	out := Render("Title", []string{"one", "two", "three"}, 1)

	lines := strings.Split(out, "\n")
	var highlighted []string
	for _, line := range lines {
		if strings.Contains(line, "\x1b[7m") {
			highlighted = append(highlighted, line)
		}
	}

	if len(highlighted) != 1 {
		t.Fatalf("got %d highlighted lines, want exactly 1: %v", len(highlighted), highlighted)
	}
	if !strings.Contains(highlighted[0], "two") {
		t.Errorf("highlighted line = %q, want it to contain the selected entry %q", highlighted[0], "two")
	}
	for _, line := range lines {
		if strings.Contains(line, "one") || strings.Contains(line, "three") {
			if strings.Contains(line, "\x1b[7m") {
				t.Errorf("unselected entry line unexpectedly highlighted: %q", line)
			}
		}
	}
}

func TestRender_TruncatesLongEntries(t *testing.T) {
	longEntry := strings.Repeat("x", 500)
	out := Render("Title", []string{longEntry}, 0)

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	width := len([]rune(stripANSI(lines[0])))
	for _, line := range lines {
		if got := len([]rune(stripANSI(line))); got != width {
			t.Errorf("line width = %d, want %d for %q", got, width, line)
		}
	}
}

func TestRunInteractive_DownDownEnterSelectsThirdEntry(t *testing.T) {
	scenarios := []Option{
		{Path: "/scenarios/a.json", Label: "a.json - scenario A"},
		{Path: "/scenarios/b.json", Label: "b.json - scenario B"},
	}
	var out bytes.Buffer

	// Entries are [default, a.json, b.json]; down, down lands on b.json.
	got, err := RunInteractive(&out, strings.NewReader("\x1b[B\x1b[B\r"), scenarios)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/scenarios/b.json" {
		t.Errorf("RunInteractive() = %q, want %q", got, "/scenarios/b.json")
	}
	if out.Len() == 0 {
		t.Error("expected menu output to be written")
	}
}

func TestRunInteractive_RedrawsInPlaceOnNavigation(t *testing.T) {
	scenarios := []Option{{Path: "/scenarios/a.json", Label: "a.json"}}
	var out bytes.Buffer

	// entries = [default, a.json] -> 2 entries -> 8 lines per draw.
	if _, err := RunInteractive(&out, strings.NewReader("\x1b[B\r"), scenarios); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "\x1b[8A") {
		t.Errorf("expected a cursor-up-8-lines escape sequence before the redraw, got:\n%s", out.String())
	}
}

func TestRunInteractive_EnterImmediatelySelectsDefault(t *testing.T) {
	scenarios := []Option{{Path: "/scenarios/a.json", Label: "a.json"}}
	var out bytes.Buffer

	got, err := RunInteractive(&out, strings.NewReader("\r"), scenarios)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("RunInteractive() = %q, want empty path for the default entry", got)
	}
}

func TestRunInteractive_UpFromDefaultWrapsToLastEntry(t *testing.T) {
	scenarios := []Option{
		{Path: "/scenarios/a.json", Label: "a.json"},
		{Path: "/scenarios/b.json", Label: "b.json"},
	}
	var out bytes.Buffer

	got, err := RunInteractive(&out, strings.NewReader("\x1b[A\r"), scenarios)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/scenarios/b.json" {
		t.Errorf("RunInteractive() = %q, want %q (wrapped to the last entry)", got, "/scenarios/b.json")
	}
}

func TestRunInteractive_QReturnsError(t *testing.T) {
	scenarios := []Option{{Path: "/scenarios/a.json", Label: "a.json"}}
	var out bytes.Buffer

	_, err := RunInteractive(&out, strings.NewReader("q"), scenarios)
	if err == nil {
		t.Fatal("expected error for cancel key, got nil")
	}
}

func TestRunInteractive_EOFWithoutEnterReturnsError(t *testing.T) {
	scenarios := []Option{{Path: "/scenarios/a.json", Label: "a.json"}}
	var out bytes.Buffer

	_, err := RunInteractive(&out, strings.NewReader("\x1b[B"), scenarios)
	if err == nil {
		t.Fatal("expected error when input ends before Enter, got nil")
	}
}

// stripANSI removes SGR escape sequences (e.g. "\x1b[7m", "\x1b[0m") so
// tests can measure visible width and compare plain text.
func stripANSI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}
			i = j
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
