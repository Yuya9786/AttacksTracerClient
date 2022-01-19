package main

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type InetAddress struct {
	octet1 byte
	octet2 byte
	octet3 byte
	octet4 byte
}

func inetAddress(address string) (*InetAddress, error) {
	octets := strings.Split(address, ".")
	if len(octets) != 4 {
		return nil, errors.New("invalid format")
	}

	var convertedOctets = make([]int, 4)
	for i, v := range octets {
		octet, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.Wrap(err, "failed to atoi")
		}
		convertedOctets[i] = octet
	}

	addr := &InetAddress{
		octet1: byte(convertedOctets[0]),
		octet2: byte(convertedOctets[1]),
		octet3: byte(convertedOctets[2]),
		octet4: byte(convertedOctets[3]),
	}

	return addr, nil
}

func (a *InetAddress) toIntAddress() int {
	return (int(a.octet1) & 0xff) << 24 | (int(a.octet2) & 0xff) << 16 | 
		(int(a.octet3) & 0xff) << 8 | (int(a.octet4) & 0xff)
}

func (a *InetAddress) isSameNetwork(b *InetAddress, mask int) bool {
	shift := 32 - mask
	return ((a.toIntAddress() >> shift) << shift) == ((b.toIntAddress() >> shift) << shift)
}