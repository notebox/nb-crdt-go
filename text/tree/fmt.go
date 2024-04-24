package tree

import (
	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/span"
)

func (n *Node) FMT(span *span.FMTSpan, minIndex uint32) error {
	curr := n.Span
	currIndex := minIndex + n.leftLength()
	nextMinIndex := currIndex + curr.Length()
	cmp, err := curr.Compare(span)
	if err != nil {
		return err
	}

	switch cmp {
	case common.Less, common.Prependable:
		err := n.fmtRight(span, nextMinIndex)
		if err != nil {
			return err
		}
		return nil
	case common.Greater, common.Appendable:
		err := n.fmtLeft(span, minIndex)
		if err != nil {
			return err
		}
		return nil
	case common.IncludingLeft, common.IncludingRight, common.IncludingMiddle, common.Equal:
		subIndex := span.LowerPoint().Nonce() - curr.LowerPoint().Nonce()
		curr.Content.FMT(subIndex, span.Content)
		return nil
	case common.RightOverlap:
		err := n.fmtRight(span, nextMinIndex)
		if err != nil {
			return err
		}
		subIndex := span.LowerPoint().Nonce() - curr.LowerPoint().Nonce()
		content := span.Content.Slice(
			0,
			curr.UpperPoint().Nonce()-span.LowerPoint().Nonce()+1,
		)
		n.Span.Content.FMT(subIndex, content)
		return nil
	case common.LeftOverlap:
		err := n.fmtLeft(span, minIndex)
		if err != nil {
			return err
		}
		content := span.Content.Slice(
			curr.LowerPoint().Nonce()-span.LowerPoint().Nonce(),
			span.Content.Length(),
		)
		n.Span.Content.FMT(0, content)
		return nil
	case common.IncludedLeft:
		err := n.fmtRight(span, nextMinIndex)
		if err != nil {
			return err
		}
		content := span.Content.Slice(0, curr.Length())
		n.Span.Content.FMT(0, content)
		return nil
	case common.IncludedRight:
		err := n.fmtLeft(span, minIndex)
		if err != nil {
			return err
		}
		content := span.Content.Slice(
			span.Length()-curr.Length(),
			span.Length(),
		)
		n.Span.Content.FMT(0, content)
		return nil
	case common.IncludedMiddle:
		err := n.fmtLeft(span, minIndex)
		if err != nil {
			return err
		}
		err = n.fmtRight(span, nextMinIndex)
		if err != nil {
			return err
		}
		startIndex := curr.LowerPoint().Nonce() - span.LowerPoint().Nonce()
		content := span.Content.Slice(
			startIndex,
			startIndex+curr.Length(),
		)
		n.Span.Content.FMT(0, content)
		return nil
	case common.Splitting, common.Splitted:
		return nil
	}
	return nil
}

func (n *Node) fmtLeft(span *span.FMTSpan, minIndex uint32) error {
	if n.Left == nil {
		return nil
	}

	return n.Left.FMT(span, minIndex)
}

func (n *Node) fmtRight(span *span.FMTSpan, minIndex uint32) error {
	if n.Right == nil {
		return nil
	}

	return n.Right.FMT(span, minIndex)
}
