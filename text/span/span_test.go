package span_test

import (
	"encoding/json"
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/point"
	"github.com/notebox/nb-crdt-go/text/span"
	"github.com/notebox/nb-crdt-go/text/test"
	"github.com/stretchr/testify/assert"
)

func TestSpan(t *testing.T) {
	subject := test.MODSpanFrom(test.Cases[common.Equal])

	t.Run("JSON", func(t *testing.T) {
		stringifiedJSON := "[[[1,2,3]],5]"
		var subject span.DELSpan
		err := json.Unmarshal([]byte(stringifiedJSON), &subject)
		assert.NoError(t, err)
		assert.Equal(t, point.Point{{Priority: 1, ReplicaID: 2, Nonce: 3}}, subject.LowerPoint())
		assert.Equal(t, uint32(5), subject.Content.Length())
		encoded, err := json.Marshal(subject)
		assert.NoError(t, err)
		assert.Equal(t, stringifiedJSON, string(encoded))
	})

	t.Run("getters", func(t *testing.T) {
		other := test.MODSpanFrom(test.Cases[common.Equal])
		assert.True(t, subject.Equals(other))
		assert.Equal(t, uint32(5), subject.ReplicaID())
		assert.Equal(t, uint32(3), subject.Length())
		assert.True(t, subject.UpperPoint().Equals(test.PointFrom([][3]uint32{{1, 1, 1}, {5, 5, 7}})))
		assert.True(t, subject.NthPoint(5).Equals(test.PointFrom([][3]uint32{{1, 1, 1}, {5, 5, 10}})))
		assert.Equal(t, subject.NonceRange(), common.ClosedRange{Lower: 5, Length: 3})
	})

	t.Run("Append", func(t *testing.T) {
		assert.Equal(t, "456789", subject.Append(test.MODSpanFrom(test.Cases[common.Prependable])).Content.Text())
	})

	t.Run("splitting", func(t *testing.T) {
		assert.Equal(t, "4", subject.LeftSplitAt(1).Content.Text())
		assert.Equal(t, "56", subject.RightSplitAt(1).Content.Text())

		l, r, err := subject.SplitAt(1)
		assert.NoError(t, err)
		assert.Equal(t, "4", l.Content.Text())
		assert.Equal(t, "56", r.Content.Text())

		l, r, err = subject.SplitWith(test.MODSpanFrom(test.Cases[common.Splitted]))
		assert.NoError(t, err)
		assert.Equal(t, "4", l.Content.Text())
		assert.Equal(t, "56", r.Content.Text())
	})

	t.Run("AppendableSegmentTo", func(t *testing.T) {
		subject := test.MODSpanFrom(test.Cases[common.Equal])

		for _, c := range []struct {
			order common.Order
			text  string
			point [][3]uint32
			err   error
		}{
			{common.Appendable, "456", [][3]uint32{{1, 1, 1}, {5, 5, 5}}, nil},
			{common.LeftOverlap, "56", [][3]uint32{{1, 1, 1}, {5, 5, 6}}, nil},
			{common.IncludingLeft, "56", [][3]uint32{{1, 1, 1}, {5, 5, 6}}, nil},
			{common.IncludingMiddle, "6", [][3]uint32{{1, 1, 1}, {5, 5, 7}}, nil},
			{common.Greater, "", nil, common.UnAppendable},
		} {
			seg, err := subject.AppendableSegmentTo(test.MODSpanFrom(test.Cases[c.order]))
			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
				continue
			}
			assert.NoError(t, err)
			assert.True(t, seg.LowerPoint().Equals(test.PointFrom(c.point)))
			assert.Equal(t, c.text, seg.Content.Text())
		}
	})

	t.Run("PrependableSegmentTo", func(t *testing.T) {
		subject := test.MODSpanFrom(test.Cases[common.Equal])

		for _, c := range []struct {
			order common.Order
			text  string
			point [][3]uint32
			err   error
		}{
			{common.Prependable, "456", [][3]uint32{{1, 1, 1}, {5, 5, 5}}, nil},
			{common.RightOverlap, "45", [][3]uint32{{1, 1, 1}, {5, 5, 5}}, nil},
			{common.IncludingRight, "45", [][3]uint32{{1, 1, 1}, {5, 5, 5}}, nil},
			{common.IncludingMiddle, "4", [][3]uint32{{1, 1, 1}, {5, 5, 5}}, nil},
			{common.Less, "", nil, common.UnPrependable},
		} {
			seg, err := subject.PrependableSegmentTo(test.MODSpanFrom(test.Cases[c.order]))
			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
				continue
			}
			assert.NoError(t, err)
			assert.True(t, seg.LowerPoint().Equals(test.PointFrom(c.point)))
			assert.Equal(t, c.text, seg.Content.Text())
		}
	})

	t.Run("Intersection", func(t *testing.T) {
		subject := test.MODSpanFrom(test.Cases[common.Equal])

		for _, c := range []struct {
			order common.Order
			text  string
			point [][3]uint32
			err   error
		}{
			{common.RightOverlap, "6", test.Cases[common.RightOverlap].Point, nil},
			{common.IncludingRight, "6", test.Cases[common.IncludingRight].Point, nil},
			{common.IncludingMiddle, "5", test.Cases[common.IncludingMiddle].Point, nil},
			{common.IncludingLeft, "4", test.Cases[common.IncludingLeft].Point, nil},
			{common.Equal, "456", test.Cases[common.Equal].Point, nil},
			{common.IncludedLeft, "456", test.Cases[common.Equal].Point, nil},
			{common.IncludedMiddle, "456", test.Cases[common.Equal].Point, nil},
			{common.IncludedRight, "456", test.Cases[common.Equal].Point, nil},
			{common.LeftOverlap, "4", test.Cases[common.Equal].Point, nil},
			{common.Splitted, "", nil, common.NoIntersection},
			{common.Less, "", nil, common.NoIntersection},
			{common.Prependable, "", nil, common.NoIntersection},
			{common.Appendable, "", nil, common.NoIntersection},
			{common.Greater, "", nil, common.NoIntersection},
			{common.Splitting, "", nil, common.NoIntersection},
		} {
			seg, err := subject.Intersection(test.MODSpanFrom(test.Cases[c.order]))
			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
				continue
			}
			assert.NoError(t, err)
			assert.True(t, seg.LowerPoint().Equals(test.PointFrom(c.point)))
			assert.Equal(t, c.text, seg.Content.Text())
		}
	})

	t.Run("Compare", func(t *testing.T) {
		subject := test.MODSpanFrom(test.Cases[common.Equal])

		var order common.Order
		var err error
		for k := range test.Cases {
			order, err = subject.Compare(test.MODSpanFrom(test.Cases[k]))
			assert.NoError(t, err)
			assert.Equal(t, k, order)
		}

		order, err = subject.Compare(test.MODSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 6, 5}}, Text: "456"}))
		assert.NoError(t, err)
		assert.Equal(t, common.Less, order)

		order, err = subject.Compare(test.MODSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 4, 5}}, Text: "456"}))
		assert.NoError(t, err)
		assert.Equal(t, common.Greater, order)

		order, err = subject.Compare(test.MODSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 5, 4}, {3, 3, 3}}, Text: "456"}))
		assert.NoError(t, err)
		assert.Equal(t, common.Greater, order)

		order, err = subject.Compare(test.MODSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 5, 8}, {3, 3, 3}}, Text: "456"}))
		assert.NoError(t, err)
		assert.Equal(t, common.Less, order)

		order, err = subject.Compare(test.MODSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}}, Text: "456"}))
		assert.NoError(t, err)
		assert.Equal(t, common.Splitting, order)

		order, err = subject.Compare(test.MODSpanFrom(test.RawCase{Point: [][3]uint32{{1, 1, 1}, {5, 5, 5}, {3, 3, 3}}, Text: "456"}))
		assert.NoError(t, err)
		assert.Equal(t, common.Splitted, order)

		subject = test.MODSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 5}, {7, 7, 7}}, Text: "456"})
		order, err = subject.Compare(test.MODSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 9}}, Text: "456"}))
		assert.NoError(t, err)
		assert.Equal(t, common.Less, order)

		order, err = subject.Compare(test.MODSpanFrom(test.RawCase{Point: [][3]uint32{{5, 5, 1}}, Text: "456"}))
		assert.NoError(t, err)
		assert.Equal(t, common.Greater, order)
	})
}
