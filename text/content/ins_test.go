package content

import (
	"encoding/json"
	"testing"

	"github.com/notebox/nb-crdt-go/text/content/attrs"
	"github.com/stretchr/testify/assert"
)

func TestINSContent(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		stringifiedJSON := `[[[6]],"foobar"]`
		var subject INSContent
		err := json.Unmarshal([]byte(stringifiedJSON), &subject)
		assert.NoError(t, err)
		assert.Equal(t, INSContent{text: "foobar", attrs: attrs.Attrs{{Length: 6}}}, subject)
		encoded, err := json.Marshal(subject)
		assert.NoError(t, err)
		assert.Equal(t, stringifiedJSON, string(encoded))
	})

	t.Run("Slice", func(t *testing.T) {
		subject := INSContent{text: "foobar", attrs: attrs.Attrs{{Length: 6}}}
		assert.True(t, subject.Slice(1, 3).Equals(&INSContent{text: "oo", attrs: attrs.Attrs{{Length: 2}}}))
	})

	t.Run("Concat", func(t *testing.T) {
		a := INSContent{text: "foo", attrs: attrs.Attrs{{Length: 3}}}
		b := INSContent{text: "bar", attrs: attrs.Attrs{{Length: 3}}}
		assert.True(t, a.Concat(&b).Equals(&INSContent{text: "foobar", attrs: attrs.Attrs{{Length: 6}}}))
	})

	t.Run("NewINSContent", func(t *testing.T) {
		assert.True(t, NewINSContent("foobar").Equals(&INSContent{text: "foobar", attrs: attrs.Attrs{{Length: 6}}}))
	})

	t.Run("others", func(t *testing.T) {
		subject := INSContent{text: "aa😀bb🤖cc👍🏿dd강ee", attrs: attrs.Attrs{{Length: 19}}}
		assert.Equal(t, uint32(19), subject.Length())
		assert.Equal(t, attrs.Attrs{{Length: 19}}, subject.Attrs())
	})

	t.Run("FMT", func(t *testing.T) {
		subject := INSContent{text: "foobar", attrs: attrs.Attrs{{Length: 6}}}
		subject.FMT(2, &FMTContent{length: 3, attrs: attrs.Attrs{{Length: 3, Props: attrs.TextProps{"B": true}}}})
		assert.True(t, subject.Equals(&INSContent{text: "foobar", attrs: attrs.Attrs{{Length: 2}, {Length: 3, Props: attrs.TextProps{"B": true}}, {Length: 1}}}))
	})

	t.Run("MOD", func(t *testing.T) {
		subject := INSContent{text: "foobar", attrs: attrs.Attrs{{Length: 6}}}
		subject.MOD(2, &MODContent{text: "xyz"})
		assert.True(t, subject.Equals(&INSContent{text: "foxyzr", attrs: attrs.Attrs{{Length: 6}}}))
	})
}
