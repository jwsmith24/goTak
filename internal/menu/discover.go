// Package menu presents an interactive terminal menu for choosing which
// scenario file to simulate.
package menu

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jwsmith24/goTak/internal/scenario"
)

// Option is one scenario a user can choose from the menu.
type Option struct {
	Path  string // scenario file path
	Label string // display text, shown next to the option's number
}

// DiscoverScenarios scans dir for *.json scenario files that parse
// successfully, returning one Option per file sorted by filename.
// Invalid or non-scenario JSON files are skipped. A missing directory is
// not an error: it simply yields no options.
func DiscoverScenarios(dir string) []Option {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)

	options := make([]Option, 0, len(names))
	for _, name := range names {
		path := filepath.Join(dir, name)
		sc, err := scenario.Load(path)
		if err != nil {
			continue
		}
		options = append(options, Option{Path: path, Label: label(name, sc)})
	}
	return options
}

func label(filename string, sc scenario.Scenario) string {
	description := sc.Description
	if description == "" {
		description = fmt.Sprintf("%d track(s)", len(sc.Tracks))
	}
	return fmt.Sprintf("%s - %s", filename, description)
}
