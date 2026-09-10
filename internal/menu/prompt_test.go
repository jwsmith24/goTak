package menu

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrompt_SelectingDefaultReturnsEmptyPath(t *testing.T) {
	scenarios := []Option{
		{Path: "/scenarios/a.json", Label: "a.json - scenario A"},
		{Path: "/scenarios/b.json", Label: "b.json - scenario B"},
	}
	var out bytes.Buffer

	got, err := Prompt(&out, strings.NewReader("1\n"), scenarios)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("Prompt() = %q, want empty path for the default option", got)
	}

	printed := out.String()
	if !strings.Contains(printed, "a.json - scenario A") || !strings.Contains(printed, "b.json - scenario B") {
		t.Errorf("printed menu = %q, want it to list both scenarios", printed)
	}
}

func TestPrompt_SelectingScenarioReturnsItsPath(t *testing.T) {
	scenarios := []Option{
		{Path: "/scenarios/a.json", Label: "a.json - scenario A"},
		{Path: "/scenarios/b.json", Label: "b.json - scenario B"},
	}
	var out bytes.Buffer

	// Option 1 is the built-in default, so 3 selects the second scenario.
	got, err := Prompt(&out, strings.NewReader("3\n"), scenarios)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/scenarios/b.json" {
		t.Errorf("Prompt() = %q, want %q", got, "/scenarios/b.json")
	}
}

func TestPrompt_NoScenarios_StillOffersDefault(t *testing.T) {
	var out bytes.Buffer

	got, err := Prompt(&out, strings.NewReader("1\n"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("Prompt() = %q, want empty path", got)
	}
}

func TestPrompt_OutOfRangeSelectionReturnsError(t *testing.T) {
	scenarios := []Option{{Path: "/scenarios/a.json", Label: "a.json"}}
	var out bytes.Buffer

	_, err := Prompt(&out, strings.NewReader("99\n"), scenarios)
	if err == nil {
		t.Fatal("expected error for out-of-range selection, got nil")
	}
}

func TestPrompt_NonNumericSelectionReturnsError(t *testing.T) {
	scenarios := []Option{{Path: "/scenarios/a.json", Label: "a.json"}}
	var out bytes.Buffer

	_, err := Prompt(&out, strings.NewReader("banana\n"), scenarios)
	if err == nil {
		t.Fatal("expected error for non-numeric selection, got nil")
	}
}

func TestPrompt_EOFReturnsError(t *testing.T) {
	scenarios := []Option{{Path: "/scenarios/a.json", Label: "a.json"}}
	var out bytes.Buffer

	_, err := Prompt(&out, strings.NewReader(""), scenarios)
	if err == nil {
		t.Fatal("expected error on EOF with no input, got nil")
	}
}
