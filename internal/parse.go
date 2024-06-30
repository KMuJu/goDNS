package internal

import (
	"encoding/binary"
	"fmt"
)

var ()

func Parse(response []byte) Message {

	return Message{}
}

func ParseHeader(h []byte) (Header, error) {
	if len(h) != 12 {
		return Header{}, fmt.Errorf("Expected length of 12, got: %d", len(h))
	}
	id := binary.BigEndian.Uint16(h[:2])
	flag := binary.BigEndian.Uint16(h[2:4])
	qdcount := binary.BigEndian.Uint16(h[4:6])
	ancount := binary.BigEndian.Uint16(h[6:8])
	nscount := binary.BigEndian.Uint16(h[8:10])
	arcount := binary.BigEndian.Uint16(h[10:12])
	return Header{
		id:      id,
		flags:   flag,
		qdcount: qdcount,
		ancount: ancount,
		nscount: nscount,
		arcount: arcount,
	}, nil
}

func parseQuestion(_ []byte) (Question, error) {
	return Question{}, nil
}
