// Command gotak simulates air tracks on a TAK server for development use.
package main

import (
	"fmt"
	"os"

	"github.com/jwsmith24/goTak/internal/config"
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "gotak:", err)
		os.Exit(1)
	}

	fmt.Printf("Connecting to TAK server %s as %s...\n", cfg.ServerAddress, cfg.Username)
}
