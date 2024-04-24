package tree

import (
	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/span"
)

func (n *Node) MOD(span *span.MODSpan, minIndex uint32) error {
	curr := n.Span
	currIndex := minIndex + n.leftLength()
	nextMinIndex := currIndex + curr.Length()
	cmp, err := curr.Compare(span)
	if err != nil {
		return err
	}

	switch cmp {
	case common.Less, common.Prependable:
		return n.modRight(span, nextMinIndex)
	case common.Greater, common.Appendable:
		return n.modLeft(span, minIndex)
	case common.IncludingLeft, common.IncludingRight, common.IncludingMiddle, common.Equal:
		subIndex := span.LowerPoint().Nonce() - curr.LowerPoint().Nonce()
		curr.Content.MOD(subIndex, span.Content)
		return nil
	case common.RightOverlap:
		err := n.modRight(span, nextMinIndex)
		if err != nil {
			return err
		}
		subIndex := span.LowerPoint().Nonce() - curr.LowerPoint().Nonce()
		content := span.Content.Slice(
			0,
			curr.UpperPoint().Nonce()-span.LowerPoint().Nonce()+1,
		)
		curr.Content.MOD(subIndex, content)
		return nil
	case common.LeftOverlap:
		err := n.modLeft(span, minIndex)
		if err != nil {
			return err
		}
		content := span.Content.Slice(
			curr.LowerPoint().Nonce()-span.LowerPoint().Nonce(),
			span.Content.Length(),
		)
		curr.Content.MOD(0, content)
		return nil
	case common.IncludedLeft:
		err := n.modRight(span, nextMinIndex)
		if err != nil {
			return err
		}
		content := span.Content.Slice(0, curr.Length())
		curr.Content.MOD(0, content)
		return nil
	case common.IncludedRight:
		err := n.modLeft(span, minIndex)
		if err != nil {
			return err
		}
		content := span.Content.Slice(
			span.Length()-curr.Length(),
			span.Length(),
		)
		curr.Content.MOD(0, content)
		return nil
	case common.IncludedMiddle:
		err := n.modLeft(span, minIndex)
		if err != nil {
			return err
		}
		err = n.modRight(span, nextMinIndex)
		if err != nil {
			return err
		}
		startIndex := curr.LowerPoint().Nonce() - span.LowerPoint().Nonce()
		content := span.Content.Slice(
			startIndex,
			startIndex+curr.Length(),
		)
		curr.Content.MOD(0, content)
		return nil
	case common.Splitting, common.Splitted:
		return nil
	}
	return nil
}

func (n *Node) modLeft(span *span.MODSpan, minIndex uint32) error {
	if n.Left == nil {
		return nil
	}

	return n.Left.MOD(span, minIndex)
}

func (n *Node) modRight(span *span.MODSpan, minIndex uint32) error {
	if n.Right == nil {
		return nil
	}

	return n.Right.MOD(span, minIndex)
}
