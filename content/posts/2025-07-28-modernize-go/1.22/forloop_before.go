package main

import "fmt"

func main() {
	done := make(chan bool)

	// << begin >>
	values := []string{"a", "b", "c"}
	for _, v := range values {
		v := v
		go func() {
			fmt.Println(v)
			done <- true
		}()
	}
	// << end >>

	for _ = range values {
		<-done
	}
}
