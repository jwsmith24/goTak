// Package config parses the connection settings needed to talk to a TAK
// server: the server address plus the username and password used to
// authenticate and enroll for a client certificate.
package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config holds the connection settings for a TAK server session.
type Config struct {
	ServerAddress string
	Username      string
	Password      string
	ScenarioPath  string // optional; empty means use the built-in default track
}

// dotEnvPath is the env file ParseFlags looks for in the working
// directory. It is optional: flags always take precedence, and a missing
// file is not an error.
const dotEnvPath = ".env"

// Env file keys read as fallback values for the matching flag.
const (
	envServerKey   = "GOTAK_SERVER"
	envUsernameKey = "GOTAK_USERNAME"
	envPasswordKey = "GOTAK_PASSWORD"
	envScenarioKey = "GOTAK_SCENARIO"
)

// ParseFlags parses args (excluding the program name) into a Config,
// falling back to a .env file (GOTAK_SERVER, GOTAK_USERNAME,
// GOTAK_PASSWORD, GOTAK_SCENARIO) in the working directory for any flag
// not given on the command line. It returns an error naming every
// required field still missing once both sources are applied;
// ScenarioPath is optional.
func ParseFlags(args []string) (Config, error) {
	return parseFlags(args, dotEnvPath)
}

func parseFlags(args []string, envFilePath string) (Config, error) {
	envValues, err := LoadEnvFile(envFilePath)
	if err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config: reading env file %s: %w", envFilePath, err)
	}

	fs := flag.NewFlagSet("gotak", flag.ContinueOnError)
	server := fs.String("server", envValues[envServerKey], "TAK server IP address or hostname")
	username := fs.String("username", envValues[envUsernameKey], "username for certificate enrollment")
	password := fs.String("password", envValues[envPasswordKey], "password for certificate enrollment")
	scenarioPath := fs.String("scenario", envValues[envScenarioKey], "path to a JSON scenario file (optional; defaults to a single built-in track)")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	var missing []string
	if *server == "" {
		missing = append(missing, "server")
	}
	if *username == "" {
		missing = append(missing, "username")
	}
	if *password == "" {
		missing = append(missing, "password")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required field(s): %s", strings.Join(missing, ", "))
	}

	return Config{
		ServerAddress: *server,
		Username:      *username,
		Password:      *password,
		ScenarioPath:  *scenarioPath,
	}, nil
}
