package main

import "fmt"

// << snippet begin >>
func main() {
	a, b := 4, -1

	h := min(a, b)
	m := min(a, b)

	fmt.Println(h, m)
}

// << snippet end >>
