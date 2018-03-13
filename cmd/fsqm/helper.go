package main

import (
	"errors"
	"strings"
	"unicode"

	"github.com/anexia-it/fsquota"
	"github.com/dustin/go-humanize"
	"github.com/speijnik/go-errortree"
	"github.com/spf13/cobra"
)

func humanizeInodes(inodes uint64) string {
	return strings.TrimSuffix(humanize.Bytes(inodes), "B")
}

func printInfo(cmd *cobra.Command, info *fsquota.Info, prefix string) {
	cmd.Println(prefix + "bytes:")
	cmd.Printf(prefix+"  - soft: %s\n", humanize.IBytes(info.Bytes.GetSoft()))
	cmd.Printf(prefix+"  - hard: %s\n", humanize.IBytes(info.Bytes.GetHard()))
	cmd.Printf(prefix+"  - used: %s\n", humanize.IBytes(info.BytesUsed))
	cmd.Println(prefix + "files:")
	cmd.Printf(prefix+"  - soft: %s\n", humanizeInodes(info.Files.GetSoft()))
	cmd.Printf(prefix+"  - hard: %s\n", humanizeInodes(info.Files.GetHard()))
	cmd.Printf(prefix+"  - used: %s\n", humanizeInodes(info.FilesUsed))
}

func printQuotaInfo(cmd *cobra.Command, info *fsquota.Info) {
	printInfo(cmd, info, "")
}

func noopLookup(s string) string {
	return s
}

func printReport(cmd *cobra.Command, report *fsquota.Report, reportType string, lookupFn func(string) string) {
	for identifier, info := range report.Infos {
		cmd.Printf("%s %s:\n", reportType, lookupFn(identifier))
		printInfo(cmd, info, "  ")
	}
}

func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}

func parseLimitsFlag(cmd *cobra.Command, flagName string) (soft, hard uint64, present bool, err error) {
	var flagString string
	if flagString, err = cmd.Flags().GetString(flagName); err != nil {
		return
	}

	if flagString == "" {
		return
	}
	present = true

	valueParts := strings.Split(flagString, ",")
	if len(valueParts) != 2 {
		err = errors.New("expected format is soft,hard")
		return
	}

	var convErr error
	if soft, convErr = humanize.ParseBytes(valueParts[0]); convErr != nil {
		err = errortree.Add(err, "soft", convErr)
	}

	if hard, convErr = humanize.ParseBytes(valueParts[1]); convErr != nil {
		err = errortree.Add(err, "hard", convErr)
	}

	return
}
