package internal

import (
	"fmt"
	"net"
)

const (
	TYPE    = "udp"
	ADDRESS = "8.8.8.8:53"
	PORT    = 53
)

func QueryDomain() error {

	conn, err := net.Dial(TYPE, ADDRESS)
	if err != nil {
		fmt.Printf("Could not connect to address\n%s\n", err.Error())
		return err
	}
	defer conn.Close()

	b, err := ExampleMessage(22).Bytes()
	if err != nil {
		return err
	}
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

	return err
}
