package fsquota_test

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/anexia-it/fsquota"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareIntegrationTest(t *testing.T) string {
	testMountPoint := os.Getenv("TEST_MOUNTPOINT")

	if testMountPoint == "" {
		t.Skip("Skipping integration tests: TEST_MOUNTPOINT environment variable not set")
	} else if os.Getuid() != 0 {
		t.Skip("Skipping integration tests: not running as root")
	}

	// Bail out if the tests should be skipped
	if t.Skipped() {
		t.SkipNow()
	}

	testMountPoint = filepath.Join(testMountPoint, "child")
	require.NoError(t, os.MkdirAll(testMountPoint, 0755))
	return testMountPoint
}

func TestSetAndGetUserQuota(t *testing.T) {
	testMountPoint := prepareIntegrationTest(t)

	// Test against user 10000
	testUser := &user.User{
		Uid: "10000",
	}

	limits := fsquota.Limits{}
	limits.Bytes.SetSoft(10 * 1024 * 1024)  // 10MiB soft limit
	limits.Bytes.SetHard(500 * 1024 * 1024) // 500MiB hard limit
	limits.Files.SetSoft(1000)              // 1000 files soft limit
	limits.Files.SetHard(5000)              // 5000 files hard limit

	quotaInfo, err := fsquota.SetUserQuota(testMountPoint, testUser, limits)
	require.NoError(t, err)
	require.NotNil(t, quotaInfo)

	assert.EqualValues(t, 10*1024*1024, quotaInfo.Bytes.GetSoft())
	assert.EqualValues(t, 500*1024*1024, quotaInfo.Bytes.GetHard())
	assert.EqualValues(t, 1000, quotaInfo.Files.GetSoft())
	assert.EqualValues(t, 5000, quotaInfo.Files.GetHard())

	// Retrieve the quota information again, testing GetUserQuota as well
	quotaInfo, err = fsquota.GetUserInfo(testMountPoint, testUser)
	require.NoError(t, err)
	require.NotNil(t, quotaInfo)

	// The values should still be the same, meaning the information was persisted to the filesystem
	assert.EqualValues(t, 10*1024*1024, quotaInfo.Bytes.GetSoft())
	assert.EqualValues(t, 500*1024*1024, quotaInfo.Bytes.GetHard())
	assert.EqualValues(t, 1000, quotaInfo.Files.GetSoft())
	assert.EqualValues(t, 5000, quotaInfo.Files.GetHard())
}

func TestSetAndGetGroupQuota(t *testing.T) {
	testMountPoint := prepareIntegrationTest(t)

	// Test against group 10000
	testGroup := &user.Group{
		Gid: "10000",
	}

	limits := fsquota.Limits{}
	limits.Bytes.SetSoft(10 * 1024 * 1024)  // 10MiB soft limit
	limits.Bytes.SetHard(500 * 1024 * 1024) // 500MiB hard limit
	limits.Files.SetSoft(1000)              // 1000 files soft limit
	limits.Files.SetHard(5000)              // 5000 files hard limit

	quotaInfo, err := fsquota.SetGroupQuota(testMountPoint, testGroup, limits)
	require.NoError(t, err)

	assert.EqualValues(t, 10*1024*1024, quotaInfo.Bytes.GetSoft())
	assert.EqualValues(t, 500*1024*1024, quotaInfo.Bytes.GetHard())
	assert.EqualValues(t, 1000, quotaInfo.Files.GetSoft())
	assert.EqualValues(t, 5000, quotaInfo.Files.GetHard())

	// Retrieve the quota information again, testing GetUserQuota as well
	quotaInfo, err = fsquota.GetGroupInfo(testMountPoint, testGroup)
	require.NoError(t, err)
	require.NotNil(t, quotaInfo)

	// The values should still be the same, meaning the information was persisted to the filesystem
	assert.EqualValues(t, 10*1024*1024, quotaInfo.Bytes.GetSoft())
	assert.EqualValues(t, 500*1024*1024, quotaInfo.Bytes.GetHard())
	assert.EqualValues(t, 1000, quotaInfo.Files.GetSoft())
	assert.EqualValues(t, 5000, quotaInfo.Files.GetHard())
}

func TestUserQuotasSupported(t *testing.T) {
	testMountPoint := prepareIntegrationTest(t)

	supported, err := fsquota.UserQuotasSupported(testMountPoint)
	assert.True(t, supported)
	assert.NoError(t, err)

	// We expect / not to have quota support enabled
	supported, err = fsquota.UserQuotasSupported("/")
	assert.False(t, supported)
	assert.Error(t, err)
}

func TestGroupQuotasSupported(t *testing.T) {
	testMountPoint := prepareIntegrationTest(t)

	supported, err := fsquota.GroupQuotasSupported(testMountPoint)
	assert.True(t, supported)
	assert.NoError(t, err)

	// We expect / not to have quota support enabled
	supported, err = fsquota.GroupQuotasSupported("/")
	assert.False(t, supported)
	assert.Error(t, err)
}
