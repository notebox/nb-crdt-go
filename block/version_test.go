package block_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/block"
	"github.com/stretchr/testify/assert"
)

func TestBlockVersion(t *testing.T) {
	newSubject := func() block.Version {
		return block.Version{59: block.ReplicaNonce{8, 7}, 1: block.ReplicaNonce{5, 9}, 2: block.ReplicaNonce{5, 9}, 3: block.ReplicaNonce{5, 9}}
	}

	t.Run("Merge", func(t *testing.T) {
		subject := newSubject()
		assert.True(t, subject.Merge(block.Version{59: block.ReplicaNonce{7, 6}, 1: block.ReplicaNonce{5, 10}, 2: block.ReplicaNonce{6, 9}, 19: block.ReplicaNonce{1, 1}}))
		assert.Equal(t, subject, block.Version{59: block.ReplicaNonce{8, 7}, 1: block.ReplicaNonce{5, 9}, 2: block.ReplicaNonce{6, 9}, 3: block.ReplicaNonce{5, 9}, 19: block.ReplicaNonce{1, 1}})
		assert.False(t, newSubject().Merge(block.Version{59: block.ReplicaNonce{8, 7}, 1: block.ReplicaNonce{5, 9}, 2: block.ReplicaNonce{5, 9}}))
	})

	t.Run("Add", func(t *testing.T) {
		subject := newSubject()
		assert.True(t, subject.Add(59, block.ReplicaNonce{9, 7}))
		assert.Equal(t, subject, block.Version{59: block.ReplicaNonce{9, 7}, 1: block.ReplicaNonce{5, 9}, 2: block.ReplicaNonce{5, 9}, 3: block.ReplicaNonce{5, 9}})
		assert.True(t, subject.Add(4, block.ReplicaNonce{1, 0}))
		assert.Equal(t, subject, block.Version{59: block.ReplicaNonce{9, 7}, 1: block.ReplicaNonce{5, 9}, 2: block.ReplicaNonce{5, 9}, 3: block.ReplicaNonce{5, 9}, 4: block.ReplicaNonce{1, 0}})
		assert.False(t, subject.Add(59, block.ReplicaNonce{8, 8}))
	})

	t.Run("IsNewerOrEqualThan", func(t *testing.T) {
		subject := newSubject()
		assert.True(t, subject.IsNewerOrEqualThan(nil))
		assert.True(t, subject.IsNewerOrEqualThan(subject))
		assert.True(t, subject.IsNewerOrEqualThan(block.Version{59: block.ReplicaNonce{8, 8}}))
		assert.True(t, subject.IsNewerOrEqualThan(block.Version{59: block.ReplicaNonce{7, 7}}))
		assert.True(t, subject.IsNewerOrEqualThan(block.Version{59: block.ReplicaNonce{7, 6}, 1: block.ReplicaNonce{5, 9}}))
		assert.False(t, subject.IsNewerOrEqualThan(block.Version{59: block.ReplicaNonce{9, 0}}))
		assert.False(t, subject.IsNewerOrEqualThan(block.Version{59: block.ReplicaNonce{7, 6}, 1: block.ReplicaNonce{6, 8}}))
		assert.False(t, subject.IsNewerOrEqualThan(block.Version{4: block.ReplicaNonce{1, 0}}))
	})
}
