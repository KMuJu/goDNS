package internal

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSlist(t *testing.T) {
	cases := []struct {
		name     string
		target   string
		ips      []net.IP
		domains  []string
		expected sList
	}{
		{
			name:   "Simple init",
			target: "domain.com",
			ips: []net.IP{
				net.IPv4(0, 0, 0, 0),
				net.IPv4(0, 0, 0, 1),
			},
			domains: []string{
				".com",
				"domain.com",
			},
			expected: sList{
				servers: [MAXSIZE]server{
					{
						score: 1,
						ip:    net.IPv4(0, 0, 0, 0),
					},
					{
						score: 2,
						ip:    net.IPv4(0, 0, 0, 1),
					},
				},
				size:   2,
				target: "domain.com",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output := newSlist(tc.target, tc.ips, tc.domains)
			assert.Equal(t, tc.expected, output)
		})
	}
}

func TestGetBestIP(t *testing.T) {
	cases := []struct {
		name       string
		slist      sList
		expectedip net.IP
		index      int
	}{
		{
			name: "two ips",
			slist: sList{
				servers: [MAXSIZE]server{
					{
						score: 1,
						ip:    net.IPv4(0, 0, 0, 0),
					},
					{
						score: 2,
						ip:    net.IPv4(0, 0, 0, 1),
					},
				},
				size:   2,
				target: "domain.com",
			},
			expectedip: net.IPv4(0, 0, 0, 1),
			index:      1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ip, index := tc.slist.getBestServer()
			assert.Equal(t, tc.expectedip, ip)
			assert.Equal(t, tc.index, index)
		})
	}
}

func TestScore(t *testing.T) {
	cases := []struct {
		name          string
		target        string
		input         string
		expectedscore int
	}{
		{
			name:          "Single label 1",
			target:        "a",
			input:         "a",
			expectedscore: 1,
		},
		{
			name:          "Single label 0",
			target:        "a",
			input:         "b",
			expectedscore: 0,
		},
		{
			name:          "two equal labels",
			target:        "a.b",
			input:         "a.b",
			expectedscore: 2,
		},
		{
			name:          "two labels one equal",
			target:        "a.b",
			input:         "b.b",
			expectedscore: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			score := score(tc.target, tc.input)
			assert.Equal(t, tc.expectedscore, score)
		})
	}

}
