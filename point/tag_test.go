package point_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common" // Replace with the correct package path
	"github.com/notebox/nb-crdt-go/point"
	"github.com/stretchr/testify/assert"
)

func TestPointTag(t *testing.T) {
	var subject point.PointTag

	t.Run("MinTag", func(t *testing.T) {
		subject = point.MinTag
		assert.Equal(t, common.Priority(common.UInt32Min), subject.Priority)
		assert.Equal(t, common.ReplicaID(0), subject.ReplicaID)
		assert.Equal(t, uint32(1), subject.Nonce)
	})

	t.Run("MidTag", func(t *testing.T) {
		subject = point.MidTag
		assert.Equal(t, common.Priority(common.UInt32Mid), subject.Priority)
		assert.Equal(t, common.ReplicaID(0), subject.ReplicaID)
		assert.Equal(t, uint32(3), subject.Nonce)
	})

	t.Run("MaxTag", func(t *testing.T) {
		subject = point.MaxTag
		assert.Equal(t, common.Priority(common.UInt32Max), subject.Priority)
		assert.Equal(t, common.ReplicaID(0), subject.ReplicaID)
		assert.Equal(t, uint32(2), subject.Nonce)
	})

	subject = point.PointTag{
		Priority:  1,
		ReplicaID: 2,
		Nonce:     3,
	}

	t.Run("WithNonce", func(t *testing.T) {
		target := subject.WithNonce(4)
		assert.Equal(t, subject.Priority, target.Priority)
		assert.Equal(t, subject.ReplicaID, target.ReplicaID)
		assert.Equal(t, uint32(4), target.Nonce)
		assert.NotEqual(t, subject.Nonce, target.Nonce)
	})

	t.Run("CompareBase", func(t *testing.T) {
		cases := []struct {
			expected common.Order
			tag      point.PointTag
		}{
			{common.Equal, point.PointTag{Priority: subject.Priority, ReplicaID: subject.ReplicaID, Nonce: subject.Nonce}},
			{common.Equal, point.PointTag{Priority: subject.Priority, ReplicaID: subject.ReplicaID, Nonce: subject.Nonce + 1}},
			{common.Greater, point.PointTag{Priority: subject.Priority - 1, ReplicaID: subject.ReplicaID, Nonce: subject.Nonce + 1}},
			{common.Greater, point.PointTag{Priority: subject.Priority, ReplicaID: subject.ReplicaID - 1, Nonce: subject.Nonce + 1}},
			{common.Less, point.PointTag{Priority: subject.Priority + 1, ReplicaID: subject.ReplicaID, Nonce: subject.Nonce - 1}},
			{common.Less, point.PointTag{Priority: subject.Priority, ReplicaID: subject.ReplicaID + 1, Nonce: subject.Nonce - 1}},
		}
		for _, c := range cases {
			assert.Equal(t, c.expected, subject.CompareBase(c.tag))
		}
	})

	t.Run("Compare", func(t *testing.T) {
		cases := []struct {
			expected common.Order
			tag      point.PointTag
		}{
			{common.Equal, point.PointTag{Priority: subject.Priority, ReplicaID: subject.ReplicaID, Nonce: subject.Nonce}},
			{common.Less, point.PointTag{Priority: subject.Priority, ReplicaID: subject.ReplicaID, Nonce: subject.Nonce + 1}},
			{common.Greater, point.PointTag{Priority: subject.Priority, ReplicaID: subject.ReplicaID, Nonce: subject.Nonce - 1}},
			{common.Greater, point.PointTag{Priority: subject.Priority - 1, ReplicaID: subject.ReplicaID, Nonce: subject.Nonce + 1}},
			{common.Greater, point.PointTag{Priority: subject.Priority, ReplicaID: subject.ReplicaID - 1, Nonce: subject.Nonce + 1}},
			{common.Less, point.PointTag{Priority: subject.Priority + 1, ReplicaID: subject.ReplicaID, Nonce: subject.Nonce - 1}},
			{common.Less, point.PointTag{Priority: subject.Priority, ReplicaID: subject.ReplicaID + 1, Nonce: subject.Nonce - 1}},
		}
		for _, c := range cases {
			assert.Equal(t, c.expected, subject.Compare(c.tag))
		}
	})
}
