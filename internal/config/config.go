// Package config parses the connection settings needed to talk to a TAK
// server: the server address plus the username and password used to
// authenticate and enroll for a client certificate.
package config

import (
	"flag"
	"fmt"
	"strings"
)

// Config holds the connection settings for a TAK server session.
type Config struct {
	ServerAddress string
	Username      string
	Password      string
}

// ParseFlags parses args (excluding the program name) into a Config. It
// returns an error naming every missing required field.
func ParseFlags(args []string) (Config, error) {
	fs := flag.NewFlagSet("gotak", flag.ContinueOnError)
	server := fs.String("server", "", "TAK server IP address or hostname")
	username := fs.String("username", "", "username for certificate enrollment")
	password := fs.String("password", "", "password for certificate enrollment")

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
	}, nil
}
