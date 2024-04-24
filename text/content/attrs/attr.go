package attrs

import (
	"encoding/json"

	"github.com/notebox/nb-crdt-go/common"
)

type Attr struct {
	Length uint32
	Props  TextProps
	Stamp  *common.Stamp
}

func (leaf Attr) MarshalJSON() ([]byte, error) {
	if leaf.Stamp != nil {
		return json.Marshal([]any{leaf.Length, leaf.Props, leaf.Stamp})
	} else if leaf.Props != nil && len(leaf.Props) > 0 {
		return json.Marshal([]any{leaf.Length, leaf.Props})
	}
	return json.Marshal([]any{leaf.Length})
}

func (leaf *Attr) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[0], &leaf.Length)
	if err != nil {
		return err
	}
	if len(raw) > 1 {
		err = json.Unmarshal(raw[1], &leaf.Props)
		if err != nil {
			return err
		}
	}
	if len(raw) > 2 {
		err = json.Unmarshal(raw[2], &leaf.Stamp)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO deprecated
// func (l *Attr) Clone() *Attr {
// 	props := make(TextProps)
// 	maps.Copy(l.Props, props)

// 	return &Attr{
// 		Length: l.Length,
// 		Props:  props,
// 		Stamp:  l.Stamp,
// 	}
// }

func (l *Attr) EqualsExceptForLength(other *Attr) bool {
	if (l.Stamp != nil) && (other.Stamp != nil) && *l.Stamp != *other.Stamp {
		return false
	}
	if len(l.Props) != len(other.Props) {
		return false
	}
	for k, v := range l.Props {
		if other.Props[k] != v {
			return false
		}
	}
	return true
}

func (l *Attr) Apply(props TextProps, stamp *common.Stamp) {
	for k, v := range props {
		if v == nil {
			delete(l.Props, k)
		} else {
			l.Props[k] = v
		}
	}
	l.Stamp = stamp
}

type TextProps = map[string]any
