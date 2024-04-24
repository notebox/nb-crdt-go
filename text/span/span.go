package span

import (
	"encoding/json"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/point"
	"github.com/notebox/nb-crdt-go/text/content"
)

type Span[C content.Content[C]] struct {
	Content C

	lowerPoint point.Point
}

func New[C content.Content[C]](lowerPoint point.Point, content C) Span[C] {
	return Span[C]{lowerPoint: lowerPoint, Content: content}
}

func (s Span[C]) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{s.lowerPoint, s.Content})
}

func (s *Span[C]) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[0], &s.lowerPoint)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[1], &s.Content)
	if err != nil {
		return err
	}
	return nil
}

func (s *Span[C]) Equals(other *Span[C]) bool {
	return s.lowerPoint.Equals(other.LowerPoint()) && s.Content.Equals(other.Content)
}

func (s *Span[C]) LowerPoint() point.Point {
	return s.lowerPoint
}

func (s *Span[C]) ReplicaID() common.ReplicaID {
	return s.lowerPoint.ReplicaID()
}

func (s *Span[C]) Length() uint32 {
	return s.Content.Length()
}

func (s *Span[C]) UpperPoint() point.Point {
	return s.NthPoint(s.Length() - 1)
}

func (s *Span[C]) NthPoint(nth uint32) point.Point {
	return s.lowerPoint.Offset(nth)
}

func (s *Span[C]) NonceRange() common.ClosedRange {
	return common.ClosedRange{
		Lower:  uint32(s.lowerPoint.Nonce()),
		Length: s.Length(),
	}
}

func (s *Span[C]) Compare(other AnySpan) (common.Order, error) {
	baseCmp := s.lowerPoint.CompareBase(other.LowerPoint())
	if baseCmp == common.Less {
		return common.Less, nil
	}
	if baseCmp == common.Greater {
		return common.Greater, nil
	}

	// EQUAL - only diff nonces
	if baseCmp == common.Equal {
		return s.NonceRange().Compare(other.NonceRange()), nil
	}

	dist, order, err := s.lowerPoint.DistanceFrom(other.LowerPoint())
	if err != nil {
		return common.Order(-1), err
	}

	// TAGGED
	if baseCmp == common.Tagged {
		if order == common.Greater {
			return common.Greater, nil
		}

		if dist >= s.Length()-1 {
			return common.Less, nil
		}

		return common.Splitted, nil
	}

	// TAGGING
	if order == common.Less {
		return common.Less, nil
	}
	if dist >= other.Length()-1 {
		return common.Greater, nil
	}

	return common.Splitting, nil
}

func (s *Span[C]) Append(other *Span[C]) *Span[C] {
	return &Span[C]{
		lowerPoint: s.lowerPoint,
		Content:    s.Content.Concat(other.Content),
	}
}

func (s *Span[C]) LeftSplitAt(index uint32) *Span[C] {
	return &Span[C]{
		lowerPoint: s.lowerPoint,
		Content:    s.Content.Slice(0, index),
	}
}

func (s *Span[C]) RightSplitAt(index uint32) *Span[C] {
	return &Span[C]{
		lowerPoint: s.NthPoint(index),
		Content:    s.Content.Slice(index, s.Length()),
	}
}

func (s *Span[C]) SplitAt(index uint32) (*Span[C], *Span[C], error) {
	return s.LeftSplitAt(index), s.RightSplitAt(index), nil
}

func (s *Span[C]) SplitWith(other AnySpan) (*Span[C], *Span[C], error) {
	dist, _, err := s.lowerPoint.DistanceFrom(other.LowerPoint())
	if err != nil {
		return nil, nil, err
	}

	return s.SplitAt(dist + 1)
}

func (s *Span[C]) AppendableSegmentTo(other AnySpan) (*Span[C], error) {
	dist, order, err := s.lowerPoint.DistanceFrom(other.LowerPoint())
	if err != nil {
		return nil, err
	}
	if order == common.Less {
		return s.RightSplitAt(other.Length() + dist), nil
	}

	if other.Length() < dist {
		return nil, common.UnAppendable
	}
	index := other.Length() - dist
	return s.RightSplitAt(index), nil
}

func (s *Span[C]) PrependableSegmentTo(other AnySpan) (*Span[C], error) {
	dist, order, err := s.lowerPoint.DistanceFrom(other.LowerPoint())
	if err != nil {
		return nil, err
	}
	if order != common.Less || dist > s.Length() {
		return nil, common.UnPrependable
	}
	return s.LeftSplitAt(dist), nil
}

func (s *Span[C]) Intersection(other AnySpan) (*Span[C], error) {
	cmp, err := s.Compare(other)
	if err != nil {
		return nil, err
	}
	switch cmp {
	case common.Splitted, common.Less, common.Prependable, common.Appendable, common.Greater, common.Splitting:
		return nil, common.NoIntersection
	default:
		break
	}

	dist, order, err := s.lowerPoint.DistanceFrom(other.LowerPoint())
	if err != nil {
		return nil, err
	}
	if order == common.Less {
		end := dist + min(s.Length()-dist, other.Length())
		return &Span[C]{
			lowerPoint: other.LowerPoint(),
			Content:    s.Content.Slice(dist, end),
		}, nil
	}

	end := min(s.Length(), other.Length()-dist)
	return &Span[C]{
		lowerPoint: s.lowerPoint,
		Content:    s.Content.Slice(0, end),
	}, nil
}

type INSSpan = Span[*content.INSContent]
type DELSpan = Span[*content.DELContent]
type FMTSpan = Span[*content.FMTContent]
type MODSpan = Span[*content.MODContent]
type AnySpan interface {
	LowerPoint() point.Point
	NonceRange() common.ClosedRange
	Length() uint32
}
