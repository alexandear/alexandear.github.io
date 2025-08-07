package main_test

import (
	"io/ioutil"
	"os"
	"testing"
)

// << snippet begin >>
func TestSomeFunc(t *testing.T) {
	tmp, err := ioutil.TempDir("", "pattern")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Error(err)
		}
	}()

	t.Log("Test logic for SomeFunc that uses temporary directory:", tmp)
}

// << snippet end >>
