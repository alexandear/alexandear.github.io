package main_test

import (
	"os"
	"testing"
)

// << snippet begin >>
func TestSomeFunc(t *testing.T) {
	orig := os.Getenv("SOME_ENV")
	if err := os.Setenv("SOME_ENV", "new_value"); err != nil {
		t.Fatal(err)
	}

	t.Log("Test logic for SomeFunc depends on the env variable SOME_ENV")

	t.Cleanup(func() {
		// set to original value after the test
		if err := os.Setenv("SOME_ENV", orig); err != nil {
			t.Error(err)
		}
	})
}

// << snippet end >>
