package common

const UInt32Min = uint32(0)
const UInt32Mid = (uint32(^uint32(0)) - uint32(0)) / 2
const UInt32Max = uint32(^uint32(0))

func CompareNumber[T ~uint32 | ~int](a, b T) Order {
	if a == b {
		return Equal
	}
	if a < b {
		return Less
	}
	return Greater
}
