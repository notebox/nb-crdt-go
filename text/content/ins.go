package content

import (
	"encoding/json"
	"reflect"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/content/attrs"
)

type Meta map[string]*string

type INSContent struct {
	text  string
	attrs attrs.Attrs
}

func NewINSContent(text string) *INSContent {
	return &INSContent{text, attrs.Attrs{{Length: common.UTF16Length(text)}}}
}

func (c INSContent) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{c.attrs, c.text})
}

func (c *INSContent) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[0], &c.attrs)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[1], &c.text)
	if err != nil {
		return err
	}
	return nil
}

func (c *INSContent) Length() uint32 {
	return common.UTF16Length(c.text)
}

func (c *INSContent) Text() string {
	return c.text
}

func (c *INSContent) Attrs() attrs.Attrs {
	return c.attrs
}

// assuming content is not meta
func (c *INSContent) Concat(other *INSContent) *INSContent {
	return &INSContent{c.text + other.Text(), c.attrs.Concat(other.Attrs())}
}

// assuming content is not meta
func (c *INSContent) Slice(start uint32, end uint32) *INSContent {
	return &INSContent{common.UTF16Slice(c.text, start, end), c.attrs.Slice(start, end)}
}

func (c *INSContent) Equals(other *INSContent) bool {
	return c.text == other.Text() && reflect.DeepEqual(c.attrs, other.Attrs())
}

func (c *INSContent) FMT(index uint32, other *FMTContent) {
	left := c.attrs.Slice(0, index)
	affected := c.attrs.Slice(index, index+other.Length())
	right := c.attrs.Slice(index+other.Length(), c.Length())

	affected.Merge(other.Attrs())

	n := left.Concat(affected).Concat(right)
	c.attrs = n
}

// assuming content is not meta
func (c *INSContent) MOD(index uint32, other *MODContent) {
	c.text = common.UTF16Slice(c.text, 0, index) +
		other.Text() +
		common.UTF16Slice(c.text, index+other.Length(), c.Length())
}
