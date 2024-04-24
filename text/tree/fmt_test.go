package tree_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/content/attrs"
	"github.com/notebox/nb-crdt-go/text/test"
	"github.com/notebox/nb-crdt-go/text/tree"
	"github.com/stretchr/testify/assert"
)

func TestNodeFMT(t *testing.T) {
	props := attrs.TextProps{"COL": "green"}

	t.Run("less", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			tree.New(*test.INSSpanFrom(test.Cases[common.Less]), nil, nil),
		)
		err := subject.FMT(test.FMTSpanFrom(test.Cases[common.Less], props), 0)
		assert.NoError(t, err)
		assert.Nil(t, subject.Left)
		assert.Equal(t, attrs.Attrs{{Length: 3, Props: props}}, subject.Right.Span.Content.Attrs())
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("prependable", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.Cases[common.Prependable], props), 0)
		assert.NoError(t, err)
		assert.Nil(t, subject.Left)
		assert.Nil(t, subject.Right)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("greater", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			tree.New(*test.INSSpanFrom(test.Cases[common.Greater]), nil, nil),
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.Cases[common.Greater], props), 0)
		assert.NoError(t, err)
		assert.Equal(t, attrs.Attrs{{Length: 3, Props: props}}, subject.Left.Span.Content.Attrs())
		assert.Nil(t, subject.Right)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("appendable", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.Cases[common.Appendable], props), 0)
		assert.NoError(t, err)
		assert.Nil(t, subject.Left)
		assert.Nil(t, subject.Right)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("includingLeft", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.RawCase{Point: test.Cases[common.IncludingLeft].Point, Text: "x"}, attrs.TextProps{"COL": "green"}), 0)
		assert.NoError(t, err)
		assert.Equal(t, attrs.Attrs{{Length: 1, Props: props}, {Length: 2}}, subject.Span.Content.Attrs())
	})

	t.Run("includingRight", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.RawCase{Point: test.Cases[common.IncludingRight].Point, Text: "x"}, attrs.TextProps{"COL": "green"}), 0)
		assert.NoError(t, err)
		assert.Equal(t, attrs.Attrs{{Length: 2}, {Length: 1, Props: props}}, subject.Span.Content.Attrs())
	})

	t.Run("includingMiddle", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.RawCase{Point: test.Cases[common.IncludingMiddle].Point, Text: "x"}, attrs.TextProps{"COL": "green"}), 0)
		assert.NoError(t, err)
		assert.Equal(t, attrs.Attrs{{Length: 1}, {Length: 1, Props: props}, {Length: 1}}, subject.Span.Content.Attrs())
	})

	t.Run("rightOverlap", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.RawCase{Point: test.Cases[common.RightOverlap].Point, Text: "xyz"}, attrs.TextProps{"COL": "green"}), 0)
		assert.NoError(t, err)
		assert.Equal(t, attrs.Attrs{{Length: 2}, {Length: 1, Props: props}}, subject.Span.Content.Attrs())
	})

	t.Run("leftOverlap", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.RawCase{Point: test.Cases[common.LeftOverlap].Point, Text: "xyz"}, attrs.TextProps{"COL": "green"}), 0)
		assert.NoError(t, err)
		assert.Equal(t, attrs.Attrs{{Length: 1, Props: props}, {Length: 2}}, subject.Span.Content.Attrs())
	})

	t.Run("splitted", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.RawCase{Point: test.Cases[common.Splitted].Point, Text: "xyz"}, attrs.TextProps{"COL": "green"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("splitting", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.RawCase{Point: test.Cases[common.Splitting].Point, Text: "xyz"}, attrs.TextProps{"COL": "green"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("includedLeft", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.RawCase{Point: test.Cases[common.IncludedLeft].Point, Text: "vwxy"}, attrs.TextProps{"COL": "green"}), 0)
		assert.NoError(t, err)
		assert.Equal(t, attrs.Attrs{{Length: 3, Props: props}}, subject.Span.Content.Attrs())
	})

	t.Run("includedMiddle", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.RawCase{Point: test.Cases[common.IncludedRight].Point, Text: "vwxyz"}, attrs.TextProps{"COL": "green"}), 0)
		assert.NoError(t, err)
		assert.Equal(t, attrs.Attrs{{Length: 3, Props: props}}, subject.Span.Content.Attrs())
	})

	t.Run("includedRight", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.RawCase{Point: test.Cases[common.IncludedRight].Point, Text: "wxyz"}, attrs.TextProps{"COL": "green"}), 0)
		assert.NoError(t, err)
		assert.Equal(t, attrs.Attrs{{Length: 3, Props: props}}, subject.Span.Content.Attrs())
	})

	t.Run("equal", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.FMT(test.FMTSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "xyz"}, attrs.TextProps{"COL": "green"}), 0)
		assert.NoError(t, err)
		assert.Equal(t, attrs.Attrs{{Length: 3, Props: props}}, subject.Span.Content.Attrs())
	})
}
