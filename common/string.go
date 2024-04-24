package common

import "unicode/utf16"

func UTF16Slice(s string, start, end uint32) string {
	encoded := utf16.Encode([]rune(s))
	return string(utf16.Decode(encoded[start:end]))
}

func UTF16Length(s string) uint32 {
	return uint32(len(utf16.Encode([]rune(s))))
}
