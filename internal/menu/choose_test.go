package menu

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func bufReader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}

func TestChooseLine_ReturnsSelectedIndex(t *testing.T) {
	var out bytes.Buffer

	got, err := chooseLine(&out, bufReader("1\n"), "Pick one:", []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0 {
		t.Errorf("chooseLine() = %d, want 0", got)
	}

	printed := out.String()
	if !strings.Contains(printed, "a") || !strings.Contains(printed, "b") {
		t.Errorf("printed menu = %q, want it to list both entries", printed)
	}
}

func TestChooseLine_SelectsLaterEntry(t *testing.T) {
	var out bytes.Buffer

	got, err := chooseLine(&out, bufReader("2\n"), "Pick one:", []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 1 {
		t.Errorf("chooseLine() = %d, want 1", got)
	}
}

func TestChooseLine_OutOfRangeSelectionReturnsError(t *testing.T) {
	var out bytes.Buffer

	_, err := chooseLine(&out, bufReader("99\n"), "Pick one:", []string{"a"})
	if err == nil {
		t.Fatal("expected error for out-of-range selection, got nil")
	}
}

func TestChooseLine_NonNumericSelectionReturnsError(t *testing.T) {
	var out bytes.Buffer

	_, err := chooseLine(&out, bufReader("banana\n"), "Pick one:", []string{"a"})
	if err == nil {
		t.Fatal("expected error for non-numeric selection, got nil")
	}
}

func TestChooseLine_EOFReturnsError(t *testing.T) {
	var out bytes.Buffer

	_, err := chooseLine(&out, bufReader(""), "Pick one:", []string{"a"})
	if err == nil {
		t.Fatal("expected error on EOF with no input, got nil")
	}
}

func TestChooseLine_LeavesLaterInputForASubsequentCallOnTheSameReader(t *testing.T) {
	var out bytes.Buffer
	shared := bufReader("1\n2\n")

	first, err := chooseLine(&out, shared, "First:", []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error on first call: %v", err)
	}
	if first != 0 {
		t.Fatalf("first = %d, want 0", first)
	}

	second, err := chooseLine(&out, shared, "Second:", []string{"x", "y"})
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}
	if second != 1 {
		t.Errorf("second = %d, want 1 (the second line, not lost between calls)", second)
	}
}

func TestChooseInteractive_DownDownEnterSelectsThirdEntry(t *testing.T) {
	var out bytes.Buffer

	got, err := chooseInteractive(&out, bufReader("\x1b[B\x1b[B\r"), "Pick one:", []string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 2 {
		t.Errorf("chooseInteractive() = %d, want 2", got)
	}
	if out.Len() == 0 {
		t.Error("expected menu output to be written")
	}
}

func TestChooseInteractive_EnterImmediatelySelectsFirstEntry(t *testing.T) {
	var out bytes.Buffer

	got, err := chooseInteractive(&out, bufReader("\r"), "Pick one:", []string{"a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0 {
		t.Errorf("chooseInteractive() = %d, want 0", got)
	}
}

func TestChooseInteractive_UpFromFirstWrapsToLastEntry(t *testing.T) {
	var out bytes.Buffer

	got, err := chooseInteractive(&out, bufReader("\x1b[A\r"), "Pick one:", []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 1 {
		t.Errorf("chooseInteractive() = %d, want 1 (wrapped to the last entry)", got)
	}
}

func TestChooseInteractive_QReturnsCancelled(t *testing.T) {
	var out bytes.Buffer

	_, err := chooseInteractive(&out, bufReader("q"), "Pick one:", []string{"a"})
	if err != ErrCancelled {
		t.Fatalf("err = %v, want ErrCancelled", err)
	}
}

func TestChooseInteractive_EOFWithoutEnterReturnsError(t *testing.T) {
	var out bytes.Buffer

	_, err := chooseInteractive(&out, bufReader("\x1b[B"), "Pick one:", []string{"a"})
	if err == nil {
		t.Fatal("expected error when input ends before Enter, got nil")
	}
}

func TestChooseInteractive_RedrawsInPlaceOnNavigation(t *testing.T) {
	var out bytes.Buffer

	// entries = ["a"] -> 1 entry -> 7 lines per draw.
	if _, err := chooseInteractive(&out, bufReader("\x1b[B\r"), "Pick one:", []string{"a"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "\x1b[7A") {
		t.Errorf("expected a cursor-up-7-lines escape sequence before the redraw, got:\n%s", out.String())
	}
}

func TestChooseInteractive_LeavesLaterInputForASubsequentCallOnTheSameReader(t *testing.T) {
	var out bytes.Buffer
	// First menu: Enter selects entry 0 immediately. Second menu's bytes
	// ("down, enter") must survive on the same shared reader.
	shared := bufReader("\r\x1b[B\r")

	first, err := chooseInteractive(&out, shared, "First:", []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error on first call: %v", err)
	}
	if first != 0 {
		t.Fatalf("first = %d, want 0", first)
	}

	second, err := chooseInteractive(&out, shared, "Second:", []string{"x", "y"})
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}
	if second != 1 {
		t.Errorf("second = %d, want 1 (the down+enter bytes, not lost between calls)", second)
	}
}
