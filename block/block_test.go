package block_test

import (
	"encoding/json"
	"testing"

	"github.com/notebox/nb-crdt-go/block"
	"github.com/notebox/nb-crdt-go/common"
	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		cases := []string{
			`["86905c41-46bc-402e-8e72-f298be4e72e9",{},[[0,0,1]],{"TYPE":[null,"NOTE"]},false,[[[[2147483648,9777,1]],[[[1]],"1"]]],"8aa83876-57ca-422b-a290-b070dd07d2f7"]`,
			`["86905c41-46bc-402e-8e72-f298be4e72e9",{},[[0,0,1]],{"TYPE":[null,"NOTE"]},false,null,"8aa83876-57ca-422b-a290-b070dd07d2f7"]`,
			`["86905c41-46bc-402e-8e72-f298be4e72e9",{},[[0,0,1]],{"TYPE":[null,"NOTE"]},false]`,
		}
		for _, stringifiedJSON := range cases {
			var subject block.Block
			err := json.Unmarshal([]byte(stringifiedJSON), &subject)
			assert.NoError(t, err)
			encoded, err := json.Marshal(subject)
			assert.NoError(t, err)
			assert.Equal(t, stringifiedJSON, string(encoded))
		}
	})

	t.Run("Apply", func(t *testing.T) {
		stringifiedBlock := `["86905C41-46BC-402E-8E72-F298BE4E72E9",{},[[0,0,1]],{"TYPE":[null,"NOTE"]},false,[[[[2147483648,9777,1]],[[[1]],"1"]]],"8AA83876-57CA-422B-A290-B070DD07D2F7"]`
		var subject block.Block
		err := json.Unmarshal([]byte(stringifiedBlock), &subject)
		assert.NoError(t, err)

		stringifiedOPs := `{"bDEL":false,"bINS":["76905C41-46BC-402E-8E72-F298BE4E72E9",{},[[2147483647,0,3]],{},false,null,null],"bMOV":["66905C41-46BC-402E-8E72-F298BE4E72E9",[[1,2,3]]],"bSET":{"TYPE":"LINE","SRC":"null"},"tDEL":[[[[1,2,3]],9]],"tFMT":[[[[1,2,3]],[1,[[1,{"B":true},[8,7]]]]]],"tINS":[[[[1,2,3]],[[[3,null,null]],"abc"]]],"tMOD":[[[[1,2,3]],"abc"]]}`
		var ops block.Operations
		err = json.Unmarshal([]byte(stringifiedOPs), &ops)
		assert.NoError(t, err)
		ctrb := block.Contribution{
			BlockID:    subject.BlockID,
			Nonce:      block.ReplicaNonce{87, 59},
			Stamp:      common.Stamp{ReplicaID: 11, Timestamp: 12},
			Operations: ops,
		}
		err = subject.Apply(ctrb)
		assert.NoError(t, err)
		assert.Equal(t, block.ReplicaNonce{87, 59}, subject.Version[11])
	})

	t.Run("BlockPoint JSON", func(t *testing.T) {
		stringifiedJSON := `["86905c41-46bc-402e-8e72-f298be4e72e9",[[1,2,3]]]`
		var subject block.BlockPoint
		err := json.Unmarshal([]byte(stringifiedJSON), &subject)
		assert.NoError(t, err)
		encoded, err := json.Marshal(subject)
		assert.NoError(t, err)
		assert.Equal(t, stringifiedJSON, string(encoded))
	})
}
