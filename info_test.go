package fsquota

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoIsEmpty(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		i := Info{}
		assert.True(t, i.isEmpty())
	})

	t.Run("BytesSoftSet", func(t *testing.T) {
		l := uint64(1)
		i := Info{
			Limits: Limits{
				Bytes: Limit{
					soft: &l,
				},
			},
		}
		assert.False(t, i.isEmpty())
	})

	t.Run("BytesHardSet", func(t *testing.T) {
		l := uint64(1)
		i := Info{
			Limits: Limits{
				Bytes: Limit{
					hard: &l,
				},
			},
		}
		assert.False(t, i.isEmpty())
	})

	t.Run("BytesUsedSet", func(t *testing.T) {
		i := Info{
			BytesUsed: 1,
		}
		assert.False(t, i.isEmpty())
	})

	t.Run("FilesSoftSet", func(t *testing.T) {
		l := uint64(1)
		i := Info{
			Limits: Limits{
				Files: Limit{
					soft: &l,
				},
			},
		}
		assert.False(t, i.isEmpty())
	})

	t.Run("FilesHardSet", func(t *testing.T) {
		l := uint64(1)
		i := Info{
			Limits: Limits{
				Files: Limit{
					hard: &l,
				},
			},
		}
		assert.False(t, i.isEmpty())
	})

	t.Run("FilesUsedSet", func(t *testing.T) {
		i := Info{
			FilesUsed: 1,
		}
		assert.False(t, i.isEmpty())
	})
}
