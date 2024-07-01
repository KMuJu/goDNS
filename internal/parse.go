package internal

import (
	"encoding/binary"
	"fmt"
)

var ()

func Parse(response []byte) (Message, error) {
	var qd Question
	var an ResourceRecord
	var ns ResourceRecord
	var ar ResourceRecord
	currLen := 12
	if len(response) < 12 {
		return Message{}, fmt.Errorf("Response must be more than 12 bytes (for header), but was %d", len(response))
	}
	h, err := parseHeader(response[:12])
	if err != nil {
		return Message{}, err
	}

	fmt.Printf("Header: %+v\n", h)

	if h.qdcount != 0 {
		l := 0
		qd, l, err = parseQuestion(response[currLen:])
		if err != nil {
			return Message{}, err
		}
		currLen += l
	}

	if h.ancount != 0 {
		l := 0
		an, l, err = parseRR(response[currLen:])
		if err != nil {
			return Message{}, err
		}
		currLen += l
	}

	if h.nscount == 1 {
		l := 0
		ns, l, err = parseRR(response[currLen:])
		if err != nil {
			return Message{}, err
		}
		currLen += l
	}

	if h.arcount == 1 {
		l := 0
		ar, l, err = parseRR(response[12+currLen:])
		if err != nil {
			return Message{}, err
		}
		currLen += l
	}

	return Message{
		Header:     h,
		Question:   qd,
		Answer:     an,
		Authority:  ns,
		Additional: ar,
	}, nil
}

func parseHeader(h []byte) (Header, error) {
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

func parseRR(a []byte) (ResourceRecord, int, error) {
	name, l := getPointerOrLabels(a)
	if len(a) < l+10 {
		return ResourceRecord{}, 0, fmt.Errorf("Length of RR is too short")
	}
	t := binary.BigEndian.Uint16(a[l : l+2])
	class := binary.BigEndian.Uint16(a[l+2 : l+4])
	ttl := binary.BigEndian.Uint32(a[l+4 : l+8])
	rdlength := binary.BigEndian.Uint16(a[l+8 : l+10])
	rdata := a[l+10 : l+10+int(rdlength)]
	return ResourceRecord{
		notempty: true,
		name:     name,
		t:        t,
		class:    class,
		ttl:      ttl,
		rdlength: rdlength,
		rdata:    rdata,
	}, l + 10 + int(rdlength), nil
}

func parseQuestion(q []byte) (Question, int, error) {
	qname, l := getPointerOrLabels(q)
	if len(q) < l+4 {
		return Question{}, 0, fmt.Errorf("Length of RR is too short")
	}
	qtype := binary.BigEndian.Uint16(q[l : l+2])
	qclass := binary.BigEndian.Uint16(q[l+2 : l+4])
	return Question{
		empty:  false,
		qname:  qname,
		qtype:  qtype,
		qclass: qclass,
	}, l + 4, nil
}
