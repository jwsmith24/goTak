package menu

import (
	"bytes"
	"strings"
	"testing"
)

func TestPromptCustomLocation_BuildsLocationFromInput(t *testing.T) {
	var out bytes.Buffer

	loc, err := promptCustomLocation(bufReader("Home Base\n30.0\n-97.0\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.Name != "Home Base" || loc.Lat != 30.0 || loc.Lon != -97.0 {
		t.Errorf("promptCustomLocation() = %+v, unexpected", loc)
	}
}

func TestPromptCustomLocation_DefaultsNameWhenBlank(t *testing.T) {
	var out bytes.Buffer

	loc, err := promptCustomLocation(bufReader("\n30.0\n-97.0\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.Name != defaultCustomLocationName {
		t.Errorf("Name = %q, want default %q", loc.Name, defaultCustomLocationName)
	}
}

func TestPromptCustomLocation_RejectsNonNumericLatitude(t *testing.T) {
	var out bytes.Buffer

	_, err := promptCustomLocation(bufReader("Home\nnotanumber\n-97.0\n"), &out)
	if err == nil {
		t.Fatal("expected error for non-numeric latitude, got nil")
	}
}

func TestPromptCustomLocation_RejectsOutOfRangeLatitude(t *testing.T) {
	var out bytes.Buffer

	_, err := promptCustomLocation(bufReader("Home\n120\n-97.0\n"), &out)
	if err == nil {
		t.Fatal("expected error for out-of-range latitude, got nil")
	}
	if !strings.Contains(err.Error(), "latitude") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "latitude")
	}
}

func TestPromptCustomLocation_PrintsPrompts(t *testing.T) {
	var out bytes.Buffer

	if _, err := promptCustomLocation(bufReader("Home\n30.0\n-97.0\n"), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	printed := out.String()
	for _, want := range []string{"Name", "Latitude", "Longitude"} {
		if !strings.Contains(printed, want) {
			t.Errorf("printed prompts = %q, want it to mention %q", printed, want)
		}
	}
}
