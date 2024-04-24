package common

type FatalError string

const (
	NoIntersection                   FatalError = "NoIntersection"
	ExistingSpanOverwrite            FatalError = "ExistingSpanOverwrite"
	UnAppendable                     FatalError = "UnAppendable"
	UnPrependable                    FatalError = "UnPrependable"
	InvalidDistanceBetweenNoRelation FatalError = "InvalidDistanceBetweenNoRelation"
)

func (fatal FatalError) Error() string {
	return string(fatal)
}
