package common_test

import (
	"encoding/json"
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/stretchr/testify/assert"
)

func TestStamp(t *testing.T) {
	t.Run("CheckIsNewerThan", func(t *testing.T) {
		subject := &common.Stamp{ReplicaID: 1, Timestamp: 0}
		assert.True(t, subject.IsOlderThan(&common.Stamp{ReplicaID: 0, Timestamp: 1}))
		assert.False(t, subject.IsOlderThan(&common.Stamp{ReplicaID: 0, Timestamp: 0}))

		assert.False(t, subject.IsOlderThan(nil))
		subject = nil
		assert.True(t, subject.IsOlderThan(&common.Stamp{ReplicaID: 0, Timestamp: 0}))
	})

	t.Run("JSON", func(t *testing.T) {
		subject := &common.Stamp{ReplicaID: 1, Timestamp: 2}
		encoded, err := json.Marshal(subject)
		assert.NoError(t, err)
		assert.Equal(t, "[1,2]", string(encoded))

		var decoded common.Stamp
		json.Unmarshal(encoded, &decoded)
		assert.Equal(t, subject, &decoded)
	})
}
