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
