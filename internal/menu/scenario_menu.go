package menu

import (
	"bufio"
	"io"
	"os"
)

// RunScenarioMenu shows a menu of scenarios (plus a built-in "default
// track" choice) and returns the selected scenario's Path (empty for
// the default). See choose for how the menu is presented and why input
// must be shared across every menu shown during a run.
func RunScenarioMenu(stdin *os.File, input *bufio.Reader, stdout io.Writer, scenarios []Option) (string, error) {
	entries := make([]string, len(scenarios)+1)
	entries[0] = defaultLabel
	for i, opt := range scenarios {
		entries[i+1] = opt.Label
	}

	idx, err := choose(stdin, input, stdout, "Select a configuration to run:", entries)
	if err != nil {
		return "", err
	}
	if idx == 0 {
		return "", nil
	}
	return scenarios[idx-1].Path, nil
}
