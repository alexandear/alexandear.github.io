package main

import (
	"fmt"
	"math"
)

// << snippet begin >>
func main() {
	a, b := 4, -1

	h := min(a, b)
	m := int(math.Min(float64(a), float64(b)))

	fmt.Println(h, m)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// << snippet end >>
