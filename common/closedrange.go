package common

type ClosedRange struct {
	Lower  uint32
	Length uint32
}

func (cr ClosedRange) Upper() uint32 {
	return cr.Lower + cr.Length - 1
}

func (cr ClosedRange) Compare(other ClosedRange) Order {
	upper := cr.Upper()
	otherUpper := other.Upper()

	if upper < other.Lower {
		if upper+1 == other.Lower {
			return Prependable
		}

		return Less
	}

	if otherUpper < cr.Lower {
		if otherUpper+1 == cr.Lower {
			return Appendable
		}

		return Greater
	}

	if cr.Lower == other.Lower {
		if upper == otherUpper {
			return Equal
		}

		if upper < otherUpper {
			return IncludedLeft
		}

		return IncludingLeft
	}

	if cr.Lower < other.Lower {
		if otherUpper == upper {
			return IncludingRight
		}

		if otherUpper < upper {
			return IncludingMiddle
		}

		return RightOverlap
	}

	if upper == otherUpper {
		return IncludedRight
	}

	if upper < otherUpper {
		return IncludedMiddle
	}

	return LeftOverlap
}

func (cr ClosedRange) Intersection(other ClosedRange) (ClosedRange, error) {
	if other.Upper() < cr.Lower || cr.Upper() < other.Lower {
		return ClosedRange{}, NoIntersection
	}

	lower := max(cr.Lower, other.Lower)
	upper := min(cr.Upper(), other.Upper())

	return ClosedRange{Lower: lower, Length: upper - lower + 1}, nil
}
