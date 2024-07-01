package internal

type Message struct {
	Header     Header
	Question   []Question
	Answer     []ResourceRecord
	Authority  []ResourceRecord
	Additional []ResourceRecord
}

type Header struct {
	id      uint16
	flags   uint16
	qdcount uint16
	ancount uint16
	nscount uint16
	arcount uint16
}

type Question struct {
	empty  bool
	qname  []byte
	qtype  uint16
	qclass uint16
}

type ResourceRecord struct {
	notempty bool
	name     []byte
	t        uint16
	class    uint16
	ttl      uint32
	rdlength uint16
	rdata    []byte
}
