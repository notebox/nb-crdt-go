package attrs

import (
	"slices"

	"github.com/notebox/nb-crdt-go/common"
)

type Attrs []Attr

func (attrs Attrs) Clone() Attrs {
	return slices.Clone(attrs)
}

func (attrs Attrs) Concat(other Attrs) Attrs {
	boundaryIndex := len(attrs)
	leaves := append(attrs, other...)

	if boundaryIndex < 1 || boundaryIndex >= len(leaves) {
		return leaves
	}

	if leaves[boundaryIndex-1].EqualsExceptForLength(&leaves[boundaryIndex]) {
		leaves[boundaryIndex-1].Length += leaves[boundaryIndex].Length
		leaves = append(leaves[:boundaryIndex], leaves[boundaryIndex+1:]...)
	}

	return leaves
}

func (attrs Attrs) Slice(start, end uint32) Attrs {
	if end-start < 1 {
		return make(Attrs, 0)
	}

	leaves := attrs.Clone()
	nextIndex := leaves[0].Length

	for nextIndex <= start {
		leaves = leaves[1:]
		nextIndex += leaves[0].Length
	}
	leaves[0].Length = nextIndex - start

	index := 0
	for nextIndex < end {
		index++
		nextIndex += leaves[index].Length
	}

	if nextIndex > end {
		leaves[index].Length -= (nextIndex - end)
	}

	return leaves[:index+1]
}

// TODO deprecated
// func (attrs *Attrs) Apply(props TextProps, stamp common.Stamp) {
// 	beforeIDX := -1
// 	newLeaves := make([]Attr, 0)
// 	for idx, leaf := range *attrs {
// 		leaf.Apply(props, &stamp)
// 		if beforeIDX > -1 && newLeaves[beforeIDX].EqualsExceptForLength(&leaf) {
// 			newLeaves[beforeIDX].Length += leaf.Length
// 		} else {
// 			newLeaves = append(newLeaves, leaf)
// 		}
// 		beforeIDX = idx
// 	}
// 	*attrs = newLeaves
// }

func (attrs *Attrs) Merge(other Attrs) {
	newLeaves := make([]Attr, 0)
	count := len(*attrs)
	otherCount := len(other)
	if count == 0 {
		return
	}

	idx := 0
	otherIDX := 0
	leaf := (*attrs)[idx]
	otherAttr := other[otherIDX]
	leafLength := leaf.Length
	otherAttrLength := otherAttr.Length

	for leafLength > 0 {
		var newProps *TextProps
		var newStamp *common.Stamp
		var newLength uint32

		if leaf.Stamp.IsOlderThan(otherAttr.Stamp) {
			newProps = &otherAttr.Props
			newStamp = otherAttr.Stamp
		} else {
			newProps = &leaf.Props
			newStamp = leaf.Stamp
		}

		if leafLength < otherAttrLength {
			newLength = leafLength
		} else {
			newLength = otherAttrLength
		}

		newAttr := Attr{Length: newLength, Props: *newProps, Stamp: newStamp}
		lastNewAttrIndex := len(newLeaves) - 1
		if lastNewAttrIndex > -1 && newLeaves[lastNewAttrIndex].EqualsExceptForLength(&newAttr) {
			newLeaves[lastNewAttrIndex].Length += newAttr.Length
		} else {
			newLeaves = append(newLeaves, newAttr)
		}

		leafLength -= newLength
		otherAttrLength -= newLength

		if leafLength < 1 {
			idx++
			if idx < count {
				leaf = (*attrs)[idx]
				leafLength = leaf.Length
			} else {
				leafLength = 0
			}
		}

		if otherAttrLength < 1 {
			otherIDX++
			if otherIDX < otherCount {
				otherAttr = other[otherIDX]
				otherAttrLength = otherAttr.Length
			} else {
				otherAttrLength = 0
			}
		}
	}
	*attrs = newLeaves
}
