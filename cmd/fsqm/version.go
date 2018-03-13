package main

import (
	"github.com/anexia-it/fsquota"
	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Show the fsqm version information",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("fsqm v%s\n", fsquota.VersionString())
		cmd.Println("Copyright (C) 2018 Anexia Internetdienstleistungs GmbH")
		cmd.Println("License: MIT")
	},
}

func init() {
	cmdRoot.AddCommand(cmdVersion)
}
