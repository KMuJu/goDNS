package internal

import (
	"fmt"
	"net"
	"time"
)

const (
	TYPE    = "udp"
	ADDRESS = "198.41.0.4"
	PORT    = 53
)

func QueryDomain(address string, mess Message) (Message, error) {

	conn, err := net.Dial(TYPE, fmt.Sprintf("%s:53", address))
	if err != nil {
		fmt.Printf("Could not connect to address\n%s\n", err.Error())
		return Message{}, err
	}
	defer conn.Close()
	err = conn.SetDeadline(time.Now().Add(time.Second))

	b, err := mess.Bytes()
	if err != nil {
		return Message{}, err
	}
	_, err = conn.Write(b)
	if err != nil {
		return Message{}, err
	}

	response := make([]byte, 1024) // Adjust buffer size as needed
	n, err := conn.Read(response)
	if err != nil {
		return Message{}, err
	}

	// fmt.Printf("Received %d bytes response:\n%x\n", n, response[:n])
	m, err := Parse(response[:n])
	if err != nil {
		return Message{}, err
	}

	return m, err
}
