package network

import (
	"net"
	"testing"
)

func TestAllocation(t *testing.T) {
	ip, ipnet, _ := net.ParseCIDR("192.168.0.0/24")

	aip, _ := ipAllocator.Allocate(ipnet)

	t.Logf("first ip: %v, allocate IP: %s", ip, aip)
}

func TestRelease(t *testing.T) {
	ip, ipnet, _ := net.ParseCIDR("192.168.0.1/24")

	err := ipAllocator.Release(ipnet, &ip)

	t.Logf("first ip: %v, err: %v", ip, err)
}
