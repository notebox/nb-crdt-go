package point

import (
	"encoding/json"

	"github.com/notebox/nb-crdt-go/common"
)

type PointTag struct {
	Priority  common.Priority
	ReplicaID common.ReplicaID
	Nonce     common.Nonce
}

func (pt PointTag) MarshalJSON() ([]byte, error) {
	return json.Marshal([]uint32{pt.Priority, pt.ReplicaID, pt.Nonce})
}

func (pt *PointTag) UnmarshalJSON(data []byte) error {
	var raw []uint32
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	pt.Priority = raw[0]
	pt.ReplicaID = raw[1]
	pt.Nonce = raw[2]
	return nil
}

func (pt *PointTag) WithNonce(nonce uint32) PointTag {
	return PointTag{
		Priority:  pt.Priority,
		ReplicaID: pt.ReplicaID,
		Nonce:     nonce,
	}
}

func (pt *PointTag) CompareBase(other PointTag) common.Order {
	result := common.CompareNumber(pt.Priority, other.Priority)

	if result == common.Equal {
		return common.CompareNumber(pt.ReplicaID, other.ReplicaID)
	}

	return result
}

func (pt *PointTag) Compare(other PointTag) common.Order {
	result := pt.CompareBase(other)

	if result == common.Equal {
		return common.CompareNumber(pt.Nonce, other.Nonce)
	}

	return result
}

var (
	MinTag = PointTag{
		Priority:  common.UInt32Min,
		ReplicaID: 0,
		Nonce:     1,
	}
	MidTag = PointTag{
		Priority:  common.UInt32Mid,
		ReplicaID: 0,
		Nonce:     3,
	}
	MaxTag = PointTag{
		Priority:  common.UInt32Max,
		ReplicaID: 0,
		Nonce:     2,
	}
)
