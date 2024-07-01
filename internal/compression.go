package internal

import (
	"strings"
)

func CompressSingleDomain(domain string) []byte {
	labels := strings.Split(domain, ".")
	compressed := make([]byte, len(domain)+2)

	index := 0
	for _, s := range labels {
		compressed[index] = byte(len(s))
		index++
		for _, c := range s {
			compressed[index] = byte(c)
			index++
		}
	}
	compressed[len(domain)+1] = 0

	return compressed
}

func DecompressSingleDomain(compressed []byte) string {
	b := strings.Builder{}

	index := 0
	maxlen := len(compressed) - 1
	for {
		if index >= maxlen {
			break
		}
		labelLen := int(compressed[index])
		index++
		for i := index; i < index+labelLen; i++ {
			b.WriteByte(compressed[i])
		}
		index += labelLen
		if index+1 < maxlen {
			b.WriteByte('.')
		}
	}

	return b.String()
}
