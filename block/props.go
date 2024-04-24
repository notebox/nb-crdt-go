package block

import (
	"encoding/json"

	"github.com/notebox/nb-crdt-go/common"
)

type Props map[string]Prop
type Prop struct {
	Nested map[string]Prop
	Stamp  *common.Stamp
	Value  any
}

func (prop *Prop) IsNotLeaf() bool {
	return prop.Nested != nil && prop.Nested["LEAF"].Value != true
}

// assuming LEAF flagged nested is not used
func (prop Prop) MarshalJSON() ([]byte, error) {
	if prop.IsNotLeaf() {
		return json.Marshal(prop.Nested)
	}
	if prop.Value != nil {
		return json.Marshal([]any{prop.Stamp, prop.Value})
	}
	return json.Marshal([]any{prop.Stamp})
}

// assuming LEAF flagged nested is not used
func (prop *Prop) UnmarshalJSON(data []byte) error {
	if data[0] == '{' {
		return json.Unmarshal(data, &prop.Nested)
	}
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[0], &prop.Stamp)
	if err != nil {
		return err
	}
	if len(raw) == 2 {
		return json.Unmarshal(raw[1], &prop.Value)
	}
	return nil
}

type PropsDelta = map[string]PropDelta
type PropDelta struct {
	Nested map[string]PropDelta
	Value  any
}

func (prop *PropDelta) IsNotLeaf() bool {
	return prop.Nested != nil && prop.Nested["LEAF"].Value != true
}

// assuming LEAF flagged nested is not used
func (prop PropDelta) MarshalJSON() ([]byte, error) {
	if prop.IsNotLeaf() {
		return json.Marshal(prop.Nested)
	}
	return json.Marshal(prop.Value)
}

// assuming LEAF flagged nested is not used
func (prop *PropDelta) UnmarshalJSON(data []byte) error {
	if data[0] == '{' {
		return json.Unmarshal(data, &prop.Nested)
	}
	return json.Unmarshal(data, &prop.Value)
}

// assuming LEAF flagged nested is not used and delta is not nil
func UpdateProps(props Props, delta PropsDelta, stamp common.Stamp) Props {
	if props == nil {
		props = make(Props)
	}

	for key, d := range delta {
		if d.IsNotLeaf() {
			props[key] = Prop{
				Nested: UpdateProps(props[key].Nested, d.Nested, stamp),
				Stamp:  &stamp,
			}
		} else if props[key].Stamp.IsOlderThan(&stamp) {
			props[key] = Prop{
				Value: d.Value,
				Stamp: &stamp,
			}
		}
	}

	return props
}

const (
	DEL_PROP_KEY = "DEL"
	MOV_PROP_KEY = "MOV"
)
