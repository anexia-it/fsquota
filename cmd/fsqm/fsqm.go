package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := cmdRoot.Execute(); err != nil {
		// Return exit code 1 on error
		os.Exit(1)
	}
}

var cmdRoot = &cobra.Command{
	Use:   "fsqm",
	Short: "filesystem quota manager",
}
