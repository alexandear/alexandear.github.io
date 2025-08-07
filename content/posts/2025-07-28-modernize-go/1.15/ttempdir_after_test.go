package main_test

import (
	"testing"
)

// << snippet begin >>
func TestSomeFunc(t *testing.T) {
	tmp := t.TempDir()

	t.Log("Test logic for SomeFunc that uses temporary directory:", tmp)
}

// << snippet end >>
