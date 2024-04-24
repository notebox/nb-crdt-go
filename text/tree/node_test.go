package tree

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/span"
	"github.com/notebox/nb-crdt-go/text/test"
	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	t.Run("NewFromSpans, Spans", func(t *testing.T) {
		ls := test.INSSpanFrom(test.Cases[common.Greater])
		s := test.INSSpanFrom(test.Cases[common.Equal])
		rs := test.INSSpanFrom(test.Cases[common.Less])
		subject := NewFromSpans([]*span.INSSpan{ls, s, rs})

		assert.Equal(t, *New(*s, New(*ls, nil, nil), New(*rs, nil, nil)), *subject)

		spans := subject.Spans()
		assert.Equal(t, []*span.INSSpan{ls, s, rs}, spans)
	})

	t.Run("deleteSelf", func(t *testing.T) {
		l := test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 1}}, Text: "a"})
		c := test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 3}}, Text: "b"})
		r := test.INSSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 5}}, Text: "c"})
		subject := NewFromSpans([]*span.INSSpan{l, c, r})
		subject.deleteSelf()
		assert.Equal(t, []*span.INSSpan{l, r}, subject.Spans())
		subject.deleteSelf()
		assert.Equal(t, []*span.INSSpan{l}, subject.Spans())
		assert.False(t, subject.ShouldBeDeleted)
		subject.deleteSelf()
		assert.Equal(t, []*span.INSSpan{l}, subject.Spans())
		assert.True(t, subject.ShouldBeDeleted)
	})

	t.Run("Balance", func(t *testing.T) {
		// 1 / 1 + 2 = 3 / 3 + 4 = 7 / 7 + 8 = 15 /

		span := test.INSSpanFrom(test.Cases[common.Equal])
		subject := New(*span, nil, nil)
		assert.Equal(t, 1, subject.Rank)

		for i := 0; i < 2; i++ {
			subject.insertSuccessor(*span)
		}
		subject = subject.Balance()
		assert.Equal(t, 2, subject.Rank)

		for i := 0; i < 6; i++ {
			subject.insertSuccessor(*span)
		}
		subject = subject.Balance()
		assert.Equal(t, 3+1, subject.Rank)

		for i := 0; i < 14; i++ {
			subject.insertSuccessor(*span)
		}
		subject = subject.Balance()
		assert.Equal(t, 4+1, subject.Rank)
	})

	t.Run("ShouldBeDeleted", func(t *testing.T) {
		subject := New(*test.INSSpanFrom(test.Cases[common.Equal]), nil, nil)
		assert.NotNil(t, subject.Balance())
		subject.ShouldBeDeleted = true
		assert.Nil(t, subject.Balance())
	})

	t.Run("rank, length, insert, delete", func(t *testing.T) {
		cases := make(map[string]*span.INSSpan)
		for _, k := range []string{"min", "left", "pred", "mid", "succ", "right", "max"} {
			cases[k] = test.INSSpanFrom(test.RawCase{Point: test.Cases[common.Equal].Point, Text: k})
		}

		subject := New(*cases["mid"], nil, nil)
		assert.Equal(t, 1, subject.Rank)
		assert.Equal(t, uint32(3), subject.Length)
		assert.Nil(t, subject.predecessorSpan())
		assert.Nil(t, subject.successorSpan())

		key := "min"
		subject.insertPredecessor(*cases[key])
		assert.Equal(t, 2, subject.Rank)
		assert.Equal(t, uint32(6), subject.Length)
		assert.Equal(t, key, subject.predecessorSpan().Content.Text())
		assert.Nil(t, subject.successorSpan())

		key = "max"
		subject.insertSuccessor(*cases[key])
		assert.Equal(t, 2, subject.Rank)
		assert.Equal(t, uint32(9), subject.Length)
		assert.Equal(t, key, subject.successorSpan().Content.Text())

		key = "left"
		subject.insertPredecessor(*cases[key])
		assert.Equal(t, 3, subject.Rank)
		assert.Equal(t, uint32(13), subject.Length)
		assert.Equal(t, key, subject.predecessorSpan().Content.Text())

		key = "right"
		subject.insertSuccessor(*cases[key])
		assert.Equal(t, 3, subject.Rank)
		assert.Equal(t, uint32(18), subject.Length)
		assert.Equal(t, key, subject.successorSpan().Content.Text())

		key = "pred"
		subject.insertPredecessor(*cases[key])
		assert.Equal(t, 3, subject.Rank)
		assert.Equal(t, uint32(22), subject.Length)
		assert.Equal(t, key, subject.predecessorSpan().Content.Text())

		key = "succ"
		subject.insertSuccessor(*cases[key])
		assert.Equal(t, 3, subject.Rank)
		assert.Equal(t, uint32(26), subject.Length)
		assert.Equal(t, key, subject.successorSpan().Content.Text())

		subject.deleteSuccessor()
		assert.Equal(t, 3, subject.Rank)
		assert.Equal(t, uint32(22), subject.Length)
		assert.Equal(t, "right", subject.successorSpan().Content.Text())

		subject.deletePredecessor()
		assert.Equal(t, 3, subject.Rank)
		assert.Equal(t, uint32(18), subject.Length)
		assert.Equal(t, "left", subject.predecessorSpan().Content.Text())

		subject.deleteSuccessor()
		assert.Equal(t, 3, subject.Rank)
		assert.Equal(t, uint32(13), subject.Length)
		assert.Equal(t, "max", subject.successorSpan().Content.Text())

		subject.deletePredecessor()
		assert.Equal(t, 2, subject.Rank)
		assert.Equal(t, uint32(9), subject.Length)
		assert.Equal(t, "min", subject.predecessorSpan().Content.Text())

		subject.deleteSuccessor()
		assert.Equal(t, 2, subject.Rank)
		assert.Equal(t, uint32(6), subject.Length)
		assert.Nil(t, subject.successorSpan())

		subject.deletePredecessor()
		assert.Equal(t, 1, subject.Rank)
		assert.Equal(t, uint32(3), subject.Length)
		assert.Nil(t, subject.predecessorSpan())
	})
}
