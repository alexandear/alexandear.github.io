package main

import (
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractShContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    template.HTML
		wantErr error
	}{
		{
			name: "valid shell script with content",
			content: `#!/bin/sh

# << snippet begin >>
echo "hello world"
ls -la
# << snippet end >>

exit 0`,
			want: "echo \"hello world\"\nls -la",
		},
		{
			name: "empty content between markers",
			content: `#!/bin/sh

# << snippet begin >>
# << snippet end >>`,
			want: "",
		},
		{
			name: "content with whitespace",
			content: `#!/bin/sh

# << snippet begin >>

   echo "test"

# << snippet end >>`,
			want: "echo \"test\"",
		},
		{
			name: "missing begin marker",
			content: `#!/bin/sh
echo "test"
# << snippet end >>`,
			wantErr: errors.New(`missed "# << snippet begin >>\n" in file "`),
		},
		{
			name: "missing end marker",
			content: `#!/bin/sh
# << snippet begin >>
echo "test"`,
			wantErr: errors.New(`missed "# << snippet end >>" in file "`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := filepath.Join(t.TempDir(), "test.sh")
			if err := os.WriteFile(script, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}

			result, err := extractShContent(script)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("Expected error, got <nil>")
				}
				if !strings.HasPrefix(err.Error(), tt.wantErr.Error()) {
					t.Fatalf("Mismatch error: got %q, want %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tt.want {
				t.Errorf("Unexpected result: got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestExtractGoContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    template.HTML
		wantErr error
	}{
		{
			name: "valid Go file with content",
			content: `package main

import "fmt"

// << snippet begin >>
func main() {
	fmt.Println("hello")
}
// << snippet end >>`,
			want: "func main() {\n\tfmt.Println(\"hello\")\n}\n",
		},
		{
			name: "empty content between markers",
			content: `package main

// << snippet begin >>
// << snippet end >>`,
			want: "",
		},
		{
			name: "content with variables",
			content: `package main

// << snippet begin >>
var x = 42
var y = "test"
// << snippet end >>`,
			want: "var x = 42\nvar y = \"test\"\n",
		},
		{
			name: "missing begin marker",
			content: `package main

var x = 42
// << snippet end >>`,
			want:    "",
			wantErr: errors.New(`missed "// << snippet begin >>\n" in file "`),
		},
		{
			name: "missing end marker",
			content: `package main

// << snippet begin >>
var x = 42`,
			want:    "",
			wantErr: errors.New(`missed "// << snippet end >>" in file "`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goFile := filepath.Join(t.TempDir(), "test.go")
			err := os.WriteFile(goFile, []byte(tt.content), 0o644)
			if err != nil {
				t.Fatal(err)
			}

			result, err := extractGoContent(goFile)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("Expected error, got <nil>")
				}
				if !strings.HasPrefix(err.Error(), tt.wantErr.Error()) {
					t.Fatalf("Mismatch error: got %q, want %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tt.want {
				t.Errorf("Unexpected result: got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestFormatGitHubLink(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected template.HTML
	}{
		{
			name:     "pull request URL",
			input:    "https://github.com/goreleaser/goreleaser/pull/4856/files#diff-3756619488c8c0f0c0300fc0cdcfecbb39c2a7bcb4fe4b3ac5305c6057512986L486",
			expected: "[goreleaser/goreleaser](https://github.com/goreleaser/goreleaser/pull/4856/files#diff-3756619488c8c0f0c0300fc0cdcfecbb39c2a7bcb4fe4b3ac5305c6057512986L486)",
		},
		{
			name:     "simple repo URL",
			input:    "https://github.com/kubernetes-sigs/kueue",
			expected: "[kubernetes-sigs/kueue](https://github.com/kubernetes-sigs/kueue)",
		},
		{
			name:     "issues URL",
			input:    "https://github.com/golang/go/issues/69820",
			expected: "[golang/go](https://github.com/golang/go/issues/69820)",
		},
		{
			name:     "blob URL",
			input:    "https://github.com/tailscale/tailscale/blob/5bb42e3018a0543467a332322f438cda98530c3a/net/connstats/stats_test.go#L28",
			expected: "[tailscale/tailscale](https://github.com/tailscale/tailscale/blob/5bb42e3018a0543467a332322f438cda98530c3a/net/connstats/stats_test.go#L28)",
		},
		{
			name:     "repo with hyphens and numbers",
			input:    "https://github.com/99designs/gqlgen/pull/3387/files",
			expected: "[99designs/gqlgen](https://github.com/99designs/gqlgen/pull/3387/files)",
		},
		{
			name:     "non-GitHub URL",
			input:    "https://example.com/some/path",
			expected: "[https://example.com/some/path](https://example.com/some/path)",
		},
		{
			name:     "malformed URL",
			input:    "not-a-url",
			expected: "[not-a-url](not-a-url)",
		},
		{
			name:     "GitHub URL without repo",
			input:    "https://github.com/",
			expected: "[https://github.com/](https://github.com/)",
		},
		{
			name:     "GitHub URL with only owner",
			input:    "https://github.com/golang",
			expected: "[https://github.com/golang](https://github.com/golang)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatGitHubLink(tt.input)
			if result != tt.expected {
				t.Errorf("formatGitHubLink(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
