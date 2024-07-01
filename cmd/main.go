package main

import (
	"flag"
	"fmt"

	"github.com/kmuju/goDNS/internal"
)

func main() {
	_ = flag.String("domain", "nrk.no", "")
	flag.Parse()
	fmt.Printf("%+v\n", flag.Args())
	if flag.NArg() == 0 {
		fmt.Printf("need more arguments")
		return
	}
	domain := flag.Arg(0)
	ip, err := internal.GetAddress(domain)
	if err == internal.OverMaxQueries {
		fmt.Printf("Overstepped maxqueries, last ip %s\n", ip)
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Ip address: %s\n", ip)
}
