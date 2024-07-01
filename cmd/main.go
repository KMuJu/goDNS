package main

import (
	"fmt"

	"github.com/kmuju/goDNS/internal"
)

func main() {
	err := internal.QueryDomain()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
