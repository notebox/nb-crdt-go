package attrs_test

import (
	"encoding/json"
	"testing"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/content/attrs"
	"github.com/stretchr/testify/assert"
)

func TestAttr(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		cases := []struct {
			expected        attrs.Attr
			stringifiedJSON string
		}{
			{
				attrs.Attr{Length: 5, Props: map[string]any{"B": true}, Stamp: &common.Stamp{ReplicaID: 3, Timestamp: 9}},
				`[5,{"B":true},[3,9]]`,
			},
			{
				attrs.Attr{Length: 5, Stamp: &common.Stamp{ReplicaID: 3, Timestamp: 9}},
				`[5,null,[3,9]]`,
			},
			{
				attrs.Attr{Length: 5, Props: map[string]any{"B": true}},
				`[5,{"B":true}]`,
			},
			{
				attrs.Attr{Length: 5},
				`[5]`,
			},
		}

		for _, c := range cases {
			var leaf attrs.Attr
			err := json.Unmarshal([]byte(c.stringifiedJSON), &leaf)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, leaf)
			encoded, err := json.Marshal(leaf)
			assert.NoError(t, err)
			assert.Equal(t, c.stringifiedJSON, string(encoded))
		}
	})

	t.Run("EqualsExceptForLength", func(t *testing.T) {
		subject := attrs.Attr{Length: 5, Props: map[string]any{"B": true}, Stamp: &common.Stamp{ReplicaID: 3, Timestamp: 9}}

		assert.True(t, subject.EqualsExceptForLength(&attrs.Attr{Length: 5, Props: map[string]any{"B": true}, Stamp: &common.Stamp{ReplicaID: 3, Timestamp: 9}}))
		assert.True(t, subject.EqualsExceptForLength(&attrs.Attr{Length: 3, Props: map[string]any{"B": true}, Stamp: &common.Stamp{ReplicaID: 3, Timestamp: 9}}))
		assert.False(t, subject.EqualsExceptForLength(&attrs.Attr{Length: 5, Props: map[string]any{"B": true}, Stamp: &common.Stamp{ReplicaID: 4, Timestamp: 9}}))
		assert.False(t, subject.EqualsExceptForLength(&attrs.Attr{Length: 5, Props: map[string]any{"B": false}, Stamp: &common.Stamp{ReplicaID: 3, Timestamp: 9}}))
		assert.False(t, subject.EqualsExceptForLength(&attrs.Attr{Length: 5, Stamp: &common.Stamp{ReplicaID: 3, Timestamp: 9}}))
	})

	t.Run("Apply", func(t *testing.T) {
		subject := attrs.Attr{
			Length: 5,
			Props:  map[string]any{"B": true, "S": true, "COLOR": "red"},
			Stamp:  &common.Stamp{ReplicaID: 3, Timestamp: 9},
		}
		stamp := common.Stamp{ReplicaID: 0, Timestamp: 0}

		subject.Apply(map[string]any{"COLOR": "blue", "S": nil}, &stamp)
		assert.Equal(t, map[string]any{"COLOR": "blue", "B": true}, subject.Props)
		assert.Equal(t, stamp, *subject.Stamp)

		subject.Apply(map[string]any{"COLOR": nil, "B": nil}, &stamp)
		assert.Equal(t, 0, len(subject.Props))
		assert.Equal(t, stamp, *subject.Stamp)
	})
}
