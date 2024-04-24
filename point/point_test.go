package point_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/point"
	"github.com/stretchr/testify/assert"
)

func TestPoint(t *testing.T) {
	t.Run("Clone", func(t *testing.T) {
		subject := &point.Point{{}}
		cloned := subject.Clone()
		cloned[0].Nonce = 1
		assert.NotEqual(t, (*subject)[0].Nonce, cloned[0].Nonce)
	})

	t.Run("WithNonce", func(t *testing.T) {
		subject := point.Point{point.MidTag, point.MinTag}
		withNonce := subject.WithNonce(point.MinTag.Nonce + 1)
		assert.Equal(t, subject[0], withNonce[0])
		assert.Equal(t, subject[1].Priority, withNonce[1].Priority)
		assert.Equal(t, subject[1].ReplicaID, withNonce[1].ReplicaID)
		assert.Equal(t, subject[1].Nonce, withNonce[1].Nonce-1)
	})

	t.Run("Offset", func(t *testing.T) {
		subject := point.Point{point.MidTag, point.MinTag}
		withNonce := subject.Offset(59)
		assert.Equal(t, subject[0], withNonce[0])
		assert.Equal(t, subject[1].Priority, withNonce[1].Priority)
		assert.Equal(t, subject[1].ReplicaID, withNonce[1].ReplicaID)
		assert.Equal(t, subject[1].Nonce, withNonce[1].Nonce-59)
	})

	t.Run("CompareBase", func(t *testing.T) {
		testCases := []struct {
			name     string
			point    point.Point
			other    point.Point
			expected common.Order
		}{
			{
				name:     "Equal points",
				point:    point.Point{point.MidTag},
				other:    point.Point{point.MidTag},
				expected: common.Equal,
			},
			{
				name:     "Different points",
				point:    point.Point{point.MidTag},
				other:    point.Point{point.MaxTag},
				expected: common.Less,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := tc.point.CompareBase(tc.other)
				assert.Equal(t, tc.expected, result)
			})
		}
	})

	t.Run("CompareBase", func(t *testing.T) {
		t.Run("Equal points", func(t *testing.T) {
			a := point.Point{{5, 5, 5}}
			b := point.Point{{6, 5, 5}}
			result := a.CompareBase(b)
			assert.Equal(t, common.Equal, result)

			a = point.Point{{5, 5, 5}}
			b = point.Point{{5, 6, 5}, {1, 5, 5}}
			result = a.CompareBase(b)
			assert.Equal(t, common.Equal, result)

			a = point.Point{{5, 5, 5}, {3, 3, 3}}
			b = point.Point{{5, 5, 5}, {3, 3, 6}}
			result = a.CompareBase(b)
			assert.Equal(t, common.Equal, result)
		})

		t.Run("Less and Greater", func(t *testing.T) {
			a := point.Point{{5, 5, 5}}
			b := point.Point{{6, 5, 4}}
			assert.Equal(t, common.Less, a.CompareBase(b))
			assert.Equal(t, common.Greater, b.CompareBase(a))

			a = point.Point{{5, 5, 5}}
			b = point.Point{{5, 6, 5}, {1, 2, 3}}
			assert.Equal(t, common.Less, a.CompareBase(b))
			assert.Equal(t, common.Greater, b.CompareBase(a))

			a = point.Point{{5, 5, 5}, {3, 3, 3}}
			b = point.Point{{5, 5, 5}, {4, 4, 4}}
			assert.Equal(t, common.Less, a.CompareBase(b))
			assert.Equal(t, common.Greater, b.CompareBase(a))
		})

		t.Run("Tagging and Tagged", func(t *testing.T) {
			a := point.Point{{5, 5, 5}}
			b := point.Point{{5, 5, 5}, {1, 2, 3}}
			assert.Equal(t, common.Tagged, a.CompareBase(b))
			assert.Equal(t, common.Tagging, b.CompareBase(a))

			a = point.Point{{5, 5, 5}}
			b = point.Point{{5, 5, 1}, {1, 2, 3}}
			assert.Equal(t, common.Tagged, a.CompareBase(b))
			assert.Equal(t, common.Tagging, b.CompareBase(a))

			a = point.Point{{5, 5, 5}}
			b = point.Point{{5, 5, 9}, {1, 2, 3}}
			assert.Equal(t, common.Tagged, a.CompareBase(b))
			assert.Equal(t, common.Tagging, b.CompareBase(a))

			a = point.Point{{5, 5, 5}, {3, 3, 3}}
			b = point.Point{{5, 5, 5}, {3, 3, 3}, {1, 2, 3}}
			assert.Equal(t, common.Tagged, a.CompareBase(b))
			assert.Equal(t, common.Tagging, b.CompareBase(a))

			a = point.Point{{5, 5, 5}, {3, 3, 3}}
			b = point.Point{{5, 5, 5}, {3, 3, 1}, {1, 2, 3}}
			assert.Equal(t, common.Tagged, a.CompareBase(b))
			assert.Equal(t, common.Tagging, b.CompareBase(a))

			a = point.Point{{5, 5, 5}, {3, 3, 3}}
			b = point.Point{{5, 5, 5}, {3, 3, 9}, {1, 2, 3}}
			assert.Equal(t, common.Tagged, a.CompareBase(b))
			assert.Equal(t, common.Tagging, b.CompareBase(a))
		})
	})

	t.Run("Compare", func(t *testing.T) {
		subject := point.Point{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

		t.Run("Equal", func(t *testing.T) {
			a := point.Point{{5, 5, 5}}
			b := point.Point{{6, 5, 5}}
			result := a.Compare(b)
			assert.Equal(t, common.Equal, result)

			a = point.Point{{5, 5, 5}}
			b = point.Point{{5, 6, 5}, {1, 5, 5}}
			result = a.Compare(b)
			assert.Equal(t, common.Equal, result)

			a = point.Point{{5, 5, 5}, {3, 3, 3}}
			b = point.Point{{5, 5, 5}, {4, 3, 3}}
			result = a.Compare(b)
			assert.Equal(t, common.Equal, result)
		})

		t.Run("Less and Greater", func(t *testing.T) {
			a := point.Point{{5, 5, 5}}
			b := point.Point{{5, 5, 5}, {3, 3, 3}}
			assert.Equal(t, common.Less, a.Compare(b))
			assert.Equal(t, common.Greater, b.Compare(a))
		})

		t.Run("CompareFromTheFirst", func(t *testing.T) {
			other := point.Point{{1, 2, 2}, {5, 6, 7}}
			assert.Equal(t, common.Greater, subject.Compare(other))
			assert.Equal(t, common.Less, other.Compare(subject))
		})

		t.Run("TaggingReturnsLess", func(t *testing.T) {
			other := point.Point{{1, 2, 3}, {4, 5, 6}}
			assert.Equal(t, common.Greater, subject.Compare(other))
			assert.Equal(t, common.Less, other.Compare(subject))
		})
	})

	t.Run("DistanceFrom", func(t *testing.T) {
		subject := point.Point{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

		t.Run("BaseTaggingTagged", func(t *testing.T) {
			other := point.Point{{1, 2, 3}, {4, 5, 3}}
			dist, order, err := subject.DistanceFrom(other)

			assert.Equal(t, uint32(3), dist)
			assert.Equal(t, common.Greater, order)
			assert.NoError(t, err)
		})

		t.Run("BaseEqual", func(t *testing.T) {
			other := point.Point{{1, 2, 3}, {4, 5, 6}, {7, 8, 120}}
			dist, order, err := subject.DistanceFrom(other)

			assert.Equal(t, uint32(111), dist)
			assert.Equal(t, common.Less, order)
			assert.NoError(t, err)
		})

		t.Run("ErrorIfPointBasesAreNotEqual", func(t *testing.T) {
			other := point.Point{{1, 2, 3}, {4, 5, 6}, {7, 9, 9}}
			_, _, err := subject.DistanceFrom(other)
			assert.ErrorIs(t, common.InvalidDistanceBetweenNoRelation, err)
			_, _, err = other.DistanceFrom(subject)
			assert.ErrorIs(t, common.InvalidDistanceBetweenNoRelation, err)
		})
	})
}
