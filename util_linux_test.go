package fsquota

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIDsFromUserOrGroupFile(t *testing.T) {
	t.Run("FileNotFound", func(t *testing.T) {
		dirName, err := ioutil.TempDir("", "fsquota-test-")
		require.NoError(t, err)
		defer os.RemoveAll(dirName)

		ids, err := getIDsFromUserOrGroupFile(filepath.Join(dirName, "non-existent"))
		assert.Nil(t, ids)
		if assert.Error(t, err) {
			assert.True(t, os.IsNotExist(err))
		}
	})

	t.Run("OK", func(t *testing.T) {
		dirName, err := ioutil.TempDir("", "fsquota-test-")
		require.NoError(t, err)
		defer os.RemoveAll(dirName)

		fileData := `
# line without colon
too:few:parts
ignored:ignored:unparsable:ignored
ignored:ignored:1:ok1
ignored:ignored:2:ok2
ignored:ignored:1000:ok1000
`

		fileName := filepath.Join(dirName, "ids")
		require.NoError(t, ioutil.WriteFile(fileName, []byte(fileData), 0640))

		ids, err := getIDsFromUserOrGroupFile(fileName)
		assert.NoError(t, err)
		assert.EqualValues(t, []uint32{1, 2, 1000}, ids)

	})
}
