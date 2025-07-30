package main

import "fmt"

func main() {
	done := make(chan bool)

	// << snippet begin >>
	values := []string{"a", "b", "c"}
	for _, v := range values {
		v := v
		go func() {
			fmt.Println(v)
			done <- true
		}()
	}
	// << snippet end >>

	for range values {
		<-done
	}
}
