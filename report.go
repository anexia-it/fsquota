package fsquota

// Report contains a quota report
type Report struct {
	// Map of user or group to info structure
	Infos map[string]*Info
}
