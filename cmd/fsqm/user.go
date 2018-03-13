package main

import (
	"os/user"

	"github.com/spf13/cobra"
)

var cmdUser = &cobra.Command{
	Use:   "user",
	Short: "User quota management",
}

func init() {
	cmdRoot.AddCommand(cmdUser)
}

func lookupUser(userIdOrUsername string) (usr *user.User, err error) {
	if isNumeric(userIdOrUsername) {
		usr = &user.User{
			Uid: userIdOrUsername,
		}
		return
	}
	return user.Lookup(userIdOrUsername)
}
