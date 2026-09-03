package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadEnvFile reads a simple KEY=VALUE env file: blank lines and lines
// starting with '#' are ignored, and a value may be wrapped in matching
// single or double quotes. It returns an *os.PathError satisfying
// os.IsNotExist when path does not exist.
func LoadEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("config: invalid line in env file %s: %q", path, line)
		}

		values[strings.TrimSpace(key)] = unquote(strings.TrimSpace(value))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func unquote(s string) string {
	if len(s) >= 2 {
		quote := s[0]
		if (quote == '"' || quote == '\'') && s[len(s)-1] == quote {
			return s[1 : len(s)-1]
		}
	}
	return s
}
