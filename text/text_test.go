package text_test

import (
	"encoding/json"
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text"
	"github.com/notebox/nb-crdt-go/text/content/attrs"
	"github.com/notebox/nb-crdt-go/text/span"
	"github.com/notebox/nb-crdt-go/text/test"
	"github.com/notebox/nb-crdt-go/text/tree"
	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		stringifiedJSON := "[[[[1,2,3]],[[[6]],\"foobar\"]],[[[4,5,6]],[[[4]],\"kang\"]]]"
		var subject text.Text
		err := json.Unmarshal([]byte(stringifiedJSON), &subject)
		assert.NoError(t, err)
		assert.Equal(t, text.Text{Node: tree.NewFromSpans([]*span.INSSpan{
			test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{1, 2, 3}}, Text: "foobar"}),
			test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{4, 5, 6}}, Text: "kang"}),
		})}, subject)
		encoded, err := json.Marshal(subject)
		assert.NoError(t, err)
		assert.Equal(t, stringifiedJSON, string(encoded))
	})

	t.Run("JSON Empty Node", func(t *testing.T) {
		stringifiedJSON := "[]"
		var subject text.Text
		err := json.Unmarshal([]byte(stringifiedJSON), &subject)
		assert.NoError(t, err)
		assert.Equal(t, text.Text{}, subject)
		encoded, err := json.Marshal(subject)
		assert.NoError(t, err)
		assert.Equal(t, stringifiedJSON, string(encoded))
	})

	t.Run("String, Spans", func(t *testing.T) {
		ls := test.INSSpanFrom(test.Cases[common.Greater])
		s := test.INSSpanFrom(test.Cases[common.Equal])
		rs := test.INSSpanFrom(test.Cases[common.Less])
		subject := text.Text{Node: tree.NewFromSpans([]*span.INSSpan{ls, s, rs})}

		assert.Equal(t, "01245689a", subject.String())

		spans := subject.Spans()
		assert.Equal(t, []*span.INSSpan{ls, s, rs}, spans)
	})

	t.Run("INS", func(t *testing.T) {
		subject := text.Text{}

		subject.INS(test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 3}}, Text: "a"}))
		assert.Equal(t, "a", subject.String())

		subject.INS(test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 4}}, Text: "b"}))
		assert.Equal(t, "ab", subject.String())
	})

	t.Run("DEL", func(t *testing.T) {
		subject := text.Text{Node: tree.NewFromSpans([]*span.INSSpan{
			test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 1}}, Text: "a"}),
			test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 3}}, Text: "b"}),
			test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 5}}, Text: "c"}),
		})}

		subject.DEL(test.DELSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 3}}, Text: "b"}))
		assert.Equal(t, "ac", subject.String())
	})

	t.Run("MOD", func(t *testing.T) {
		subject := text.Text{Node: tree.NewFromSpans([]*span.INSSpan{
			test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 1}}, Text: "a"}),
			test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 3}}, Text: "b"}),
			test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 5}}, Text: "c"}),
		})}

		subject.MOD(test.MODSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 3}}, Text: "x"}))
		assert.Equal(t, "axc", subject.String())
	})

	t.Run("FMT", func(t *testing.T) {
		subject := text.Text{Node: tree.NewFromSpans([]*span.INSSpan{
			test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 1}}, Text: "a"}),
			test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 3}}, Text: "b"}),
			test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 5}}, Text: "c"}),
		})}

		subject.FMT(test.FMTSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 3}}, Text: "x"}, attrs.TextProps{"B": true}))
		encoded, err := subject.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `[[[[5,5,1]],[[[1]],"a"]],[[[5,5,3]],[[[1,{"B":true}]],"b"]],[[[5,5,5]],[[[1]],"c"]]]`, string(encoded))
	})
}
