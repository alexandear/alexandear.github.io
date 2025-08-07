package main_test

import (
	"testing"
)

// << snippet begin >>
func TestSomeFunc(t *testing.T) {
	t.Setenv("SOME_ENV", "new_value")

	t.Log("Test logic for SomeFunc depends on the env variable SOME_ENV")
}

// << snippet end >>
