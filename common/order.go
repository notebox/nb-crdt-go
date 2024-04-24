package common

type Order int

const (
	Splitting Order = iota + 1
	Tagging
	Less
	Prependable
	RightOverlap
	IncludingRight
	IncludingMiddle
	IncludingLeft
	Equal
	IncludedLeft
	IncludedMiddle
	IncludedRight
	LeftOverlap
	Appendable
	Greater
	Tagged
	Splitted
)
