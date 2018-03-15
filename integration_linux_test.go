package fsquota_test

import (
	"os"
	"os/user"
	"testing"

	"github.com/anexia-it/fsquota"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetUserQuota(t *testing.T) {
	testMountPoint := os.Getenv("TEST_MOUNTPOINT")

	if testMountPoint == "" {
		t.Skip("Skipping integration tests: TEST_MOUNTPOINT environment variable not set")
	} else if os.Getuid() != 0 {
		t.Skip("Skipping integration tests: not running as root")
	}

	// Bail out if the tests should be skipped
	if t.Skipped() {
		return
	}

	// Test against user 10000
	testUser := &user.User{
		Uid: "10000",
	}

	limits := fsquota.Limits{}
	limits.Bytes.SetSoft(10 * 1024 * 1024) // 10MiB soft limit
	limits.Bytes.SetHard(500 * 1024 * 1024) // 500MiB hard limit
	limits.Files.SetSoft(1000) // 1000 files soft limit
	limits.Files.SetHard(5000) // 5000 files hard limit

	quotaInfo, err := fsquota.SetUserQuota(testMountPoint, testUser, limits)
	require.NoError(t, err)

	assert.EqualValues(t, 10 * 1024 * 1024, quotaInfo.Bytes.GetSoft())
	assert.EqualValues(t, 500 * 1024 * 1024, quotaInfo.Bytes.GetHard())
	assert.EqualValues(t, 1000, quotaInfo.Files.GetSoft())
	assert.EqualValues(t, 5000, quotaInfo.Files.GetHard())
}