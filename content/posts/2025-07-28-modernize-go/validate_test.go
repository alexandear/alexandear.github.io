package main

import (
	"bytes"
	"cmp"
	"context"
	"go/version"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	gocmp "github.com/google/go-cmp/cmp"
)

const noFix = "# No auto fix: "

// Checks that the fix commands work correctly by running them and comparing before/after files.
func TestValidateFixCommands(t *testing.T) {
	for _, item := range gos {
		if len(item.Sections) == 0 {
			continue
		}

		t.Run(item.Version, func(t *testing.T) {
			t.Chdir(item.Version)

			for caseName := range findCaseNames(t) {
				caseScripts, err := filepath.Glob(caseName + "*.sh")
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("Found case scripts: %v", caseScripts)

				for _, script := range caseScripts {
					t.Run(script, func(t *testing.T) {
						beforeFilename, err := beforeFilename(caseName)
						if err != nil {
							t.Fatal(err)
						}
						afterFilename, err := afterFilename(caseName)
						if err != nil {
							t.Fatal(err)
						}
						afterContents, err := os.ReadFile(afterFilename)
						if err != nil {
							t.Fatal(err)
						}

						tmp := t.TempDir()
						for _, src := range slices.Concat([]string{beforeFilename, "go.mod"}, caseScripts) {
							dst := filepath.Join(tmp, src)
							copyFile(t, src, dst)
						}

						scriptContents, err := os.ReadFile(script)
						if err != nil {
							t.Fatal(err)
						}

						stdout := executeScript(t, script, 0, tmp)

						if _, after, ok := strings.Cut(string(scriptContents), noFix); ok {
							after = strings.Split(after, "\n")[0]
							t.Logf("Script %q is marked as not fixing, checking for output: %q", script, after)

							if !strings.Contains(stdout, after) {
								t.Errorf("Missed expected output in %q in script %q", stdout, script)
							}
							return
						}

						beforeContents, err := os.ReadFile(filepath.Join(tmp, beforeFilename))
						if err != nil {
							t.Fatal(err)
						}

						if diff := gocmp.Diff(beforeContents, afterContents); diff != "" {
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
	const compileTimeout = 10 * time.Second

	for _, item := range gos {
		if len(item.Sections) == 0 {
			continue
		}

		t.Run(item.Version, func(t *testing.T) {
			t.Chdir(item.Version)
			installCompiler(t, item.compilerVersion)
			for caseName := range findCaseNames(t) {
				const script = "compile.sh"

				t.Run("before", func(t *testing.T) {
					beforeFilename, err := beforeFilename(caseName)
					if err != nil {
						t.Fatal(err)
					}
					beforeDir := t.TempDir()
					for _, src := range []string{beforeFilename, "go.mod", script} {
						dst := filepath.Join(beforeDir, src)
						copyFile(t, src, dst)
					}
					executeScript(t, script, compileTimeout, beforeDir)
				})

				t.Run("after", func(t *testing.T) {
					afterFilename, err := afterFilename(caseName)
					if err != nil {
						t.Fatal(err)
					}
					afterDir := t.TempDir()
					for _, src := range []string{afterFilename, "go.mod", script} {
						dst := filepath.Join(afterDir, src)
						copyFile(t, src, dst)
					}
					executeScript(t, script, compileTimeout, afterDir)
				})
			}
		})
	}
}

func installCompiler(t *testing.T, compilerVersion string) {
	langVersion := version.Lang("go" + compilerVersion)
	installCompilerFn := installCompiler116
	if version.Compare("go1.15", "go"+langVersion) <= 0 {
		installCompilerFn = installCompiler115
	}
	installCompilerFn(t, compilerVersion)
}

// installCompiler116 installs go1.16 or above by running the following commands:
//
//	go install golang.org/dl/go1.16@latest
//	go1.16 download
func installCompiler116(t *testing.T, compilerVersion string) {
	t.Helper()
	if compilerVersion == "" {
		t.Fatal("empty compilerVersion")
	}

	_, err := exec.LookPath("go" + compilerVersion)
	if err != nil {
		execCommand(t, 20*time.Second, nil, "go", "install", "golang.org/dl/go"+compilerVersion+"@latest")
	}

	execCommand(t, time.Minute, nil, "go"+compilerVersion, "download")
}

// installCompiler115 installs go1.15 or below by running the following commands:
//
//	export GOARCH=amd64
//	go run golang.org/dl/go1.15.15@latest download
//	go install golang.org/dl/go1.15.15@latest
//
// See # https://alexandear.github.io/posts/2024-07-12-old-go-darwin-arm64/
func installCompiler115(t *testing.T, compilerVersion string) {
	t.Helper()

	if compilerVersion == "" {
		t.Fatal("empty compilerVersion")
	}

	env := []string{"GOARCH=amd64"}
	execCommand(t, time.Minute, env, "go", "run", "golang.org/dl/go"+compilerVersion+"@latest", "download")
	execCommand(t, time.Minute, env, "go", "install", "golang.org/dl/go"+compilerVersion+"@latest")
}

func execCommand(t *testing.T, timeout time.Duration, env []string, cmd string, args ...string) (stdout string) {
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()

	var stdoutBuf, stderrBuf bytes.Buffer
	command := exec.CommandContext(ctx, cmd, args...)
	t.Log("Running", command, "with timeout", timeout)
	command.Stdout = &stdoutBuf
	command.Env = slices.Concat(command.Env, os.Environ(), env)
	command.Stderr = &stderrBuf
	if err := command.Run(); err != nil {
		t.Fatalf("Failed to download go via cmd %q: %v\nStdout: %s\nStderr: %s", command, err, stdoutBuf.String(), stderrBuf.String())
	}
	return stdoutBuf.String()
}

func findCaseNames(t *testing.T) map[string]struct{} {
	t.Helper()

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
	if len(caseNames) == 0 {
		t.Fatal("Missing cases")
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

func executeScript(t *testing.T, script string, timeout time.Duration, workDir string) (stdout string) {
	timeout = cmp.Or(timeout, 5*time.Second)

	t.Logf("Executing the script %q in %q with timeout %v", script, workDir, timeout)

	t.Chdir(workDir)
	return execCommand(t, timeout, nil, "sh", script)
}
