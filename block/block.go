package block

import (
	"encoding/json"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/point"
	"github.com/notebox/nb-crdt-go/text"
)

type Block struct {
	BlockID       common.BlockID
	Version       Version
	Props         Props
	Text          *text.Text
	ParentBlockID *common.BlockID
	Point         point.Point
	IsDeleted     bool
}

func (block Block) MarshalJSON() ([]byte, error) {
	arr := []any{block.BlockID, block.Version, block.Point, block.Props, block.IsDeleted}
	if block.ParentBlockID != nil {
		arr = append(arr, block.Text, block.ParentBlockID)
	} else if block.Text != nil {
		arr = append(arr, block.Text)
	}
	return json.Marshal(arr)
}

func (block *Block) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[0], &block.BlockID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[1], &block.Version)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[2], &block.Point)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[3], &block.Props)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[4], &block.IsDeleted)
	if err != nil {
		return err
	}
	switch len(raw) {
	case 7:
		err = json.Unmarshal(raw[6], &block.ParentBlockID)
		if err != nil {
			return err
		}
		fallthrough
	case 6:
		err = json.Unmarshal(raw[5], &block.Text)
		if err != nil {
			return err
		}
	}
	return nil
}

func (block *Block) Apply(ctrb Contribution) error {
	if ctrb.Operations.BDEL != nil {
		block.IsDeleted = *ctrb.Operations.BDEL
		block.updateProp(DEL_PROP_KEY, &ctrb.Stamp)
	}

	if ctrb.Operations.BSET != nil {
		block.Props = UpdateProps(block.Props, ctrb.Operations.BSET, ctrb.Stamp)
	}

	if ctrb.Operations.BMOV != nil {
		block.ParentBlockID = &ctrb.Operations.BMOV.ParentBlockID
		block.Point = ctrb.Operations.BMOV.Point
		block.updateProp(MOV_PROP_KEY, &ctrb.Stamp)
	}

	if block.Text != nil {
		for _, ins := range ctrb.Operations.TINS {
			err := block.Text.INS(ins)
			if err != nil {
				return err
			}
		}

		for _, del := range ctrb.Operations.TDEL {
			err := block.Text.DEL(del)
			if err != nil {
				return err
			}
		}

		for _, fmt := range ctrb.Operations.TFMT {
			err := block.Text.FMT(fmt)
			if err != nil {
				return err
			}
		}

		for _, mod := range ctrb.Operations.TMOD {
			err := block.Text.MOD(mod)
			if err != nil {
				return err
			}
		}
	}

	block.Version[ctrb.Stamp.ReplicaID] = ctrb.Nonce
	return nil
}

func (block *Block) updateProp(key string, stamp *common.Stamp) {
	if block.Props == nil {
		block.Props = make(Props)
	}
	block.Props[key] = Prop{Stamp: stamp}
}
