// Package fsquota provides functions for working with filesystem quotas
package fsquota

import "os/user"

// SetUserQuota configures a user's quota
func SetUserQuota(path string, user *user.User, limits Limits) (info *Info, err error) {
	return setUserQuota(path, user, &limits)
}

// GetUserInfo retrieves a user's quota information
func GetUserInfo(path string, user *user.User) (info *Info, err error) {
	return getUserInfo(path, user)
}

// GetUserReport retrieves a report of all user quotas present at the given path
func GetUserReport(path string) (report *Report, err error) {
	return getUserReport(path)
}

// SetGroupQuota configures a group's quota
func SetGroupQuota(path string, group *user.Group, limits Limits) (info *Info, err error) {
	return setGroupQuota(path, group, &limits)
}

// GetGroupInfo retrieves a group's quota information
func GetGroupInfo(path string, group *user.Group) (info *Info, err error) {
	return getGroupInfo(path, group)
}

// GetGroupReport retrieves a report of all group quotas present at the given path
func GetGroupReport(path string) (report *Report, err error) {
	return getGroupReport(path)
}

// UserQuotasSupported checks if quotas are supported on a given path
func UserQuotasSupported(path string) (supported bool, err error) {
	return userQuotasSupported(path)
}

// GroupQuotasSupported checks if group quotas are supported on a given path
func GroupQuotasSupported(path string) (supported bool, err error) {
	return groupQuotasSupported(path)
}
