package main

import (
	"fmt"

	"github.com/kmuju/goDNS/internal"
)

func main() {
	mess := internal.ExampleMessage(22)

	fmt.Printf("%s\n", mess)
	fmt.Println("00160100000100000000000003646e7306676f6f676c6503636f6d0000010001")

	fmt.Printf("------\n")

	err := internal.QueryDomain()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
