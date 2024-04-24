package test

import (
	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nb-crdt-go/point"
	"github.com/notebox/nb-crdt-go/text/content"
	"github.com/notebox/nb-crdt-go/text/content/attrs"
	"github.com/notebox/nb-crdt-go/text/span"
)

type RawCase struct {
	Point [][3]uint32
	Text  string
}

var Cases = map[common.Order]RawCase{
	common.Splitted:        {[][3]uint32{{1, 1, 1}, {5, 5, 5}, {8, 8, 8}}, "ghi"},
	common.Less:            {[][3]uint32{{1, 1, 1}, {5, 5, 9}}, "89a"},
	common.Prependable:     {[][3]uint32{{1, 1, 1}, {5, 5, 8}}, "789"},
	common.RightOverlap:    {[][3]uint32{{1, 1, 1}, {5, 5, 7}}, "678"},
	common.IncludingRight:  {[][3]uint32{{1, 1, 1}, {5, 5, 7}}, "6"},
	common.IncludingMiddle: {[][3]uint32{{1, 1, 1}, {5, 5, 6}}, "5"},
	common.IncludingLeft:   {[][3]uint32{{1, 1, 1}, {5, 5, 5}}, "4"},
	common.Equal:           {[][3]uint32{{1, 1, 1}, {5, 5, 5}}, "456"},
	common.IncludedLeft:    {[][3]uint32{{1, 1, 1}, {5, 5, 5}}, "4567"},
	common.IncludedMiddle:  {[][3]uint32{{1, 1, 1}, {5, 5, 4}}, "34567"},
	common.IncludedRight:   {[][3]uint32{{1, 1, 1}, {5, 5, 4}}, "3456"},
	common.LeftOverlap:     {[][3]uint32{{1, 1, 1}, {5, 5, 3}}, "234"},
	common.Appendable:      {[][3]uint32{{1, 1, 1}, {5, 5, 2}}, "123"},
	common.Greater:         {[][3]uint32{{1, 1, 1}, {5, 5, 1}}, "012"},
	common.Splitting:       {[][3]uint32{{1, 1, 1}}, "jkl"},
}

func PointFrom(raw [][3]uint32) point.Point {
	var pts []point.PointTag
	for _, pt := range raw {
		pts = append(pts, point.PointTag{Priority: pt[0], ReplicaID: pt[1], Nonce: pt[2]})
	}
	return pts
}

func INSSpanFrom(raw RawCase) *span.INSSpan {
	s := span.New(PointFrom(raw.Point), content.NewINSContent(raw.Text))
	return &s
}

func MODSpanFrom(raw RawCase) *span.MODSpan {
	s := span.New(PointFrom(raw.Point), content.NewMODContent(raw.Text))
	return &s
}

func DELSpanFrom(raw RawCase) *span.DELSpan {
	s := span.New(PointFrom(raw.Point), content.NewDELContent(common.UTF16Length(raw.Text)))
	return &s
}

func FMTSpanFrom(raw RawCase, props attrs.TextProps) *span.FMTSpan {
	s := span.New(PointFrom(raw.Point), content.NewFMTContent(common.UTF16Length(raw.Text), props))
	return &s
}
