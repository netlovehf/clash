package fakeip

import (
	"errors"
	"net"
)

// Pool is a implementation about fake ip generator without storage
type Pool struct {
	max    uint32
	min    uint32
	offset uint32
}

// Get return a new fake ip
func (p *Pool) Get() net.IP {
	ip := uintToIP(p.min + p.offset)
	p.offset = (p.offset + 1) % (p.max - p.min)
	return ip
}

func ipToUint(ip net.IP) uint32 {
	v := uint32(ip[0]) << 24
	v += uint32(ip[1]) << 16
	v += uint32(ip[2]) << 8
	v += uint32(ip[3])
	return v
}

func uintToIP(v uint32) net.IP {
	return net.IPv4(byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// New return Pool instance
func New(ipnet *net.IPNet) (*Pool, error) {
	min := ipToUint(ipnet.IP) + 1

	ones, bits := ipnet.Mask.Size()
	total := 1<<uint(bits-ones) - 2

	if total <= 0 {
		return nil, errors.New("ipnet don't have valid ip")
	}

	max := min + uint32(total)
	return &Pool{
		min: min,
		max: max,
	}, nil
}
