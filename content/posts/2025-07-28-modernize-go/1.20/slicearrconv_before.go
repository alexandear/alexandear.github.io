package main

import (
	"fmt"
	"net/netip"
)

func main() {
	// << snippet begin >>
	slice := []byte{192, 168, 0, 1}
	addr := *(*[4]byte)(slice)

	fmt.Println(netip.AddrFrom4(addr))
	// << snippet end >>
}
