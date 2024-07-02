package internal

import (
	"errors"
	"fmt"
	"os"
)

var OverMaxQueries = errors.New("Went over max queries")

// Implements RFC1034 5.3.3
func GetAddress(domain string) (string, error) {

	cache, b := isCached(domain)
	if b {
		return cache, nil
	}
	servers := []string{ADDRESS}

	mess := ExampleMessage(22, domain)

	maxMessages := 15
	serverIndex := 0
	for i := 0; i < maxMessages; i++ {
		server := servers[serverIndex]
		// fmt.Printf("Sending to %s\n", server)
		fmt.Printf("Querying %s for %s\n", server, domain)
		message, err := QueryDomain(server, mess)
		if errors.Is(err, os.ErrDeadlineExceeded) {
			// fmt.Printf("Deadline exceeded\n")
			serverIndex = (serverIndex + 1) % len(servers)
			continue
		}
		if err != nil {
			return "", err
		}
		// fmt.Printf("Received Message:\n")
		// fmt.Printf("Header %+v\n\n", message.Header)
		// fmt.Printf("Question(%d)\n%+v\n", len(message.Question), message.Question)
		// fmt.Printf("Answer(%d)\n%s\n", len(message.Answer), message.Answer)
		// fmt.Printf("Authority(%d)\n%s\n", len(message.Authority), message.Authority)
		// fmt.Printf("Additional(%d)\n%s\n", len(message.Additional), message.Additional)
		if message.Error() {
			return "", fmt.Errorf("Error in rcode of response %d", message.Header.flags&0xf)
		}
		if message.Header.ancount != 0 {
			return message.Answer[0].getIp(), nil
		}
		if message.Header.nscount != 0 && message.Header.arcount != 0 {
			servers = getServers(message)
			serverIndex = 0
			continue
		}
		if message.Header.nscount != 0 && message.Header.arcount == 0 {
			name := DecompressSingleDomain(message.Authority[0].rdata)
			// fmt.Printf("Starting new query for %s\n------\n\n", name)
			s, err := GetAddress(name)
			if err != nil {
				return "IDK MAN", err
			}
			// fmt.Printf("Result: %s\n-------\n", s)
			servers = []string{s}
			continue
		}
		break
	}
	return servers[0], OverMaxQueries
}

func isCached(_ string) (string, bool) {
	return "", false
}

func getServers(m Message) []string {
	output := make([]string, m.Header.arcount)
	for i, additional := range m.Additional {
		output[i] = additional.getIp()
	}
	return output
}
