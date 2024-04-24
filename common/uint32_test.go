package common_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/stretchr/testify/assert"
)

func TestUInt32(t *testing.T) {
	t.Run("Constants", func(t *testing.T) {
		assert.Equal(t, uint32(0), common.UInt32Min)
		assert.Equal(t, uint32(2147483647), common.UInt32Mid)
		assert.Equal(t, uint32(4294967295), common.UInt32Max)
		assert.Equal(t, uint32(4294967295), common.UInt32Max)
	})

	t.Run("CompareNumber", func(t *testing.T) {
		assert.Equal(t, common.Equal, common.CompareNumber(0, 0))
		assert.Equal(t, common.Less, common.CompareNumber(0, 1))
		assert.Equal(t, common.Greater, common.CompareNumber(1, 0))
	})
}
