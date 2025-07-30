package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// Checks that the fix commands work correctly by running them and comparing before/after files.
func TestValidateFixCommands(t *testing.T) {
	for _, ver := range []string{
		"1.22",
	} {
		t.Run(ver, func(t *testing.T) {
			t.Chdir(ver)
			scripts, err := filepath.Glob("*.sh")
			if err != nil {
				t.Fatal(err)
			}
			caseNames := map[string]struct{}{}
			for _, script := range scripts {
				name, _, ok := strings.Cut(script, "_")
				if !ok {
					t.Fatalf("Script name %q should contain at least one _", script)
				}
				caseNames[name] = struct{}{}
				t.Logf("Found case: %q", name)
			}

			for caseName := range caseNames {
				caseScripts, err := filepath.Glob(caseName + "*.sh")
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("Found case scripts: %v", caseScripts)

				for _, script := range caseScripts {
					t.Run(script, func(t *testing.T) {
						beforeFilename := beforeFilename(caseName)
						afterContents, err := os.ReadFile(afterFilename(caseName))
						if err != nil {
							t.Fatal(err)
						}

						tmp := t.TempDir()
						for _, src := range slices.Concat([]string{beforeFilename, "go.mod"}, caseScripts) {
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
							t.Errorf("Diff between before and after for case %q: %v", caseName, diff)
							t.Logf("Contents of %v:\n%s", beforeFilename, beforeContents)
						}
					})
				}
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
