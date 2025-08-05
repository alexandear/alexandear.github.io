package main

import (
	"fmt"
	"net/netip"
)

func main() {
	// << snippet begin >>
	ip := []byte{192, 168, 0, 1}
	addr := [4]byte(ip)

	fmt.Println(netip.AddrFrom4(addr))
	// << snippet end >>
}
