package content

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMODContent(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		stringifiedJSON := `"foobar"`
		var subject MODContent
		err := json.Unmarshal([]byte(stringifiedJSON), &subject)
		assert.NoError(t, err)
		assert.Equal(t, MODContent{text: "foobar"}, subject)
		encoded, err := json.Marshal(subject)
		assert.NoError(t, err)
		assert.Equal(t, stringifiedJSON, string(encoded))
	})

	t.Run("Slice", func(t *testing.T) {
		subject := MODContent{text: "foobar"}
		assert.True(t, subject.Slice(1, 3).Equals(&MODContent{text: "oo"}))
	})

	t.Run("Concat", func(t *testing.T) {
		a := MODContent{text: "foo"}
		b := MODContent{text: "bar"}
		assert.True(t, a.Concat(&b).Equals(&MODContent{text: "foobar"}))
	})

	t.Run("NewMODContent", func(t *testing.T) {
		assert.True(t, NewMODContent("foobar").Equals(&MODContent{text: "foobar"}))
	})

	t.Run("others", func(t *testing.T) {
		subject := MODContent{text: "aa😀bb🤖cc👍🏿dd강ee"}
		assert.Equal(t, uint32(19), subject.Length())
	})
}
