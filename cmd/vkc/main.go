package main

import (
	"fmt"
	"os"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(cli.ExitCodeFrom(err))
	}
}
