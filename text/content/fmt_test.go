package content

import (
	"encoding/json"
	"testing"

	"github.com/notebox/nb-crdt-go/text/content/attrs"
	"github.com/stretchr/testify/assert"
)

func TestFMTContent(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		stringifiedJSON := "[5,[[5]]]"
		var subject FMTContent
		err := json.Unmarshal([]byte(stringifiedJSON), &subject)
		assert.NoError(t, err)
		assert.Equal(t, FMTContent{length: 5, attrs: attrs.Attrs{{Length: 5}}}, subject)
		encoded, err := json.Marshal(subject)
		assert.NoError(t, err)
		assert.Equal(t, stringifiedJSON, string(encoded))
	})

	t.Run("Slice", func(t *testing.T) {
		subject := FMTContent{length: 5, attrs: attrs.Attrs{{Length: 5}}}
		assert.True(t, subject.Slice(1, 3).Equals(&FMTContent{length: 2, attrs: attrs.Attrs{{Length: 2}}}))
		assert.True(t, subject.Slice(5, 1).Equals(&FMTContent{}))
		assert.True(t, subject.Slice(4, 9).Equals(&FMTContent{length: 1, attrs: attrs.Attrs{{Length: 1}}}))
	})

	t.Run("Concat", func(t *testing.T) {
		a := FMTContent{length: 3, attrs: attrs.Attrs{{Length: 3}}}
		b := FMTContent{length: 5, attrs: attrs.Attrs{{Length: 5}}}
		assert.True(t, a.Concat(&b).Equals(&FMTContent{length: 8, attrs: attrs.Attrs{{Length: 8}}}))
	})

	t.Run("NewFMTContent", func(t *testing.T) {
		assert.True(t, NewFMTContent(5, attrs.TextProps{"B": true}).Equals(&FMTContent{length: 5, attrs: attrs.Attrs{{Length: 5, Props: attrs.TextProps{"B": true}}}}))
	})
}
