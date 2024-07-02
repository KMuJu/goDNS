package internal

import (
	"strings"
)

func getLength(domains []string) (length int) {
	for _, d := range domains {
		length += len(d) + 2
	}
	return length
}

// TODO: this is wrong, the compressed strings does not end in the same byte slice
// It will be used for different questions and should maybe return a [][]byte???
func CompressDomains(domains ...string) []byte {
	compressed := make([]byte, getLength(domains))
	index := 0
	for _, domain := range domains {
		labels := strings.Split(domain, ".")

		for _, s := range labels {
			compressed[index] = byte(len(s))
			index++
			for _, c := range s {
				compressed[index] = byte(c)
				index++
			}
		}
		compressed[index] = 0
		index++
	}

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
