package tree_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/test"
	"github.com/notebox/nb-crdt-go/text/tree"
	"github.com/stretchr/testify/assert"
)

func TestNodeMOD(t *testing.T) {
	t.Run("less", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			tree.New(*test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Less].Point, Text: "000"}), nil, nil),
		)
		err := subject.MOD(test.MODSpanFrom(test.Cases[common.Less]), 0)
		assert.NoError(t, err)
		assert.Nil(t, subject.Left)
		assert.True(t, subject.Right.Span.Equals(test.INSSpanFrom(test.Cases[common.Less])))
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("prependable", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.Cases[common.Prependable]), 0)
		assert.NoError(t, err)
		assert.Nil(t, subject.Left)
		assert.Nil(t, subject.Right)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("greater", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			tree.New(*test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Greater].Point, Text: "000"}), nil, nil),
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.Cases[common.Greater]), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Left.Span.Equals(test.INSSpanFrom(test.Cases[common.Greater])))
		assert.Nil(t, subject.Right)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("appendable", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.Cases[common.Appendable]), 0)
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
		err := subject.MOD(test.MODSpanFrom(test.RawCase{Point: test.Cases[common.IncludingLeft].Point, Text: "x"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "x56"})))
	})

	t.Run("includingRight", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.RawCase{Point: test.Cases[common.IncludingRight].Point, Text: "x"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "45x"})))
	})

	t.Run("includingMiddle", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.RawCase{Point: test.Cases[common.IncludingMiddle].Point, Text: "x"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "4x6"})))
	})

	t.Run("rightOverlap", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.RawCase{Point: test.Cases[common.RightOverlap].Point, Text: "xyz"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "45x"})))
	})

	t.Run("leftOverlap", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.RawCase{Point: test.Cases[common.LeftOverlap].Point, Text: "xyz"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "z56"})))
	})

	t.Run("splitted", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.RawCase{Point: test.Cases[common.Splitted].Point, Text: "xyz"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("splitting", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.RawCase{Point: test.Cases[common.Splitting].Point, Text: "xyz"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("includedLeft", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.RawCase{Point: test.Cases[common.IncludedLeft].Point, Text: "vwxy"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "vwx"})))
	})

	t.Run("includedMiddle", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.RawCase{Point: test.Cases[common.IncludedRight].Point, Text: "vwxyz"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "wxy"})))
	})

	t.Run("includedRight", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.RawCase{Point: test.Cases[common.IncludedRight].Point, Text: "wxyz"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "xyz"})))
	})

	t.Run("equal", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.MOD(test.MODSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "xyz"}), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "xyz"})))
	})
}
