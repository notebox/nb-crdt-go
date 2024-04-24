package content

import (
	"encoding/json"
	"reflect"

	"github.com/notebox/nb-crdt-go/text/content/attrs"
)

type FMTContent struct {
	length uint32
	attrs  attrs.Attrs
}

func NewFMTContent(length uint32, props attrs.TextProps) *FMTContent {
	return &FMTContent{length, attrs.Attrs{attrs.Attr{length, props, nil}}}
}

func (c FMTContent) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{c.length, c.attrs})
}

func (c *FMTContent) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[0], &c.length)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[1], &c.attrs)
	if err != nil {
		return err
	}
	return nil
}

func (c *FMTContent) Length() uint32 {
	return c.length
}

func (c *FMTContent) Attrs() attrs.Attrs {
	return c.attrs
}

func (c *FMTContent) Concat(other *FMTContent) *FMTContent {
	return &FMTContent{c.length + other.Length(), c.attrs.Concat(other.Attrs())}
}

func (c *FMTContent) Slice(start, end uint32) *FMTContent {
	// defense code
	if start >= c.length {
		return &FMTContent{}
	}

	if end >= c.length {
		end = c.length
	}

	return &FMTContent{end - start, c.attrs.Slice(start, end)}
}

func (c *FMTContent) Equals(other *FMTContent) bool {
	return c.length == other.Length() && reflect.DeepEqual(c.attrs, other.Attrs())
}
