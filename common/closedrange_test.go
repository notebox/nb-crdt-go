package common_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/stretchr/testify/assert"
)

func TestClosedRange(t *testing.T) {
	subject := common.ClosedRange{Lower: 5, Length: 3}

	t.Run("Upper", func(t *testing.T) {
		assert.Equal(t, uint32(7), subject.Upper())
	})

	cases := map[common.Order]common.ClosedRange{
		common.Less:            {Lower: 9, Length: 2},
		common.Prependable:     {Lower: 8, Length: 1},
		common.RightOverlap:    {Lower: 7, Length: 2},
		common.IncludingRight:  {Lower: 7, Length: 1},
		common.IncludingMiddle: {Lower: 6, Length: 1},
		common.IncludingLeft:   {Lower: 5, Length: 1},
		common.Equal:           {Lower: 5, Length: 3},
		common.IncludedLeft:    {Lower: 5, Length: 5},
		common.IncludedMiddle:  {Lower: 4, Length: 5},
		common.IncludedRight:   {Lower: 3, Length: 5},
		common.LeftOverlap:     {Lower: 4, Length: 2},
		common.Appendable:      {Lower: 4, Length: 1},
		common.Greater:         {Lower: 2, Length: 2},
	}

	t.Run("Compare", func(t *testing.T) {
		for expected, c := range cases {
			assert.Equal(t, expected, subject.Compare(c))
		}
	})

	t.Run("Intersection", func(t *testing.T) {
		// right overlap
		c := cases[common.RightOverlap]
		result, err := subject.Intersection(c)
		assert.NoError(t, err)
		assert.Equal(t, c.Lower, result.Lower)
		assert.Equal(t, subject.Lower+subject.Length-c.Lower, result.Length)

		// including
		for _, order := range []common.Order{common.IncludingRight, common.IncludingMiddle, common.IncludingLeft} {
			c = cases[order]
			result, err = subject.Intersection(c)
			assert.NoError(t, err)
			assert.Equal(t, c.Lower, result.Lower)
			assert.Equal(t, c.Length, result.Length)
		}

		// included
		for _, order := range []common.Order{common.IncludedRight, common.IncludedMiddle, common.IncludedLeft} {
			c = cases[order]
			result, err = subject.Intersection(c)
			assert.NoError(t, err)
			assert.Equal(t, subject.Lower, result.Lower)
			assert.Equal(t, subject.Length, result.Length)
		}

		// left overlap
		c = cases[common.LeftOverlap]
		result, err = subject.Intersection(c)
		assert.NoError(t, err)
		assert.Equal(t, subject.Lower, result.Lower)
		assert.Equal(t, c.Lower+c.Length-subject.Lower, result.Length)
	})

	t.Run("NoIntersection", func(t *testing.T) {
		other := common.ClosedRange{Lower: 10, Length: 3}
		_, err := subject.Intersection(other)
		assert.ErrorIs(t, err, common.NoIntersection)
	})
}
