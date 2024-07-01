package internal

import (
	"fmt"
	"net"
)

const (
	TYPE    = "udp"
	ADDRESS = "198.41.0.4:53"
	PORT    = 53
)

func QueryDomain() error {

	conn, err := net.Dial(TYPE, ADDRESS)
	if err != nil {
		fmt.Printf("Could not connect to address\n%s\n", err.Error())
		return err
	}
	defer conn.Close()

	ex := ExampleMessage(22)
	b, err := ex.Bytes()
	if err != nil {
		return err
	}
	fmt.Printf("Flags: %d\n", ex.Header.flags)
	_, err = conn.Write(b)
	if err != nil {
		return err
	}

	response := make([]byte, 1024) // Adjust buffer size as needed
	n, err := conn.Read(response)
	if err != nil {
		return err
	}

	fmt.Printf("Received %d bytes response:\n%x\n", n, response[:n])
	m, err := Parse(response[:n])
	if err != nil {
		return err
	}
	b, err = m.Bytes()
	if err != nil {
		return err
	}
	fmt.Printf("Header\n %+v\n", m.Header)
	fmt.Printf("Question(%d)\n %+v\n", len(m.Question), m.Question)
	fmt.Printf("Answer(%d)\n %s\n", len(m.Answer), m.Answer)
	fmt.Printf("Authority(%d)\n %s\n", len(m.Authority), m.Authority)
	fmt.Printf("Additional(%d)\n%s\n", len(m.Additional), m.Additional)

	return err
}
