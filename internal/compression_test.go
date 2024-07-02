package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainCompression(t *testing.T) {
	cases := []struct {
		name     string
		domains  []string
		expected []byte
	}{
		{
			name:     "googleDNS",
			domains:  []string{"dns.google.com"},
			expected: []byte{0x03, 0x64, 0x6e, 0x73, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00},
		},
		{
			name:     "nrk.no",
			domains:  []string{"nrk.no"},
			expected: []byte{0x03, 0x6e, 0x72, 0x6b, 0x02, 0x6e, 0x6f, 0x00},
		},
		{
			name:     "multiple domains",
			domains:  []string{"vg.no", "nrk.no"},
			expected: []byte{0x2, 0x76, 0x67, 0x02, 0x6e, 0x6f, 0x00, 0x03, 0x6e, 0x72, 0x6b, 0x02, 0x6e, 0x6f, 0x00},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output := CompressDomains(tc.domains...)
			assert.Equal(t, tc.expected, output)
		})
	}
}

func TestDecompress(t *testing.T) {
	cases := []struct {
		name       string
		expected   string
		compressed []byte
	}{
		{
			name:       "googleDNS",
			compressed: []byte{0x03, 0x64, 0x6e, 0x73, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00},
			expected:   "dns.google.com",
		},
		{
			name:       "nrk.no",
			compressed: []byte{0x03, 0x6e, 0x72, 0x6b, 0x02, 0x6e, 0x6f, 0x00},
			expected:   "nrk.no",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ouput := DecompressSingleDomain(tc.compressed)
			assert.Equal(t, tc.expected, ouput)
		})
	}
}
