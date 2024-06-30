package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Returns a header with qdcount = 1 and the rest of the counts = 0
func NewHeader(id uint16, qr, aa, tc, rd, ra bool, rcode byte) Header {
	f := uint16(rcode & 15) // rcode ^ 0b1111
	if qr {
		f |= 1 << 15
	}
	if aa {
		f |= 1 << 10
	}
	if tc {
		f |= 1 << 9
	}
	if rd {
		f |= 1 << 8
	}
	flag := uint16(f)
	return Header{
		id:      id,
		flags:   flag,
		qdcount: 1,
		ancount: 0,
		arcount: 0,
		nscount: 0,
	}
}

func WriteRR(buf io.Writer, rr ResourceRecord) error {
	if !rr.notempty {
		return nil
	}
	if err := binary.Write(buf, binary.BigEndian, []byte(rr.name)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, rr.t); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, rr.class); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, rr.ttl); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, rr.rdlength); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, rr.rdata); err != nil {
		return err
	}
	return nil
}

func (m Message) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	//
	// Header
	//
	if err := binary.Write(buf, binary.BigEndian, m.Header.id); err != nil {
		return []byte{}, errors.New("Error in writing Id")
	}
	if err := binary.Write(buf, binary.BigEndian, m.Header.flags); err != nil {
		return []byte{}, errors.New("Error in writing flags")
	}
	if err := binary.Write(buf, binary.BigEndian, m.Header.qdcount); err != nil {
		return []byte{}, errors.New("Error in writing qdcount")
	}
	if err := binary.Write(buf, binary.BigEndian, m.Header.ancount); err != nil {
		return []byte{}, errors.New("Error in writing ancount")
	}
	if err := binary.Write(buf, binary.BigEndian, m.Header.nscount); err != nil {
		return []byte{}, errors.New("Error in writing nscount")
	}
	if err := binary.Write(buf, binary.BigEndian, m.Header.arcount); err != nil {
		return []byte{}, errors.New("Error in writing arcount")
	}

	//
	// Question
	//
	if !m.Question.empty {
		// fmt.Printf("Question name %s\n", m.Question.qname)
		// fmt.Printf("Buf %x\n", buf.Bytes())
		// fmt.Printf("Actual string %x\n", []byte(m.Question.qname))
		if err := binary.Write(buf, binary.BigEndian, []byte(m.Question.qname)); err != nil {
			return []byte{}, errors.New("Error in writing qname")
		}
		if err := binary.Write(buf, binary.BigEndian, m.Question.qtype); err != nil {
			return []byte{}, errors.New("Error in writing qtype")
		}
		if err := binary.Write(buf, binary.BigEndian, m.Question.qclass); err != nil {
			return []byte{}, errors.New("Error in writing qclass")
		}
	}

	if err := WriteRR(buf, m.Answer); err != nil {
		return []byte{}, errors.New("Error in writing Answer")
	}
	if err := WriteRR(buf, m.Authority); err != nil {
		return []byte{}, errors.New("Error in writing Authority")
	}
	if err := WriteRR(buf, m.Additional); err != nil {
		return []byte{}, errors.New("Error in writing Additional")
	}

	return buf.Bytes(), nil
}

func (m Message) String() string {
	b, err := m.Bytes()
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%x", b)
}

func ExampleMessage(id uint16) Message {
	name := CompressSingleDomain("dns.google.com")
	return Message{
		Header: NewHeader(id, false, false, false, true, false, 0),
		Question: Question{
			empty:  false,
			qname:  name,
			qtype:  1,
			qclass: 1,
		},
	}
}
