package main

import (
	"errors"
	"os/user"

	"github.com/anexia-it/fsquota"
	"github.com/spf13/cobra"
)

func lookupUsernameByUid(uid string) string {
	if u, err := user.LookupId(uid); err == nil {
		return u.Username
	}
	return uid
}

var cmdUserReport = &cobra.Command{
	Use:   "report path",
	Short: "Retrieves quota report for a given path",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			err = errors.New("exactly one argument required")
			return
		}

		var report *fsquota.Report
		if report, err = fsquota.GetUserReport(args[0]); err != nil {
			return
		}

		lookupFn := lookupUsernameByUid

		if wantNumeric, _ := cmd.Flags().GetBool("numeric"); wantNumeric {
			lookupFn = noopLookup
		}

		printReport(cmd, report, "user", lookupFn)
		return
	},
}

func init() {
	cmdUserReport.Flags().BoolP("numeric", "n", false, "Print numeric user IDs")
	cmdUser.AddCommand(cmdUserReport)
}
