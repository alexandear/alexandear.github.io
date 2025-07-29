package main_test

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRun(t *testing.T) {
	for _, ver := range []string{
		"1.22",
	} {
		t.Run(ver, func(t *testing.T) {
			t.Chdir(ver)
			scripts, err := filepath.Glob("*.sh")
			if err != nil {
				t.Fatal(err)
			}

			for _, script := range scripts {
				t.Run(script, func(t *testing.T) {
					filename := strings.TrimSuffix(script, ".sh")
					goFilename := filename + ".go"
					goldenFilename := filename + "_golden.go"

					tmp := t.TempDir()
					for _, src := range []string{script, goFilename, goldenFilename} {
						dst := filepath.Join(tmp, src)
						copyFile(t, src, dst)
					}

					t.Logf("Executing the script: %q", script)

					ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
					t.Cleanup(func() { cancel() })

					t.Chdir(tmp)
					cmd := exec.CommandContext(ctx, "sh", script)
					if err := cmd.Run(); err != nil {
						t.Fatalf("Failed to run command %q in file %q: %v", cmd, script, err)
					}

					contents, err := os.ReadFile(filename + ".go")
					if err != nil {
						t.Fatal(err)
					}
					contentsGolden, err := os.ReadFile(filename + "_golden.go")
					if err != nil {
						t.Fatal(err)
					}

					if diff := cmp.Diff(contents, contentsGolden); diff != "" {
						t.Errorf("Diff between file %q and golden: %v", filename, diff)
					}
				})
			}
		})
	}
}

func copyFile(t *testing.T, src, dst string) {
	sourceFile, err := os.Open(src)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := sourceFile.Close()
		if err != nil {
			t.Error(err)
		}
	})

	destFile, err := os.Create(dst)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := destFile.Close()
		if err != nil {
			t.Error(err)
		}
	})

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		t.Fatal(err)
	}
}
