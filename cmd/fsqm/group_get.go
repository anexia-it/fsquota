package main

import (
	"errors"
	"os/user"

	"github.com/anexia-it/fsquota"
	"github.com/spf13/cobra"
)

var cmdGroupGet = &cobra.Command{
	Use:   "get path group",
	Short: "Retrieves quota information for a given group",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 2 {
			err = errors.New("exactly two arguments required")
			return
		}

		var g *user.Group
		if g, err = lookupGroup(args[1]); err != nil {
			return
		}

		var info *fsquota.Info
		if info, err = fsquota.GetGroupInfo(args[0], g); err != nil {
			return
		}

		printQuotaInfo(cmd, info)

		return
	},
}

func init() {
	cmdGroup.AddCommand(cmdGroupGet)
}
