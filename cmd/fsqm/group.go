package main

import (
	"os/user"

	"github.com/spf13/cobra"
)

var cmdGroup = &cobra.Command{
	Use:   "group",
	Short: "Group quota management",
}

func init() {
	cmdRoot.AddCommand(cmdGroup)
}

func lookupGroup(groupIdOrGroupName string) (grp *user.Group, err error) {
	if isNumeric(groupIdOrGroupName) {
		grp = &user.Group{
			Gid: groupIdOrGroupName,
		}
		return
	}
	return user.LookupGroup(groupIdOrGroupName)
}
