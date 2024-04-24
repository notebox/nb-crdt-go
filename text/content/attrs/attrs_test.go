package attrs_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/content/attrs"
	"github.com/stretchr/testify/assert"
)

func TestAttributes(t *testing.T) {
	t.Run("Concat", func(t *testing.T) {
		var subject attrs.Attrs
		leaves := []attrs.Attr{
			{
				Length: 5,
				Props:  attrs.TextProps{"B": true},
				Stamp:  &common.Stamp{ReplicaID: 3, Timestamp: 9},
			},
			{
				Length: 10,
				Props:  attrs.TextProps{"S": true},
				Stamp:  &common.Stamp{ReplicaID: 3, Timestamp: 9},
			},
		}
		subject = subject.Concat(leaves)
		assert.Equal(t, attrs.Attrs(leaves), subject)

		concatenated := subject.Concat(leaves)
		assert.Equal(t, concatenated, attrs.Attrs(append(leaves, leaves...)))

		optimized := subject.Concat([]attrs.Attr{{Length: 2, Props: leaves[1].Props, Stamp: leaves[1].Stamp}})
		assert.Equal(t, len(optimized), 2)
		assert.Equal(t, optimized[1].Length, uint32(12))
		assert.True(t, optimized[1].EqualsExceptForLength(&leaves[1]))
	})

	t.Run("Slice", func(t *testing.T) {
		leaves := []attrs.Attr{
			{
				Length: 5,
				Props:  attrs.TextProps{"B": true},
				Stamp:  &common.Stamp{ReplicaID: 3, Timestamp: 9},
			},
			{
				Length: 10,
				Props:  attrs.TextProps{"S": true},
				Stamp:  &common.Stamp{ReplicaID: 3, Timestamp: 9},
			},
		}
		subject := attrs.Attrs(leaves)

		assert.Equal(t, 0, len(subject.Slice(3, 3)))
		assert.Equal(t, attrs.Attrs{leaves[0]}, subject.Slice(0, 5))
		assert.Equal(t, attrs.Attrs{leaves[1]}, subject.Slice(5, 15))
		assert.Equal(t, attrs.Attrs{
			{Length: 1, Props: leaves[0].Props, Stamp: leaves[0].Stamp},
			{Length: 1, Props: leaves[1].Props, Stamp: leaves[1].Stamp},
		}, subject.Slice(4, 6))
	})

	t.Run("Merge - one to nil", func(t *testing.T) {
		var subject attrs.Attrs
		subject.Merge(attrs.Attrs{{
			Length: 10,
			Props:  attrs.TextProps{"B": true},
			Stamp:  &common.Stamp{ReplicaID: 2, Timestamp: 2},
		}})
		assert.Nil(t, subject)
	})

	t.Run("Merge - one to many", func(t *testing.T) {
		subject := attrs.Attrs{
			{
				Length: 2,
				Props:  attrs.TextProps{"I": true},
				Stamp:  nil,
			},
			{
				Length: 2,
				Props:  nil,
				Stamp:  &common.Stamp{ReplicaID: 1, Timestamp: 1},
			},
			{
				Length: 2,
				Props:  attrs.TextProps{"B": true},
				Stamp:  &common.Stamp{ReplicaID: 1, Timestamp: 1},
			},
			{
				Length: 2,
				Props:  attrs.TextProps{"S": true},
				Stamp:  &common.Stamp{ReplicaID: 1, Timestamp: 1},
			},
			{
				Length: 2,
				Props:  attrs.TextProps{"CODE": true},
				Stamp:  &common.Stamp{ReplicaID: 3, Timestamp: 3},
			},
		}
		subject.Merge(attrs.Attrs{{
			Length: 10,
			Props:  attrs.TextProps{"B": true},
			Stamp:  &common.Stamp{ReplicaID: 2, Timestamp: 2},
		}})
		expected := attrs.Attrs{
			{
				Length: 8,
				Props:  attrs.TextProps{"B": true},
				Stamp:  &common.Stamp{ReplicaID: 2, Timestamp: 2},
			},
			{
				Length: 2,
				Props:  attrs.TextProps{"CODE": true},
				Stamp:  &common.Stamp{ReplicaID: 3, Timestamp: 3},
			},
		}
		assert.Equal(t, expected, subject)
	})

	t.Run("Merge - many to one", func(t *testing.T) {
		subject := attrs.Attrs{{
			Length: 10,
			Props:  attrs.TextProps{"B": true},
			Stamp:  nil,
		}}
		subject.Merge(attrs.Attrs{
			{
				Length: 2,
				Props:  attrs.TextProps{"I": true},
				Stamp:  &common.Stamp{ReplicaID: 1, Timestamp: 1},
			},
			{
				Length: 2,
				Props:  nil,
				Stamp:  &common.Stamp{ReplicaID: 1, Timestamp: 1},
			},
			{
				Length: 2,
				Props:  attrs.TextProps{"B": true},
				Stamp:  &common.Stamp{ReplicaID: 1, Timestamp: 1},
			},
			{
				Length: 2,
				Props:  attrs.TextProps{"S": true},
				Stamp:  &common.Stamp{ReplicaID: 1, Timestamp: 1},
			},
			{
				Length: 2,
				Props:  attrs.TextProps{"CODE": true},
				Stamp:  &common.Stamp{ReplicaID: 3, Timestamp: 3},
			},
		})
		expected := attrs.Attrs{
			{
				Length: 2,
				Props:  attrs.TextProps{"I": true},
				Stamp:  &common.Stamp{ReplicaID: 1, Timestamp: 1},
			},
			{
				Length: 2,
				Props:  nil,
				Stamp:  &common.Stamp{ReplicaID: 1, Timestamp: 1},
			},
			{
				Length: 2,
				Props:  attrs.TextProps{"B": true},
				Stamp:  &common.Stamp{ReplicaID: 1, Timestamp: 1},
			},
			{
				Length: 2,
				Props:  attrs.TextProps{"S": true},
				Stamp:  &common.Stamp{ReplicaID: 1, Timestamp: 1},
			},
			{
				Length: 2,
				Props:  attrs.TextProps{"CODE": true},
				Stamp:  &common.Stamp{ReplicaID: 3, Timestamp: 3},
			},
		}
		assert.Equal(t, expected, subject)
	})
}
