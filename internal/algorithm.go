package internal

import (
	"errors"
	"fmt"
	"net"
	"os"
)

var (
	OverMaxQueries = errors.New("Went over max queries")
	ADDRESS        = net.IPv4(0xc6, 0x29, 0x00, 0x04)
)

// Implements RFC1034 5.3.3
func GetAddress(domain string) (net.IP, error) {

	servers := newSlist(domain, []net.IP{ADDRESS}, []string{""})

	mess := ExampleMessage(1, domain)

	maxMessages := 15
	for i := 0; i < maxMessages; i++ {
		server, index := servers.getBestServer()
		fmt.Printf("Querying %s for %s\n", server, domain)
		response, err := QueryDomain(server, mess)
		mess.Header.id++
		if errors.Is(err, os.ErrDeadlineExceeded) {
			servers.remove(index)
			continue
		}
		if err != nil {
			return net.IP{}, err
		}
		// fmt.Printf("Received Message:\n")
		// fmt.Printf("Header %+v\n\n", message.Header)
		// fmt.Printf("Question(%d)\n%+v\n", len(message.Question), message.Question)
		// fmt.Printf("Answer(%d)\n%s\n", len(message.Answer), message.Answer)
		// fmt.Printf("Authority(%d)\n%s\n", len(message.Authority), message.Authority)
		// fmt.Printf("Additional(%d)\n%s\n", len(message.Additional), message.Additional)
		if response.Error() {
			return net.IP{}, fmt.Errorf("Error in rcode of response %d", response.Header.flags&0xf)
		}
		if response.Header.ancount != 0 {
			ip := response.Answer[0].rdata
			return ip, nil
		}

		if response.Header.nscount != 0 && response.Header.arcount != 0 {
			ips, domains, err := getServerAndDomain(response)
			if err != nil {
				return net.IP{}, fmt.Errorf("Could not get bytes from message")
			}
			servers = newSlist(domain, ips, domains)
			continue
		}
		if response.Header.nscount != 0 && response.Header.arcount == 0 {
			name := DecompressSingleDomain(response.Authority[0].rdata)
			s, err := GetAddress(name)
			if err != nil {
				return net.IP{}, err
			}
			servers = newSlist(domain, []net.IP{s}, []string{name})
			continue
		}
		break
	}
	s, _ := servers.getBestServer()
	return s, OverMaxQueries
}

func getServerAndDomain(m Message) ([]net.IP, []string, error) {
	b, err := m.Bytes()
	if err != nil {
		return []net.IP{}, []string{}, err
	}
	minlen := min(int(m.Header.nscount), int(m.Header.arcount))
	ips := make([]net.IP, minlen)
	domains := make([]string, minlen)
	for i := 0; i < minlen; i++ {
		ips[i] = m.Additional[i].rdata
		domains[i] = getDomainFromBytes(b, m.Authority[i].name)
	}

	return ips, domains, nil
}
