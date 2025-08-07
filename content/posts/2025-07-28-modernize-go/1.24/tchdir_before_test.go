package main_test

import (
	"os"
	"testing"
)

// << snippet begin >>
func TestSomeFunc(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	// Test logic for SomeFunc here

	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Error(err)
		}
	})
}

// << snippet end >>
