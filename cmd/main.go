package main

import (
	"flag"
	"fmt"

	"github.com/kmuju/goDNS/internal"
)

func main() {
	flag.Parse()
	Narg := flag.NArg()
	if Narg == 0 {
		fmt.Printf("need more arguments")
		return
	}
	domains := flag.Args()
	for i, domain := range domains {
		ip, err := internal.GetAddress(domain)
		if err == internal.OverMaxQueries {
			fmt.Printf("%s Overstepped maxqueries, last ip %s\n", domain, ip)
			continue
		}
		if err != nil {
			fmt.Printf("Error quering %s: %s", domain, err)
			continue
		}
		fmt.Printf("%s\n", ip)
		if i+1 < Narg {
			fmt.Print("\n")
		}
	}
}
