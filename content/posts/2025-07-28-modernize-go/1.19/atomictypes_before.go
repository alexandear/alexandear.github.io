package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// << snippet begin >>
type Service struct {
	running uint32 // 0 = false, 1 = true
}

func (s *Service) Start() {
	atomic.StoreUint32(&s.running, 1)
}

func (s *Service) Stop() error {
	if !atomic.CompareAndSwapUint32(&s.running, 1, 0) {
		return errors.New("service was not running")
	}
	return nil
}

func (s *Service) IsRunning() bool {
	return atomic.LoadUint32(&s.running) == 1
}

// << snippet end >>

func main() {
	var wg sync.WaitGroup
	service := &Service{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		service.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("Service running: %t\n", service.IsRunning())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := service.Stop(); err != nil {
			fmt.Println("Error in goroutine:", err)
		}
	}()

	wg.Wait()
	if err := service.Stop(); err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("Service stopped: %t\n", service.IsRunning())
}
