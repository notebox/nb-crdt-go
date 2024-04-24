package content

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDELContent(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		stringifiedJSON := "5"
		var subject DELContent
		err := json.Unmarshal([]byte(stringifiedJSON), &subject)
		assert.NoError(t, err)
		assert.Equal(t, DELContent{length: 5}, subject)
		encoded, err := json.Marshal(subject)
		assert.NoError(t, err)
		assert.Equal(t, stringifiedJSON, string(encoded))
	})

	t.Run("Slice", func(t *testing.T) {
		subject := DELContent{length: 5}
		assert.True(t, subject.Slice(1, 3).Equals(&DELContent{length: 2}))
		assert.True(t, subject.Slice(5, 1).Equals(&DELContent{}))
		assert.True(t, subject.Slice(4, 9).Equals(&DELContent{length: 1}))
	})

	t.Run("Concat", func(t *testing.T) {
		a := DELContent{length: 3}
		b := DELContent{length: 5}
		assert.True(t, a.Concat(&b).Equals(&DELContent{length: 8}))
	})

	t.Run("NewDELContent", func(t *testing.T) {
		assert.True(t, NewDELContent(5).Equals(&DELContent{length: 5}))
	})
}
