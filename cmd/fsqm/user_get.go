package main

import (
	"errors"
	"os/user"

	"github.com/anexia-it/fsquota"
	"github.com/spf13/cobra"
)

var cmdUserGet = &cobra.Command{
	Use:   "get path user",
	Short: "Retrieves quota information for a given user",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 2 {
			err = errors.New("exactly two arguments required")
			return
		}

		var u *user.User
		if u, err = lookupUser(args[1]); err != nil {
			return
		}

		var info *fsquota.Info
		if info, err = fsquota.GetUserInfo(args[0], u); err != nil {
			return
		}

		printQuotaInfo(cmd, info)

		return
	},
}

func init() {
	cmdUser.AddCommand(cmdUserGet)
}
