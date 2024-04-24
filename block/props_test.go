package block_test

import (
	"encoding/json"
	"testing"

	"github.com/notebox/nb-crdt-go/block"
	"github.com/notebox/nb-crdt-go/common"
	"github.com/stretchr/testify/assert"
)

func TestBlockProps(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		stringifiedJSON := `{"DB_RECORD":{"0-0":{"VALUE":[null,"go"]}},"MOV":[null],"SRC":[[5,9],"notebox.cloud"],"TYPE":[null,"NOTE"]}`
		var props block.Props
		err := json.Unmarshal([]byte(stringifiedJSON), &props)
		assert.NoError(t, err)
		assert.Equal(t, block.Props{
			"TYPE": block.Prop{Nested: nil, Stamp: nil, Value: "NOTE"},
			"MOV":  block.Prop{Nested: nil, Stamp: nil, Value: nil},
			"SRC":  block.Prop{Nested: nil, Stamp: &common.Stamp{ReplicaID: 5, Timestamp: 9}, Value: "notebox.cloud"},
			"DB_RECORD": block.Prop{
				Nested: map[string]block.Prop{
					"0-0": {
						Nested: map[string]block.Prop{
							"VALUE": {Nested: nil, Stamp: nil, Value: "go"},
						},
					},
				},
				Stamp: nil,
				Value: nil,
			},
		}, props)
		encoded, err := json.Marshal(props)
		assert.NoError(t, err)
		assert.Equal(t, stringifiedJSON, string(encoded))
	})

	t.Run("IsLeaf", func(t *testing.T) {
		var subject block.Prop
		subject = block.Prop{}
		assert.False(t, subject.IsNotLeaf())
		subject = block.Prop{Nested: map[string]block.Prop{}}
		assert.True(t, subject.IsNotLeaf())
		subject = block.Prop{Nested: map[string]block.Prop{"LEAF": {}}}
		assert.True(t, subject.IsNotLeaf())
		subject = block.Prop{Nested: map[string]block.Prop{"LEAF": {Value: true}}}
		assert.False(t, subject.IsNotLeaf())
	})

	t.Run("UpdateProps", func(t *testing.T) {
		cases := []struct {
			stamp    common.Stamp
			delta    string
			expected string
		}{
			{common.Stamp{ReplicaID: 1, Timestamp: 2}, `{"X":{"Y":{"Z":"A"}}}`, `{"TYPE":[null,"T"],"X":{"Y":{"Z":[[1,2],"A"]}}}`},
			{common.Stamp{ReplicaID: 2, Timestamp: 3}, `{"X":{"Y":{"Z":"B"}}}`, `{"TYPE":[null,"T"],"X":{"Y":{"Z":[[2,3],"B"]}}}`},
			{common.Stamp{ReplicaID: 3, Timestamp: 4}, `{"X":{"Y":{"Z":null}}}`, `{"TYPE":[null,"T"],"X":{"Y":{"Z":[[3,4]]}}}`},
		}
		for _, c := range cases {
			subject := block.Props{
				"TYPE": block.Prop{Nested: nil, Stamp: nil, Value: "T"},
			}
			var delta block.PropsDelta
			err := json.Unmarshal([]byte(c.delta), &delta)
			assert.NoError(t, err)
			updated := block.UpdateProps(subject, delta, c.stamp)
			encoded, err := json.Marshal(updated)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, string(encoded))
		}
	})
}

func TestBlockPropsDelta(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		stringifiedJSON := `{"X":{"Y":{"Z":"A"}}}`
		var delta block.PropsDelta
		err := json.Unmarshal([]byte(stringifiedJSON), &delta)
		assert.NoError(t, err)
		assert.Equal(t, block.PropsDelta{
			"X": block.PropDelta{
				Nested: map[string]block.PropDelta{
					"Y": {
						Nested: map[string]block.PropDelta{
							"Z": {Nested: nil, Value: "A"},
						},
					},
				},
			},
		}, delta)
		encoded, err := json.Marshal(delta)
		assert.NoError(t, err)
		assert.Equal(t, stringifiedJSON, string(encoded))
	})
}
