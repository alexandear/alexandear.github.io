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

			for caseName := range findCaseNames(t) {
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

						executeScript(t, script, tmp)

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

func TestValidateCompilation(t *testing.T) {
	for _, ver := range []string{
		"1.22",
	} {
		t.Run(ver, func(t *testing.T) {
			t.Chdir(ver)

			const script = "compile.sh"

			for caseName := range findCaseNames(t) {
				t.Run("before", func(t *testing.T) {
					beforeFilename := beforeFilename(caseName)
					tmpBefore := t.TempDir()
					for _, src := range []string{beforeFilename, "go.mod", script} {
						dst := filepath.Join(tmpBefore, src)
						copyFile(t, src, dst)
					}
					executeScript(t, script, tmpBefore)
				})

				t.Run("after", func(t *testing.T) {
					afterFilename := afterFilename(caseName)
					tmpAfter := t.TempDir()
					for _, src := range []string{afterFilename, "go.mod", script} {
						dst := filepath.Join(tmpAfter, src)
						copyFile(t, src, dst)
					}
					executeScript(t, script, tmpAfter)
				})
			}
		})
	}
}

func findCaseNames(t *testing.T) map[string]struct{} {
	scripts, err := filepath.Glob("*.sh")
	if err != nil {
		t.Fatal(err)
	}
	caseNames := map[string]struct{}{}
	for _, script := range scripts {
		if script == "compile.sh" {
			continue
		}
		name, _, ok := strings.Cut(script, "_")
		if !ok {
			t.Fatalf("Script name %q should contain at least one _", script)
		}
		caseNames[name] = struct{}{}
		t.Logf("Found case: %q", name)
	}
	return caseNames
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

func executeScript(t *testing.T, script, workDir string) {
	t.Logf("Executing the script %q in %q", script, workDir)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	t.Cleanup(func() { cancel() })

	t.Chdir(workDir)
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "sh", script)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run command %q: %v\nStdout: %s\nStderr: %s", cmd, err, stdout.String(), stderr.String())
	}
}
