// Command gotak simulates air tracks on a TAK server for development use.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jwsmith24/goTak/internal/config"
	"github.com/jwsmith24/goTak/internal/enroll"
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "gotak:", err)
		os.Exit(1)
	}

	baseURL := enroll.DefaultBaseURL(cfg.ServerAddress)
	fmt.Printf("Enrolling with %s as %s...\n", baseURL, cfg.Username)

	result, err := enroll.Enroll(context.Background(), enroll.InsecureHTTPClient(), baseURL, cfg.Username, cfg.Password)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gotak: enrollment failed:", err)
		os.Exit(1)
	}

	fmt.Printf("Enrollment succeeded: received client certificate and %d CA certificate(s).\n", len(result.CACertsPEM))
}
