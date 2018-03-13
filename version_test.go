package fsquota

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionString(t *testing.T) {
	assert.EqualValues(t, fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch), VersionString())
}
