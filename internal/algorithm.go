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

	servers := []net.IP{ADDRESS}

	mess := ExampleMessage(1, domain)

	maxMessages := 15
	serverIndex := 0
	for i := 0; i < maxMessages; i++ {
		server := servers[serverIndex]
		fmt.Printf("Querying %s for %s\n", server, domain)
		message, err := QueryDomain(server, mess)
		mess.Header.id++
		if errors.Is(err, os.ErrDeadlineExceeded) {
			serverIndex = (serverIndex + 1) % len(servers)
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
		if message.Error() {
			return net.IP{}, fmt.Errorf("Error in rcode of response %d", message.Header.flags&0xf)
		}
		if message.Header.ancount != 0 {
			ip := message.Answer[0].rdata
			return ip, nil
		}

		if message.Header.nscount != 0 && message.Header.arcount != 0 {
			servers = getServers(message)
			serverIndex = 0
			continue
		}
		if message.Header.nscount != 0 && message.Header.arcount == 0 {
			name := DecompressSingleDomain(message.Authority[0].rdata)
			s, err := GetAddress(name)
			if err != nil {
				return net.IP{}, err
			}
			servers = []net.IP{s}
			continue
		}
		break
	}
	return servers[0], OverMaxQueries
}

func getServers(m Message) []net.IP {
	output := make([]net.IP, m.Header.arcount)
	for i, additional := range m.Additional {
		output[i] = additional.rdata
	}
	return output
}
