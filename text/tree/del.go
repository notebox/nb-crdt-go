package tree

import (
	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/span"
)

func (n *Node) DEL(span *span.DELSpan, minIndex uint32) error {
	curr := n.Span
	currIndex := minIndex + n.leftLength()
	nextMinIndex := currIndex + curr.Length()
	cmp, err := curr.Compare(span)
	if err != nil {
		return err
	}

	switch cmp {
	case common.Less, common.Prependable:
		return n.delRight(span, nextMinIndex)
	case common.Greater, common.Appendable:
		return n.delLeft(span, minIndex)
	case common.IncludingLeft:
		seg, err := curr.AppendableSegmentTo(span)
		if err != nil {
			return err
		}
		n.Span = *seg
		return err
	case common.IncludingRight:
		seg, err := curr.PrependableSegmentTo(span)
		if err != nil {
			return err
		}
		n.Span = *seg
	case common.IncludingMiddle:
		seg, err := curr.PrependableSegmentTo(span)
		if err != nil {
			return err
		}
		n.Span = *seg
		seg, err = curr.AppendableSegmentTo(span)
		if err != nil {
			return err
		}
		n.insertSuccessor(*seg)
		return nil
	case common.RightOverlap:
		err := n.delRight(span, nextMinIndex)
		if err != nil {
			return err
		}
		overlapped, err := span.Intersection(&curr)
		if err != nil {
			return err
		}
		return n.DEL(overlapped, minIndex)
	case common.LeftOverlap:
		overlapped, err := span.Intersection(&curr)
		if err != nil {
			return err
		}
		err = n.DEL(overlapped, minIndex)
		if err != nil {
			return err
		}
		return n.delLeft(span, minIndex)
	case common.Equal:
		n.deleteSelf()
		return nil
	case common.IncludedLeft:
		err := n.delRight(span, nextMinIndex)
		if err != nil {
			return err
		}
		n.deleteSelf()
		return nil
	case common.IncludedRight:
		err := n.delLeft(span, minIndex)
		if err != nil {
			return err
		}
		n.deleteSelf()
		return nil
	case common.IncludedMiddle:
		err := n.delRight(span, nextMinIndex)
		if err != nil {
			return err
		}
		err = n.delLeft(span, minIndex)
		if err != nil {
			return err
		}
		n.deleteSelf()
		return nil
	case common.Splitting:
		left, right, err := span.SplitWith(&curr)
		if err != nil {
			return err
		}
		err = n.DEL(right, minIndex)
		if err != nil {
			return err
		}
		err = n.DEL(left, minIndex)
		if err != nil {
			return err
		}
		return nil
	case common.Splitted:
		return nil
	}
	return nil
}

func (n *Node) delLeft(span *span.DELSpan, minIndex uint32) error {
	if n.Left == nil {
		return nil
	}

	err := n.Left.DEL(span, minIndex)
	if err != nil {
		return err
	}
	n.Left = n.Left.Balance()
	return nil
}

func (n *Node) delRight(span *span.DELSpan, minIndex uint32) error {
	if n.Right == nil {
		return nil
	}

	err := n.Right.DEL(span, minIndex)
	if err != nil {
		return err
	}
	n.Right = n.Right.Balance()
	return nil
}
