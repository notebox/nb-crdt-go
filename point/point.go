package point

import (
	"slices"

	"github.com/notebox/nb-crdt-go/common"
)

// assuming no empty case
type Point []PointTag

func (p Point) Depth() int {
	return len(p)
}

func (p Point) ReplicaID() common.ReplicaID {
	return p[len(p)-1].ReplicaID
}

func (p Point) Nonce() uint32 {
	return p[len(p)-1].Nonce
}

func (p Point) Clone() Point {
	return slices.Clone(p)
}

func (p Point) WithNonce(nonce common.Nonce) Point {
	lastIDX := len(p) - 1
	return append(slices.Clone(p[:lastIDX]), p[lastIDX].WithNonce(nonce))
}

func (p Point) Offset(offset uint32) Point {
	return p.WithNonce(p.Nonce() + offset)
}

func (p Point) Equals(other Point) bool {
	return p.ReplicaID() == other.ReplicaID() && p.Nonce() == other.Nonce()
}

func (p Point) CompareBase(other Point) common.Order {
	if p.Equals(other) {
		return common.Equal
	}

	i := 0
	baseCmp := common.Equal
	for i < min(p.Depth(), other.Depth())-1 && baseCmp == common.Equal {
		baseCmp = p[i].Compare(other[i])
		i++
	}
	if baseCmp == common.Equal {
		baseCmp = p[i].CompareBase(other[i])
	}

	switch baseCmp {
	case common.Equal:
		if p.Depth() == other.Depth() {
			return common.Equal
		}

		if p.Depth() > other.Depth() {
			return common.Tagging
		} else {
			return common.Tagged
		}
	default: // Less or Greater
		return baseCmp
	}
}

func (p Point) Compare(other Point) common.Order {
	if p.Equals(other) {
		return common.Equal
	}

	minimum := min(p.Depth(), other.Depth())
	for i := 0; i < minimum; i++ {
		result := p[i].Compare(other[i])

		if result != common.Equal {
			return result
		}
	}

	if p.Depth() < other.Depth() {
		return common.Less
	} else {
		return common.Greater
	}
}

func (p Point) DistanceFrom(other Point) (uint32, common.Order, error) {
	cmpBase := p.CompareBase(other)
	if cmpBase == common.Less || cmpBase == common.Greater {
		return 0, common.Equal, common.InvalidDistanceBetweenNoRelation
	}

	depth := min(other.Depth(), p.Depth())
	nonce := p[depth-1].Nonce
	otherNonce := other[depth-1].Nonce

	var distance uint32
	if nonce > otherNonce {
		distance = nonce - otherNonce
	} else {
		distance = otherNonce - nonce
	}

	return uint32(distance), common.CompareNumber(nonce, otherNonce), nil
}
