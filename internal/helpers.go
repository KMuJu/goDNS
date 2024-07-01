package internal

import "fmt"

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
	fmt.Printf("Input: %x\n", input)
	if len(input) < 2 {
		return input, len(input)
	}
	if input[0]&192 != 0 {
		return input[:2], 2
	}
	return getPrefixToNull(input)
}
