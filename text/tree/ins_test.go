package tree_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/test"
	"github.com/notebox/nb-crdt-go/text/tree"
	"github.com/stretchr/testify/assert"
)

func TestNodeINS(t *testing.T) {
	t.Run("less", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			tree.New(*test.INSSpanFrom(test.Cases[common.Less]), nil, nil),
		)
		span := test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 6, 9}}, Text: "x"})
		err := subject.INS(span, 0)
		assert.NoError(t, err)
		assert.Nil(t, subject.Left)
		assert.True(t, subject.Right.Right.Span.Equals(span))
		assert.True(t, subject.Right.Span.Equals(test.INSSpanFrom(test.Cases[common.Less])))
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("prependable", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.INS(test.INSSpanFrom(test.Cases[common.Prependable]), 0)
		assert.NoError(t, err)
		assert.Nil(t, subject.Left)
		assert.Nil(t, subject.Right)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "456789"})))
	})

	t.Run("greater", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			tree.New(*test.INSSpanFrom(test.Cases[common.Greater]), nil, nil),
			nil,
		)
		span := test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 4, 1}}, Text: "x"})
		err := subject.INS(span, 0)
		assert.NoError(t, err)
		assert.True(t, subject.Left.Left.Span.Equals(span))
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
		err := subject.INS(test.INSSpanFrom(test.Cases[common.Appendable]), 0)
		assert.NoError(t, err)
		assert.Nil(t, subject.Left)
		assert.Nil(t, subject.Right)
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Appendable].Point, Text: "123456"})))
	})

	t.Run("splitted", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.INS(test.INSSpanFrom(test.Cases[common.Splitted]), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Left.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: "4"})))
		assert.True(t, subject.Right.Span.Equals(test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 5, 6}}, Text: "56"})))
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Splitted])))
	})

	t.Run("splitting", func(t *testing.T) {
		subject := tree.New(
			*test.INSSpanFrom(test.Cases[common.Equal]),
			nil,
			nil,
		)
		err := subject.INS(test.INSSpanFrom(test.Cases[common.Splitting]), 0)
		assert.NoError(t, err)
		assert.True(t, subject.Left.Span.Equals(test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Splitting].Point, Text: "j"})))
		assert.True(t, subject.Right.Span.Equals(test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 2}}, Text: "kl"})))
		assert.True(t, subject.Span.Equals(test.INSSpanFrom(test.Cases[common.Equal])))
	})

	t.Run("right overlap, left overlap, including left, including middle, including right, included left, included middle, included right, equal", func(t *testing.T) {
		for _, order := range []common.Order{common.RightOverlap, common.LeftOverlap, common.IncludingLeft, common.IncludingMiddle, common.IncludingRight, common.IncludedLeft, common.IncludedMiddle, common.IncludedRight, common.Equal} {
			subject := tree.New(
				*test.INSSpanFrom(test.Cases[common.Equal]),
				nil,
				nil,
			)
			err := subject.INS(test.INSSpanFrom(test.Cases[order]), 0)
			assert.ErrorIs(t, err, common.ExistingSpanOverwrite)
		}
	})
}
