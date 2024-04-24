package common_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/stretchr/testify/assert"
)

func TestUTF16(t *testing.T) {
	t.Run("UTF16Slice", func(t *testing.T) {
		assert.Equal(t, "😀bb🤖cc👍🏿dd강e", common.UTF16Slice("aa😀bb🤖cc👍🏿dd강ee", 2, 18))
	})

	t.Run("UTF16Length", func(t *testing.T) {
		assert.Equal(t, uint32(19), common.UTF16Length("aa😀bb🤖cc👍🏿dd강ee"))
	})
}
