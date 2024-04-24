package block

import (
	"encoding/json"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/point"
	"github.com/notebox/nb-crdt-go/text/span"
)

type Contribution struct {
	BlockID    common.BlockID `json:"blockID"`
	Nonce      ReplicaNonce   `json:"nonce"`
	Stamp      common.Stamp   `json:"stamp"`
	Operations Operations     `json:"ops"`
}

type BlockPoint struct {
	ParentBlockID common.BlockID
	Point         point.Point
}

func (bp BlockPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]any{bp.ParentBlockID, bp.Point})
}

func (bp *BlockPoint) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[0], &bp.ParentBlockID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[1], &bp.Point)
	if err != nil {
		return err
	}
	return nil
}

type Operations struct {
	BINS *Block
	BDEL *bool
	BSET PropsDelta
	BMOV *BlockPoint
	TINS []*span.INSSpan
	TDEL []*span.DELSpan
	TFMT []*span.FMTSpan
	TMOD []*span.MODSpan
}
