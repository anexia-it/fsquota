package fsquota_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/anexia-it/fsquota"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareIntegrationTest(t *testing.T) (testMountpointQuotasEnabled, testMountpointQuotasDisabled string) {
	testMountpointQuotasEnabled = os.Getenv("TEST_MOUNTPOINT_QUOTAS_ENABLED")
	testMountpointQuotasDisabled = os.Getenv("TEST_MOUNTPOINT_QUOTAS_DISABLED")

	if testMountpointQuotasEnabled == "" {
		t.Skip("Skipping integration tests: TEST_MOUNTPOINT_QUOTAS_ENABLED environment variable not set")
	} else if testMountpointQuotasDisabled == "" {
		t.Skip("Skipping integration tests: TEST_MOUNTPOINT_QUOTAS_DISABLED environment variable not set")
	} else if os.Getuid() != 0 {
		t.Skip("Skipping integration tests: not running as root")
	}

	// Bail out if the tests should be skipped
	if t.Skipped() {
		t.SkipNow()
	}

	testMountpointQuotasEnabled = filepath.Join(testMountpointQuotasEnabled, "child")
	require.NoError(t, os.MkdirAll(testMountpointQuotasEnabled, 0755))

	testMountpointQuotasDisabled = filepath.Join(testMountpointQuotasDisabled, "child")
	require.NoError(t, os.MkdirAll(testMountpointQuotasDisabled, 0755))
	return
}

func TestSetAndGetUserQuota(t *testing.T) {
	testMountPointQuotasEnabled, testMountpointQuotasDisabled := prepareIntegrationTest(t)

	// Test against user 10000
	testUser := &user.User{
		Uid: "10000",
	}

	limits := fsquota.Limits{}
	limits.Bytes.SetSoft(10 * 1024 * 1024)  // 10MiB soft limit
	limits.Bytes.SetHard(500 * 1024 * 1024) // 500MiB hard limit
	limits.Files.SetSoft(1000)              // 1000 files soft limit
	limits.Files.SetHard(5000)              // 5000 files hard limit

	t.Run("QuotasEnabled", func(t *testing.T) {
		quotaInfo, err := fsquota.SetUserQuota(testMountPointQuotasEnabled, testUser, limits)
		require.NoError(t, err)
		require.NotNil(t, quotaInfo)

		assert.EqualValues(t, 10*1024*1024, quotaInfo.Bytes.GetSoft())
		assert.EqualValues(t, 500*1024*1024, quotaInfo.Bytes.GetHard())
		assert.EqualValues(t, 1000, quotaInfo.Files.GetSoft())
		assert.EqualValues(t, 5000, quotaInfo.Files.GetHard())

		// Retrieve the quota information again, testing GetUserQuota as well
		quotaInfo, err = fsquota.GetUserInfo(testMountPointQuotasEnabled, testUser)
		require.NoError(t, err)
		require.NotNil(t, quotaInfo)

		// The values should still be the same, meaning the information was persisted to the filesystem
		assert.EqualValues(t, 10*1024*1024, quotaInfo.Bytes.GetSoft())
		assert.EqualValues(t, 500*1024*1024, quotaInfo.Bytes.GetHard())
		assert.EqualValues(t, 1000, quotaInfo.Files.GetSoft())
		assert.EqualValues(t, 5000, quotaInfo.Files.GetHard())
	})

	t.Run("QuotasDisabled", func(t *testing.T) {
		quotaInfo, err := fsquota.SetUserQuota(testMountpointQuotasDisabled, testUser, limits)
		assert.Error(t, err)
		assert.Nil(t, quotaInfo)
	})
}

func TestSetAndGetGroupQuota(t *testing.T) {
	testMountPointQuotasEnabled, testMountpointQuotasDisabled := prepareIntegrationTest(t)

	// Test against group 10000
	testGroup := &user.Group{
		Gid: "10000",
	}

	limits := fsquota.Limits{}
	limits.Bytes.SetSoft(10 * 1024 * 1024)  // 10MiB soft limit
	limits.Bytes.SetHard(500 * 1024 * 1024) // 500MiB hard limit
	limits.Files.SetSoft(1000)              // 1000 files soft limit
	limits.Files.SetHard(5000)              // 5000 files hard limit

	t.Run("QuotasEnabled", func(t *testing.T) {
		quotaInfo, err := fsquota.SetGroupQuota(testMountPointQuotasEnabled, testGroup, limits)
		require.NoError(t, err)

		assert.EqualValues(t, 10*1024*1024, quotaInfo.Bytes.GetSoft())
		assert.EqualValues(t, 500*1024*1024, quotaInfo.Bytes.GetHard())
		assert.EqualValues(t, 1000, quotaInfo.Files.GetSoft())
		assert.EqualValues(t, 5000, quotaInfo.Files.GetHard())

		// Retrieve the quota information again, testing GetUserQuota as well
		quotaInfo, err = fsquota.GetGroupInfo(testMountPointQuotasEnabled, testGroup)
		require.NoError(t, err)
		require.NotNil(t, quotaInfo)

		// The values should still be the same, meaning the information was persisted to the filesystem
		assert.EqualValues(t, 10*1024*1024, quotaInfo.Bytes.GetSoft())
		assert.EqualValues(t, 500*1024*1024, quotaInfo.Bytes.GetHard())
		assert.EqualValues(t, 1000, quotaInfo.Files.GetSoft())
		assert.EqualValues(t, 5000, quotaInfo.Files.GetHard())
	})

	t.Run("QuotasDisabled", func(t *testing.T) {
		quotaInfo, err := fsquota.SetGroupQuota(testMountpointQuotasDisabled, testGroup, limits)
		assert.Error(t, err)
		assert.Nil(t, quotaInfo)
	})

}

func TestUserQuotasSupported(t *testing.T) {
	testMountPointQuotasEnabled, testMountpointQuotasDisabled := prepareIntegrationTest(t)

	t.Run("QuotasEnabled", func(t *testing.T) {
		supported, err := fsquota.UserQuotasSupported(testMountPointQuotasEnabled)
		assert.True(t, supported)
		assert.NoError(t, err)
	})

	t.Run("QuotasDisabled", func(t *testing.T) {
		supported, err := fsquota.UserQuotasSupported(testMountpointQuotasDisabled)
		assert.False(t, supported)
		assert.Error(t, err)
	})
}

func TestGroupQuotasSupported(t *testing.T) {
	testMountPointQuotasEnabled, testMountpointQuotasDisabled := prepareIntegrationTest(t)

	t.Run("QuotasEnabled", func(t *testing.T) {
		supported, err := fsquota.GroupQuotasSupported(testMountPointQuotasEnabled)
		assert.True(t, supported)
		assert.NoError(t, err)
	})

	t.Run("QuotasDisabled", func(t *testing.T) {
		supported, err := fsquota.GroupQuotasSupported(testMountpointQuotasDisabled)
		assert.False(t, supported)
		assert.Error(t, err)
	})
}

func TestGetUserReport(t *testing.T) {
	testMountPointQuotasEnabled, testMountPointQuotasDisabled := prepareIntegrationTest(t)

	const userCount = 10
	const uidBase = 10000

	limits := fsquota.Limits{}
	limits.Bytes.SetSoft(userCount * 1024 * 1024) // userCount MiB soft limit
	limits.Bytes.SetHard(500 * 1024 * 1024)       // 500MiB hard limit
	limits.Files.SetSoft(1000)                    // 1000 files soft limit
	limits.Files.SetHard(5000)                    // 5000 files hard limit

	t.Run("QuotasEnabled", func(t *testing.T) {
		expectedUserByteUsages := make(map[int]int64, userCount)

		testDirectory := filepath.Join(testMountPointQuotasEnabled, "data_files")
		require.NoError(t, os.MkdirAll(testDirectory, 0755))
		defer os.RemoveAll(testDirectory)

		for uid := uidBase; uid < uidBase+userCount; uid++ {
			// Test against user 10000
			testUser := &user.User{
				Uid: fmt.Sprint(uid),
			}

			// Configure quota for user
			quotaInfo, err := fsquota.SetUserQuota(testMountPointQuotasEnabled, testUser, limits)
			require.NoError(t, err, "Failed to set quota for UID %d", uid)

			assert.EqualValues(t, 10*1024*1024, quotaInfo.Bytes.GetSoft())
			assert.EqualValues(t, 500*1024*1024, quotaInfo.Bytes.GetHard())
			assert.EqualValues(t, 1000, quotaInfo.Files.GetSoft())
			assert.EqualValues(t, 5000, quotaInfo.Files.GetHard())

			// Create s sparse file for the user...
			fileName := filepath.Join(testDirectory, fmt.Sprintf("user_%d.data", uid))
			fileSize := int64(((uid - uidBase) + 1) * 1024 * 1024) // user 0 gets 1MiB, user 1 gets 2MiB file, ...
			expectedUserByteUsages[uid] = fileSize

			// Run creation logic in an anonymous function so we can defer close the file...
			func() {
				f, err := os.Create(fileName)
				require.NoError(t, err)
				defer f.Close()

				// Create data file...
				data := make([]byte, fileSize)
				_, err = f.Write(data)
				require.NoError(t, err)

				// Change ownership of the file to the target user
				require.NoError(t, os.Chown(fileName, uid, uid))
			}()

			// Create additional files for every user...
			for i := uidBase; i < uid; i++ {
				fileName := filepath.Join(testDirectory, fmt.Sprintf("user_%d.empty.%d", uid, i))
				require.NoError(t, ioutil.WriteFile(fileName, []byte{}, 0644), "Failed to write user %d file #%d", uid, i-uidBase)
				require.NoError(t, os.Chown(fileName, uid, uid))
			}
		}

		// Now generate the report
		report, err := fsquota.GetUserReport(testMountPointQuotasEnabled)
		require.NoError(t, err)
		require.NotNil(t, report)

		// Check that the report for every user is correct...
		for uid, expectedBytes := range expectedUserByteUsages {
			userReport, userReportExists := report.Infos[fmt.Sprint(uid)]

			userReportFound := assert.True(t, userReportExists, "Report info for UID %d is missing", uid)
			userReportNotNil := assert.NotNil(t, userReport, "Report for UID %d is nil", uid)

			if userReportFound && userReportNotNil {
				assert.EqualValues(t, expectedBytes, userReport.BytesUsed, "Bytes used for UID %d is invalid. Expected %d, actual %d", uid, expectedBytes, userReport.BytesUsed)
				assert.EqualValues(t, 1+(uid-uidBase), userReport.FilesUsed, "Files used for UID %d is invalid. Expected %d, actual %d", uid, 1+(uid-uidBase), userReport.FilesUsed)

				// Validate that the configured quotas are correct
				assert.EqualValues(t, 10*1024*1024, userReport.Bytes.GetSoft(), "Byte soft quota for UID %d is invalid", uid)
				assert.EqualValues(t, 500*1024*1024, userReport.Bytes.GetHard(), "Byte hard quota for UID %d is invalid", uid)
				assert.EqualValues(t, 1000, userReport.Files.GetSoft(), "File soft quota for UID %d is invalid", uid)
				assert.EqualValues(t, 5000, userReport.Files.GetHard(), "File hard quota for UID %d is invalid", uid)
			}

		}
	})

	t.Run("QuotasDisabled", func(t *testing.T) {
		report, err := fsquota.GetUserReport(testMountPointQuotasDisabled)
		assert.Error(t, err)
		assert.Nil(t, report)
	})
}

func TestGetGroupReport(t *testing.T) {
	testMountPointQuotasEnabled, testMountPointQuotasDisabled := prepareIntegrationTest(t)

	const groupCount = 10
	const gidBase = 10000

	limits := fsquota.Limits{}
	limits.Bytes.SetSoft(groupCount * 1024 * 1024) // groupCount MiB soft limit
	limits.Bytes.SetHard(500 * 1024 * 1024)        // 500MiB hard limit
	limits.Files.SetSoft(1000)                     // 1000 files soft limit
	limits.Files.SetHard(5000)                     // 5000 files hard limit

	t.Run("QuotasEnabled", func(t *testing.T) {
		expectedGroupByteUsages := make(map[int]int64, groupCount)

		testDirectory := filepath.Join(testMountPointQuotasEnabled, "data_files")
		require.NoError(t, os.MkdirAll(testDirectory, 0755))
		defer os.RemoveAll(testDirectory)

		for gid := gidBase; gid < gidBase+groupCount; gid++ {
			// Test against group
			testGroup := &user.Group{
				Gid: fmt.Sprint(gid),
			}

			// Configure quota for group
			quotaInfo, err := fsquota.SetGroupQuota(testMountPointQuotasEnabled, testGroup, limits)
			require.NoError(t, err, "Failed to set quota for GID %d", gid)

			assert.EqualValues(t, 10*1024*1024, quotaInfo.Bytes.GetSoft())
			assert.EqualValues(t, 500*1024*1024, quotaInfo.Bytes.GetHard())
			assert.EqualValues(t, 1000, quotaInfo.Files.GetSoft())
			assert.EqualValues(t, 5000, quotaInfo.Files.GetHard())

			// Create s sparse file for the user...
			fileName := filepath.Join(testDirectory, fmt.Sprintf("group_%d.data", gid))
			fileSize := int64(((gid - gidBase) + 1) * 1024 * 1024) // group 0 gets 1MiB, group 1 gets 2MiB file, ...
			expectedGroupByteUsages[gid] = fileSize

			// Run creation logic in an anonymous function so we can defer close the file...
			func() {
				f, err := os.Create(fileName)
				require.NoError(t, err)
				defer f.Close()

				// Create data file...
				data := make([]byte, fileSize)
				_, err = f.Write(data)
				require.NoError(t, err)

				// Change ownership of the file to the target group
				require.NoError(t, os.Chown(fileName, gid, gid))
			}()

			// Create additional files for every user...
			for i := gidBase; i < gid; i++ {
				fileName := filepath.Join(testDirectory, fmt.Sprintf("group_%d.empty.%d", gid, i))
				require.NoError(t, ioutil.WriteFile(fileName, []byte{}, 0644), "Failed to write group %d file #%d", gid, i-gidBase)
				require.NoError(t, os.Chown(fileName, gid, gid))
			}
		}

		// Now generate the report
		report, err := fsquota.GetGroupReport(testMountPointQuotasEnabled)
		require.NoError(t, err)
		require.NotNil(t, report)

		// Check that the report for every user is correct...
		for gid, expectedBytes := range expectedGroupByteUsages {
			groupReport, groupReportExists := report.Infos[fmt.Sprint(gid)]

			groupReportFound := assert.True(t, groupReportExists, "Report info for GID %d is missing", gid)
			groupReportNotNil := assert.NotNil(t, groupReport, "Report for GID %d is nil", gid)

			if groupReportFound && groupReportNotNil {
				assert.EqualValues(t, expectedBytes, groupReport.BytesUsed, "Bytes used for GID %d is invalid. Expected %d, actual %d", gid, expectedBytes, groupReport.BytesUsed)
				assert.EqualValues(t, 1+(gid-gidBase), groupReport.FilesUsed, "Files used for GID %d is invalid. Expected %d, actual %d", gid, 1+(gid-gidBase), groupReport.FilesUsed)

				// Validate that the configured quotas are correct
				assert.EqualValues(t, 10*1024*1024, groupReport.Bytes.GetSoft(), "Byte soft quota for GID %d is invalid", gid)
				assert.EqualValues(t, 500*1024*1024, groupReport.Bytes.GetHard(), "Byte hard quota for GID %d is invalid", gid)
				assert.EqualValues(t, 1000, groupReport.Files.GetSoft(), "File soft quota for GID %d is invalid", gid)
				assert.EqualValues(t, 5000, groupReport.Files.GetHard(), "File hard quota for GID %d is invalid", gid)
			}

		}
	})

	t.Run("QuotasDisabled", func(t *testing.T) {
		report, err := fsquota.GetGroupReport(testMountPointQuotasDisabled)
		assert.Error(t, err)
		assert.Nil(t, report)
	})
}
