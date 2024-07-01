package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointerValue(t *testing.T) {
	cases := []struct {
		name     string
		pointer  [2]byte
		expected int
	}{
		{
			name:     "Test 1",
			pointer:  [2]byte{0xc0, 0x2c},
			expected: 44,
		},
		{
			name:     "Test 2",
			pointer:  [2]byte{0xc0, 0x17},
			expected: 23,
		},
		{
			name:     "Not A pointer",
			pointer:  [2]byte{0x00, 0x2c},
			expected: -1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ouput := pointerValue(tc.pointer)
			assert.Equal(t, tc.expected, ouput)
		})
	}
}

func TestPrefixToNull(t *testing.T) {
	cases := []struct {
		name          string
		input         []byte
		expectedbytes []byte
		expectedlen   int
	}{
		{
			name:          "Simple test",
			input:         []byte{0x22, 0x32, 0x00, 0x43, 0x1},
			expectedbytes: []byte{0x22, 0x32, 0x00},
			expectedlen:   3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, l := getPrefixToNull(tc.input)
			assert.Equal(t, tc.expectedbytes, b)
			assert.Equal(t, tc.expectedlen, l)
		})
	}
}

func TestIpFromBytes(t *testing.T) {
	input := [4]byte{0xc0, 0x21, 0x0e, 0x1e}
	output := ipFromBytes(input)
	expected := "192.33.14.30"
	assert.Equal(t, expected, output)
}
