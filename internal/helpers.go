package internal

import (
	"net"
)

func getPrefixToNull(input []byte) ([]byte, int) {
	output := make([]byte, len(input))
	for i, b := range input {
		if b == 0 {
			return output[:i+1], i + 1
		}
		output[i] = b
	}
	return output, len(input)
}

func getPointerOrLabels(input []byte) ([]byte, int) {
	if len(input) < 2 {
		return input, len(input)
	}
	if isPointer(input) {
		return input[:2], 2
	}
	return getPrefixToNull(input)
}

func pointerValue(p [2]byte) int {
	if !isPointer(p[:]) {
		return -1
	}
	// p[0] & 0b0011 1111
	// to remove the first two bits of the pointer
	return int(p[0]&63)<<8 | int(p[1])
}

func isPointer(b []byte) bool {
	return b[0]&192 != 0 // b[0] & 0b1100 0000
}

func getValueFromPointer(input []byte, p [2]byte) []byte {
	var output []byte
	val := pointerValue(p)
	for i := val; i < len(input); i++ {
		if input[i] == 0 {
			break
		}
		output = append(output, input[i])
	}
	return output
}

func ipFromBytes(input [4]byte) string {
	return net.IPv4(input[0], input[1], input[2], input[3]).String()
}
