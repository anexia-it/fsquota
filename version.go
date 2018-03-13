package fsquota

import "fmt"

// Major version
const VersionMajor = 0

// Minor version
const VersionMinor = 1

// Patch version
const VersionPatch = 2

// VersionString returns the complete version string
func VersionString() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}
