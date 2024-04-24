package tree_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/test"
	"github.com/notebox/nb-crdt-go/text/tree"
	"github.com/stretchr/testify/assert"
)

func TestNodeDEL(t *testing.T) {
	t.Run("less", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			tree.New(*test.INSSpanFrom(test.Cases[common.Less]), nil, nil),
		)
		err := subject.DEL(test.DELSpanFrom(test.Cases[common.Less]), 0)
		assert.NoError(t, err)
		assert.Nil(t, subject.Left)
		assert.Nil(t, subject.Right)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("prependable", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			tree.New(*test.INSSpanFrom(test.Cases[common.Prependable]), nil, nil),
		)
		err := subject.DEL(test.DELSpanFrom(test.Cases[common.Prependable]), 0)
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
		err := subject.DEL(test.DELSpanFrom(test.Cases[common.Greater]), 0)
		assert.NoError(t, err)
		assert.Nil(t, subject.Left)
		assert.Nil(t, subject.Right)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("appendable", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			tree.New(*test.INSSpanFrom(test.Cases[common.Appendable]), nil, nil),
			nil,
		)
		err := subject.DEL(test.DELSpanFrom(test.Cases[common.Appendable]), 0)
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
		err := subject.DEL(test.DELSpanFrom(test.Cases[common.IncludingLeft]), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 5, 6}}, Text: "56"})))
	})

	t.Run("includingRight", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.DEL(test.DELSpanFrom(test.Cases[common.IncludingRight]), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 5, 5}}, Text: "45"})))
	})

	t.Run("includingMiddle", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.DEL(test.DELSpanFrom(test.Cases[common.IncludingMiddle]), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 5, 5}}, Text: "4"})))
		assert.True(t, subject.Right.Span.Equals(test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 5, 7}}, Text: "6"})))
	})

	t.Run("rightOverlap", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.DEL(test.DELSpanFrom(test.Cases[common.RightOverlap]), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 5, 5}}, Text: "45"})))
	})

	t.Run("leftOverlap", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.DEL(test.DELSpanFrom(test.Cases[common.LeftOverlap]), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 5, 6}}, Text: "56"})))
	})

	t.Run("splitted", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.DEL(test.DELSpanFrom(test.Cases[common.Splitted]), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("splitting", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.DEL(test.DELSpanFrom(test.Cases[common.Splitting]), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("equal, included left, included right, included middle", func(t *testing.T) {
		for _, order := range []common.Order{common.Equal, common.IncludedLeft, common.IncludedRight, common.IncludedMiddle} {
			subject := tree.New(
				*test.INSSpanFrom(test.Cases[common.Equal]),
				nil,
				nil,
			)
			err := subject.DEL(test.DELSpanFrom(test.Cases[order]), 0)
			assert.NoError(t, err)
			assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
			assert.True(t, subject.ShouldBeDeleted)
		}
	})
}
