package content

type Content[T any] interface {
	*INSContent | *DELContent | *FMTContent | *MODContent

	Length() uint32
	Concat(other T) T
	Slice(start uint32, end uint32) T
	Equals(other T) bool
}
