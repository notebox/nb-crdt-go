package content

import (
	"encoding/json"
)

type DELContent struct {
	length uint32
}

func NewDELContent(length uint32) *DELContent {
	return &DELContent{length}
}

func (c DELContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.length)
}

func (c *DELContent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.length)
}

func (c *DELContent) Length() uint32 {
	return c.length
}

func (c *DELContent) Equals(other *DELContent) bool {
	return c.length == other.Length()
}

func (c *DELContent) Concat(other *DELContent) *DELContent {
	return &DELContent{c.length + other.Length()}
}

func (c *DELContent) Slice(start, end uint32) *DELContent {
	// defense code
	if start >= c.length {
		return &DELContent{}
	}

	if end >= c.length {
		return &DELContent{c.length - start}
	}

	return &DELContent{end - start}
}
