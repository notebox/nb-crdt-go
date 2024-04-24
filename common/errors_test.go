package common_test

import (
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/stretchr/testify/assert"
)

func TestFatalErrors(t *testing.T) {
	tests := []struct {
		name string
		err  common.FatalError
	}{
		{
			name: "NoIntersection",
			err:  common.NoIntersection,
		},
		{
			name: "ExistingSpanOverwrite",
			err:  common.ExistingSpanOverwrite,
		},
		{
			name: "UnAppendable",
			err:  common.UnAppendable,
		},
		{
			name: "UnPrependable",
			err:  common.UnPrependable,
		},
		{
			name: "InvalidDistanceBetweenNoRelation",
			err:  common.InvalidDistanceBetweenNoRelation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.name, tt.err.Error())
		})
	}
}
