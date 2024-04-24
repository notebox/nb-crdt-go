package tree

import (
	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/span"
)

func (n *Node) INS(span *span.INSSpan, minIndex uint32) error {
	curr := n.Span
	currStartIndex := minIndex + n.leftLength()
	nextMinIndex := currStartIndex + curr.Length()
	cmp, err := curr.Compare(span)
	if err != nil {
		return err
	}

	switch cmp {
	case common.Less:
		return n.insIntoRight(span, nextMinIndex)
	case common.Greater:
		return n.insIntoLeft(span, minIndex)
	case common.Prependable:
		err := n.insIntoRight(span, nextMinIndex)
		if err != nil {
			return err
		}
		n.mergeRight()
		return nil
	case common.Appendable:
		err := n.insIntoLeft(span, minIndex)
		if err != nil {
			return err
		}
		n.mergeLeft()
		return nil
	case common.Splitted:
		left, right, err := curr.SplitWith(span)
		if err != nil {
			return err
		}
		n.insertPredecessor(*left)
		n.Span = *span
		n.insertSuccessor(*right)
		return nil
	case common.Splitting:
		left, right, err := span.SplitWith(&curr)
		if err != nil {
			return err
		}
		err = n.insIntoRight(right, nextMinIndex)
		if err != nil {
			return err
		}
		err = n.insIntoLeft(left, minIndex)
		if err != nil {
			return err
		}
		return nil
	default:
		return common.ExistingSpanOverwrite
	}
}

func (n *Node) insIntoLeft(span *span.INSSpan, minIndex uint32) error {
	if n.Left != nil {
		err := n.Left.INS(span, minIndex)
		if err != nil {
			return err
		}
		n.Left = n.Left.Balance()
		return nil
	}

	n.Left = New(*span, nil, nil)
	return nil
}

func (n *Node) insIntoRight(span *span.INSSpan, minIndex uint32) error {
	if n.Right != nil {
		err := n.Right.INS(span, minIndex)
		if err != nil {
			return err
		}
		n.Right = n.Right.Balance()
		return nil
	}

	n.Right = New(*span, nil, nil)
	return nil
}
