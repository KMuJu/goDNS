package internal

import (
	"math"
	"net"
	"strings"
)

const (
	MAXSIZE = 10
)

type sList struct {
	servers [MAXSIZE]server
	size    int // points to the next index to insert
	target  string
}

type server struct {
	score int
	ip    net.IP
}

func newSlist(target string, ips []net.IP, domains []string) sList {
	sl := sList{target: target}
	sl.add(ips, domains)
	return sl
}

// gets the server with the highest score and removes it
func (sl *sList) getBestServer() net.IP {
	maxscore, index := math.MinInt, -1

	for i, s := range sl.servers {
		if s.score > maxscore {
			maxscore = s.score
			index = i
		}
	}
	sl.remove(index)

	return sl.servers[index].ip
}

// moves the element at the index to the back and shrinks size
func (sl *sList) remove(index int) bool {
	if sl.size == 0 {
		return false
	}
	sl.swap(index, sl.size-1)
	sl.size--
	return true
}

func (sl *sList) swap(aindex, bindex int) {
	a := sl.servers[aindex]
	sl.servers[aindex] = sl.servers[bindex]
	sl.servers[bindex] = a
}

func (sl *sList) addSingle(ip net.IP, domain string) bool {
	if sl.size >= MAXSIZE {
		return false
	}
	sl.servers[sl.size] = server{
		ip:    ip,
		score: score(sl.target, domain),
	}
	sl.size++
	return true
}

func (sl *sList) add(ips []net.IP, domains []string) {
	minlen := min(len(ips), len(domains))
	for i := 0; i < minlen; i++ {
		if !sl.addSingle(ips[i], domains[i]) {
			break
		}
	}
}

func score(target, input string) int {
	tlabels := strings.Split(target, ".")
	ilabels := strings.Split(input, ".")
	tindex := len(tlabels) - 1
	inindex := len(ilabels) - 1
	score := 0
	for tindex >= 0 && inindex >= 0 {
		tlabel := tlabels[tindex]
		ilabel := ilabels[inindex]
		if tlabel != ilabel {
			break
		}
		tindex--
		inindex--
		score++
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
