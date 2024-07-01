package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleDomainCompression(t *testing.T) {
	domain := "dns.google.com"
	expected := []byte{0x03, 0x64, 0x6e, 0x73, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00}
	output := CompressSingleDomain(domain)
	assert.Equal(t, expected, output)
}

func TestDecompress(t *testing.T) {

	compressed := []byte{0x03, 0x64, 0x6e, 0x73, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00}
	expected := "dns.google.com"
	ouput := DecompressSingleDomain(compressed)
	assert.Equal(t, expected, ouput)
}
