package main

import (
	"errors"
	"os/user"

	"github.com/anexia-it/fsquota"
	"github.com/spf13/cobra"
)

func lookupGroupnameByGid(gid string) string {
	if g, err := user.LookupGroupId(gid); err == nil {
		return g.Name
	}
	return gid
}

var cmdGroupReport = &cobra.Command{
	Use:   "report path",
	Short: "Retrieves quota report for a given path",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			err = errors.New("exactly one argument required")
			return
		}

		var report *fsquota.Report
		if report, err = fsquota.GetGroupReport(args[0]); err != nil {
			return
		}

		lookupFn := lookupGroupnameByGid

		if wantNumeric, _ := cmd.Flags().GetBool("numeric"); wantNumeric {
			lookupFn = noopLookup
		}

		printReport(cmd, report, "group", lookupFn)
		return
	},
}

func init() {
	cmdGroupReport.Flags().BoolP("numeric", "n", false, "Print numeric group IDs")
	cmdGroup.AddCommand(cmdGroupReport)
}
