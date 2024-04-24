package tree

import (
	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/text/span"
)

// base AVL Node
type Node struct {
	Span   span.INSSpan
	Left   *Node
	Right  *Node
	Length uint32
	Rank   int

	ShouldBeDeleted bool
}

func NewFromSpans(spans []*span.INSSpan) *Node {
	var left, right *Node
	c := len(spans)
	midIDX := c / 2
	if midIDX > 0 {
		left = NewFromSpans(spans[0:midIDX])
	}
	if midIDX < c-1 {
		right = NewFromSpans(spans[midIDX+1:])
	}
	return New(*spans[midIDX], left, right)
}

func New(span span.INSSpan, left *Node, right *Node) *Node {
	node := &Node{
		Span:  span,
		Left:  left,
		Right: right,
	}
	node.update()
	return node
}

func (n *Node) Spans() []*span.INSSpan {
	spans := make([]*span.INSSpan, 0)
	if n.Left != nil {
		spans = append(spans, n.Left.Spans()...)
	}
	spans = append(spans, &n.Span)
	if n.Right != nil {
		spans = append(spans, n.Right.Spans()...)
	}
	return spans
}

func (n *Node) Balance() *Node {
	if n.ShouldBeDeleted {
		return nil
	}

	node := n
	for node.isRightUnbalanced() {
		node = node.rotateLeft()
	}

	for node.isLeftUnbalanced() {
		if node.Left.isRightOriented() {
			node.Left = node.Left.rotateLeft()
			node.update()
		}
		node = node.rotateRight()
	}

	return node
}

func (n *Node) predecessorSpan() *span.INSSpan {
	if n.Left != nil {
		return n.Left.maxSpan()
	}
	return nil
}

func (n *Node) successorSpan() *span.INSSpan {
	if n.Right != nil {
		return n.Right.minSpan()
	}
	return nil
}

func (n *Node) insertPredecessor(span span.INSSpan) {
	if n.Left != nil {
		n.Left = n.Left.insertMax(span)
	} else {
		n.Left = New(span, nil, nil)
	}
	n.update()
}

func (n *Node) insertSuccessor(span span.INSSpan) {
	if n.Right != nil {
		n.Right = n.Right.insertMin(span)
	} else {
		n.Right = New(span, nil, nil)
	}
	n.update()
}

func (n *Node) deletePredecessor() {
	if n.Left != nil {
		n.Left = n.Left.deleteMax()
		n.update()
	}
}

func (n *Node) deleteSuccessor() {
	if n.Right != nil {
		n.Right = n.Right.deleteMin()
		n.update()
	}
}

func (n *Node) deleteSelf() {
	if succ := n.successorSpan(); succ != nil {
		n.deleteSuccessor()
		n.Span = *succ
		n.mergeLeft()
		return
	}

	if prev := n.predecessorSpan(); prev != nil {
		n.deletePredecessor()
		n.Span = *prev
		n.mergeRight()
		return
	}

	n.ShouldBeDeleted = true
}

// assuming content is not meta
func (n *Node) mergeLeft() {
	curr := n.Span
	if prev := n.predecessorSpan(); prev != nil {
		if cmp, err := prev.Compare(&curr); err == nil && cmp == common.Prependable {
			n.deletePredecessor()
			span := prev.Append(&curr)
			n.Span = *span
		}
	}
}

// assuming content is not meta
func (n *Node) mergeRight() {
	curr := n.Span
	if succ := n.successorSpan(); succ != nil {
		if cmp, err := succ.Compare(&curr); err == nil && cmp == common.Appendable {
			n.deleteSuccessor()
			span := curr.Append(succ)
			n.Span = *span
		}
	}
}

func (n *Node) leftRank() int {
	if n.Left == nil {
		return 0
	}
	return n.Left.Rank
}

func (n *Node) rightRank() int {
	if n.Right == nil {
		return 0
	}
	return n.Right.Rank
}

func (n *Node) leftLength() uint32 {
	if n.Left == nil {
		return 0
	}
	return n.Left.Length
}

func (n *Node) rightLength() uint32 {
	if n.Right == nil {
		return 0
	}
	return n.Right.Length
}

func (n *Node) update() {
	n.Rank = 1 + max(n.leftRank(), n.rightRank())
	n.Length = n.leftLength() + n.Span.Length() + n.rightLength()
}

func (n *Node) balanceFactor() int {
	return n.rightRank() - n.leftRank()
}

func (n *Node) isRightUnbalanced() bool {
	return n.balanceFactor() > 1
}

func (n *Node) isLeftUnbalanced() bool {
	return n.balanceFactor() < -1
}

func (n *Node) isRightOriented() bool {
	return n.balanceFactor() == 1
}

func (n *Node) rotateLeft() *Node {
	right := n.Right
	n.Right = right.Left
	n.update()
	right.Left = n
	right.update()
	return right
}

func (n *Node) rotateRight() *Node {
	left := n.Left
	n.Left = left.Right
	n.update()
	left.Right = n
	left.update()
	return left
}

func (n *Node) minSpan() *span.INSSpan {
	if n.Left != nil {
		return n.Left.minSpan()
	}
	return &n.Span
}

func (n *Node) maxSpan() *span.INSSpan {
	if n.Right != nil {
		return n.Right.maxSpan()
	}
	return &n.Span
}

func (n *Node) insertMin(span span.INSSpan) *Node {
	if n.Left != nil {
		n.Left = n.Left.insertMin(span)
	} else {
		n.Left = New(span, nil, nil)
	}
	n.update()
	return n.Balance()
}

func (n *Node) insertMax(span span.INSSpan) *Node {
	if n.Right != nil {
		n.Right = n.Right.insertMax(span)
	} else {
		n.Right = New(span, nil, nil)
	}
	n.update()
	return n.Balance()
}

func (n *Node) deleteMin() *Node {
	if n.Left == nil {
		return n.Right
	}
	n.Left = n.Left.deleteMin()
	n.update()
	return n.Balance()
}

func (n *Node) deleteMax() *Node {
	if n.Right == nil {
		return n.Left
	}
	n.Right = n.Right.deleteMax()
	n.update()
	return n.Balance()
}
