package main

import (
	"bytes"
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

// TestFix checks that the fix commands work correctly by running them and comparing before/after files.
func TestFix(t *testing.T) {
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
					filename := strings.TrimSuffix(script, "_command.sh")
					beforeFilename := beforeFilename(filename)
					afterContents, err := os.ReadFile(afterFilename(filename))
					if err != nil {
						t.Fatal(err)
					}

					tmp := t.TempDir()
					for _, src := range []string{script, beforeFilename, "go.mod"} {
						dst := filepath.Join(tmp, src)
						copyFile(t, src, dst)
					}

					t.Logf("Executing the script %q in %q", script, tmp)

					ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
					t.Cleanup(func() { cancel() })

					t.Chdir(tmp)
					var stdout, stderr bytes.Buffer
					cmd := exec.CommandContext(ctx, "sh", script)
					cmd.Stdout = &stdout
					cmd.Stderr = &stderr
					if err := cmd.Run(); err != nil {
						t.Fatalf("Failed to run command %q: %v\nStdout: %s\nStderr: %s", cmd, err, stdout.String(), stderr.String())
					}

					beforeContents, err := os.ReadFile(filepath.Join(tmp, beforeFilename))
					if err != nil {
						t.Fatal(err)
					}

					if diff := cmp.Diff(beforeContents, afterContents); diff != "" {
						t.Errorf("Diff between before and after for file %q: %v", filename, diff)
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

	t.Logf("Copying %q to %q", src, dst)
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		t.Fatal(err)
	}
}
