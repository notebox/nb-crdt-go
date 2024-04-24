package content

import (
	"encoding/json"

	"github.com/notebox/nb-crdt-go/common"
)

type MODContent struct {
	text string
}

func NewMODContent(text string) *MODContent {
	return &MODContent{text: text}
}

func (c MODContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.text)
}

func (c *MODContent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.text)
}

func (c *MODContent) Length() uint32 {
	return common.UTF16Length(c.text)
}

func (c *MODContent) Text() string {
	return c.text
}

func (c *MODContent) Concat(other *MODContent) *MODContent {
	return &MODContent{c.text + other.Text()}
}

func (c *MODContent) Slice(start, end uint32) *MODContent {
	return &MODContent{common.UTF16Slice(c.text, start, end)}
}

func (c *MODContent) Equals(other *MODContent) bool {
	return c.text == other.Text()
}
