package fsquota

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDqpblk_FromLimits(t *testing.T) {
	t.Run("InfoOnly", func(t *testing.T) {
		l := &Limits{}
		dq := dqblkFromLimits(l)
		assert.EqualValues(t, 0, dq.dqbValid)
	})

	t.Run("FileLimits", func(t *testing.T) {
		l := &Limits{}
		l.Files.SetHard(1024)
		l.Files.SetSoft(1000)

		dq := dqblkFromLimits(l)
		assert.EqualValues(t, qifILimits, dq.dqbValid)
		assert.EqualValues(t, 1024, dq.dqbIHardlimit)
		assert.EqualValues(t, 1000, dq.dqbISoftlimit)
	})

	t.Run("ByteLimits", func(t *testing.T) {
		l := &Limits{}
		l.Bytes.SetHard(2048)
		l.Bytes.SetSoft(1024)

		dq := dqblkFromLimits(l)
		assert.EqualValues(t, qifBLimits, dq.dqbValid)
		assert.EqualValues(t, 2, dq.dqbBHardlimit)
		assert.EqualValues(t, 1, dq.dqbBSoftlimit)
	})

	t.Run("Combined", func(t *testing.T) {
		l := &Limits{}
		l.Bytes.SetHard(2048)
		l.Bytes.SetSoft(1024)
		l.Files.SetHard(1024)
		l.Files.SetSoft(1000)

		dq := dqblkFromLimits(l)
		assert.EqualValues(t, qifBLimits|qifILimits, dq.dqbValid)
		assert.EqualValues(t, 2, dq.dqbBHardlimit)
		assert.EqualValues(t, 1, dq.dqbBSoftlimit)
		assert.EqualValues(t, 2, dq.dqbBHardlimit)
		assert.EqualValues(t, 1, dq.dqbBSoftlimit)
	})
}

func TestDqpblk_ToInfo(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		dq := &dqblk{}
		info := dq.toInfo()
		require.NotNil(t, info)
		assert.EqualValues(t, 0, info.FilesUsed)
		assert.EqualValues(t, 0, info.BytesUsed)

		require.NotNil(t, info.Files.soft)
		require.NotNil(t, info.Files.hard)
		assert.EqualValues(t, 0, *info.Files.soft)
		assert.EqualValues(t, 0, *info.Files.hard)

		require.NotNil(t, info.Bytes.soft)
		require.NotNil(t, info.Bytes.hard)

		assert.EqualValues(t, 0, *info.Bytes.soft)
		assert.EqualValues(t, 0, *info.Bytes.hard)
	})
}
