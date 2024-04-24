package text

import (
	"encoding/json"

	"github.com/notebox/nb-crdt-go/text/span"
	"github.com/notebox/nb-crdt-go/text/tree"
)

type Text struct {
	Node *tree.Node
}

func (t *Text) Spans() []*span.INSSpan {
	if t.Node == nil {
		return nil
	}
	return t.Node.Spans()
}

func (t Text) MarshalJSON() ([]byte, error) {
	spans := t.Spans()
	if spans == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(t.Spans())
}

func (t *Text) UnmarshalJSON(data []byte) error {
	var spans []*span.INSSpan
	err := json.Unmarshal(data, &spans)
	if err != nil {
		return err
	}
	if len(spans) == 0 {
		return nil
	}
	t.Node = tree.NewFromSpans(spans)
	return nil
}

// assuming content is not meta
func (t *Text) String() string {
	result := ""
	for _, span := range t.Spans() {
		result += span.Content.Text()
	}
	return result
}

func (t *Text) INS(span *span.INSSpan) error {
	if t.Node != nil {
		err := t.Node.INS(span, 0)
		if err != nil {
			return err
		}
		t.Node = t.Node.Balance()
	} else {
		t.Node = tree.New(*span, nil, nil)
	}
	return nil
}

func (t *Text) DEL(span *span.DELSpan) error {
	if t.Node != nil {
		err := t.Node.DEL(span, 0)
		if err != nil {
			return err
		}
		t.Node = t.Node.Balance()
	}
	return nil
}

func (t *Text) FMT(span *span.FMTSpan) error {
	if t.Node != nil {
		err := t.Node.FMT(span, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Text) MOD(span *span.MODSpan) error {
	if t.Node != nil {
		err := t.Node.MOD(span, 0)
		if err != nil {
			return err
		}
	}
	return nil
}
