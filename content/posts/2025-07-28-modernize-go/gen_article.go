// Program generates ../2025-07-28-modernize-go.md article from 2025-07-28-modernize-go.tmpl.
// Run `go generate ./...` in this directory to update contents.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//go:generate go run $GOFILE

type Go struct {
	Version  string
	Sections []Section
}

type Section struct {
	name string

	Header      template.HTML
	Body        template.HTML
	Benefit     template.HTML
	FixCommands []template.HTML
	Before      template.HTML
	After       template.HTML
	Examples    []string
}

var gos = []Go{
	{
		Version: "1.24",
		Sections: []Section{
			{
				name:    "tchdir",
				Header:  "Replace os.Chdir with t.Chdir",
				Body:    "TODO",
				Benefit: "Simplifies testing code.",
				Examples: []string{
					"https://gitlab.com/gitlab-org/cli/-/merge_requests/2278/diffs#3ae6db62934a4153d302b878fd33bbbbeccb2aa9_101_95",
					"https://github.com/wagoodman/dive/pull/631/files#diff-fa257f4f9311442699f3ac132c9a981e2cbb5bcd435fe0d0228a16bf6753e332R14",
					"https://github.com/containers/podman/pull/26768/files#diff-355f1954b7cc2d7116308e9ae0c106cdb7e5867b67444f0e21d495229f688968R81",
				},
			},
		},
	},
	{
		Version: "1.23",
	},
	{
		Version: "1.22",
		Sections: []Section{
			{
				name:    "forloop",
				Header:  "Remove redundant loop variables",
				Body:    "A construct like `tc := tc` in `for` loops is not needed anymore and we can remove it.\nSee [Fixing For Loops in Go 1.22](https://go.dev/blog/loopvar-preview) for details.",
				Benefit: "Saves one line of code. Avoid weird for Go newbies constructions like `v := v`.",
				Examples: []string{
					"https://github.com/goreleaser/goreleaser/pull/4856/files#diff-3756619488c8c0f0c0300fc0cdcfecbb39c2a7bcb4fe4b3ac5305c6057512986L486",
					"https://github.com/kubernetes-sigs/kueue/pull/1946/files#diff-22ad2263a86a607fd28df7741c704614d0f34e208b5270153aa39427e4325fb3L203",
					"https://github.com/IBM/sarama/pull/3214/files#diff-cb488ad8239edeaaf8b0c1f469cc15c03fde53cbf22ee996e2f3922b3cc6a0c9L426",
					"https://github.com/google/go-github/pull/3537/files#diff-0f446fb8e4e16b655368f9f1c774d667d5528c9b3103f35481f704e2e33a925fL292",
					"https://github.com/go-critic/go-critic/pull/1459/files#diff-c2dfb8c940e1232344ce37c2a5942712765d9acf23d43c89345feb81fdbeeb13L43",
					"https://github.com/99designs/gqlgen/pull/3387/files#diff-fa4826c514673a47321901386ae757f00b2faa73d1433d8dacfc836f4928829aL44",
					"https://github.com/air-verse/air/pull/682/files#diff-0c22297be1ae696feec687c4dc3d1f425a6ff6c7dfd47d1d2a2275c32d3da14aL96",
					"https://github.com/nametake/golangci-lint-langserver/pull/62/files#diff-0eb779b9e49d8e44b0f36923fdb8d87d5ee024f886eefc45deec4ec88380a087L86",
				},
			},
			{
				name:    "forrange",
				Header:  "Simplify `for` range loops",
				Body:    "\"For\" loops may now range over integers.\nSee [For statements with range clause](https://go.dev/ref/spec#For_range) for details.",
				Benefit: "Improves readability and less symbols to type.",
				Examples: []string{
					"https://github.com/kubernetes-sigs/kueue/pull/5914/files#diff-539f3fc7450aa4c1e6682c00a20c862a4d603225852fdd26bce2fbe6d60ed044R148",
					"https://github.com/lima-vm/lima/pull/3399/files#diff-4fe57274e3aa074c4ccca2967546e5ad77ec58165d477f30560bef494c637e4dR180",
					"https://github.com/mgechev/revive/pull/1282/files#diff-75fa8cea7543dbb0e07700624e2760869a23cc2004dcb834e3e5a84739d25519L157",
				},
			},
		},
	},
	{
		Version: "1.21",
		Sections: []Section{
			{
				name:    "minmax",
				Header:  "Replace handwritten `min, max` or `math.Min`, `math.Max` functions with builtin `min`, `max`",
				Body:    "TODO",
				Benefit: "Simplifies code.",
				Examples: []string{
					"https://github.com/kubernetes-sigs/scheduler-plugins/pull/835/files#diff-a9d2a24a7e8778c1edaecdbfef1d7873cd2c9df69c24a1bc00d4e504de2fb4b8R227",
					"https://github.com/getkin/kin-openapi/pull/1032/files#diff-6b3cce991b5d47ed27df8dafc6ece7b16dc90449f6a14cd1d5cb7229a9c5920cR176",
					"https://github.com/nametake/golangci-lint-langserver/pull/62/files#diff-0eb779b9e49d8e44b0f36923fdb8d87d5ee024f886eefc45deec4ec88380a087L113-L119",
				},
			},
		},
	},
	{
		Version: "1.20",
		Sections: []Section{
			{
				name:   "slicearrconv",
				Header: "Simplify slice to array conversions",
				Body: `Can be implemented in the future via https://github.com/golang/go/issues/69820.
Exceptions: https://github.com/tailscale/tailscale/blob/5bb42e3018a0543467a332322f438cda98530c3a/net/connstats/stats_test.go#L28`,
				Benefit: "Simplifies code and improves readability.",
				Examples: []string{
					"https://github.com/SagerNet/sing-box/pull/3234/files#diff-0889dc2e2f6c8f4da1975681a711a6f5f8a4c31e91f41cdf63f9a42f79d233ccR99",
					"https://github.com/photoprism/photoprism/pull/5150/files#diff-182915d5d1268b03e71189f1c14b45481d6becf38d98a7c5e664caf8b1f60472R85",
					"https://github.com/tailscale/tailscale/pull/16778/files#diff-828b3be3fe317798f01a8c2d7ee04fcbcd96273ed464682051305357d8198c9eL302",
				},
			},
		},
	},
	{
		Version: "1.19",
		Sections: []Section{
			{
				name:    "atomictypes",
				Header:  "Replace atomic operations with atomic types",
				Body:    "https://cuonglm.xyz/post/go119_atomic_types/, https://go.dev/doc/go1.19#atomic_types",
				Benefit: "The code becomes much more readable and type safe.",
				Examples: []string{
					"https://github.com/lni/dragonboat/pull/397/files#diff-ee90aaffca546905162bc3fbc92559c58172bf8b1002935aa369cd73f9671710R54",
					"https://github.com/go-pg/pg/pull/2033/files#diff-c558bb9c370938c371ecd8203b543374c6c62ebe8c909a02dd4fb7fb87715e7cR23",
					"https://github.com/valyala/fasthttp/pull/2048/files#diff-fcd23a51a7b4023a9052d240502be2d458eb35b5d16339ac96a481330d68a7a3R199",
				},
			},
		},
	},
	{
		Version: "1.18",
	},
	{
		Version: "1.17",
		Sections: []Section{
			{
				name:    "tsetenv",
				Header:  "Simplify setting of environment variables in tests with t.Setenv",
				Body:    "Applies also to benchmarks and fuzzy tests.",
				Benefit: "Simplifies testing code.",
				Examples: []string{
					"https://github.com/go-delve/delve/pull/3503/files#diff-012a9c9759129e99b4bef0bc00f8b572d4df044215b32d2a8b25acfc0c012bffR4482",
					"https://github.com/rqlite/rqlite/pull/2037/files#diff-2129ccb7b21d5d1b042238fb694c52997805126fae4d7effaea833de4b9950ebR48",
					"https://github.com/jackc/pgx/pull/2014/files#diff-05829baf37fc98f2ff986da920439177057a3b11e89b066d98ac6e7f655af093R1012",
					"https://github.com/go-swagger/go-swagger/pull/2914/files#diff-ed129a1818e6ed02001cb40223cf4fe0492abc18ab859d464f007c7c33360a28R130",
				},
			},
		},
	},
	{
		Version: "1.16",
	},
	{
		Version: "1.15",
		Sections: []Section{
			{
				name:    "ttempdir",
				Header:  "Simplify temporary file handling in tests with t.TempDir",
				Body:    "TODO",
				Benefit: "Simplifies testing code.",
				Examples: []string{
					"https://github.com/rqlite/rqlite/pull/2022/files#diff-ac1924166edf5c219c1ce8562a0b345c7358d96fc9a28a9cbcf23b084e6a5705R18",
					"https://github.com/spf13/viper/pull/1631/files#diff-cbedff684a75d0a2cfdc7d4234fb3573dd7af147878776e039d69e1694a18b34R249",
					"https://github.com/gohugoio/hugo/pull/10944/files#diff-570e9a48e4df86fdea8e1465b52882220ab7e6ee01f85ef7a0352e04315caf73R263",
					"https://gitlab.com/gitlab-org/api/client-go/-/merge_requests/2355/diffs#b65b37f51cbf5979f82400da21dd968d0ac02885_57_55",
					"https://github.com/securego/gosec/pull/1265/files#diff-67271f60deaa2e51dcee16995bc6eb5e599dc120dd47dddc41137dca5a1a6509R47",
				},
			},
		},
	},
	{
		Version: "1.14",
	},
	{
		Version: "1.13",
	},
	{
		Version: "1.12",
	},
	{
		Version: "1.11",
	},
	{
		Version: "1.10",
	},
	{
		Version: "1.9",
	},
	{
		Version: "1.8",
	},
	{
		Version: "1.7",
	},
	{
		Version: "1.6",
	},
	{
		Version: "1.5",
	},
	{
		Version: "1.4",
	},
	{
		Version: "1.3",
	},
	{
		Version: "1.2",
	},
	{
		Version: "1.1",
	},
}

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "If print debug statements")
	flag.Parse()

	funcMap := template.FuncMap{
		"formatRepoLink": formatRepoLink,
	}

	const articleFilename = "2025-07-28-modernize-go"
	tmplFilename := articleFilename + ".tmpl"
	tmpl, err := template.New(tmplFilename).Funcs(funcMap).ParseFiles(tmplFilename)
	if err != nil {
		log.Fatal(err)
	}

	for _, goElem := range gos {
		if len(goElem.Sections) == 0 {
			log.Printf("No sections for Go version %q\n", goElem.Version)
			continue
		}
		func() {
			if err := os.Chdir(goElem.Version); err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := os.Chdir(".."); err != nil {
					log.Panic(err)
				}
			}()

			for i := range goElem.Sections {
				section := goElem.Sections[i]

				scripts, err := filepath.Glob(section.name + "*.sh")
				if err != nil {
					log.Fatal(err)
				}

				if len(scripts) == 0 {
					log.Fatalf("Missed scripts for section %q in version %q", section.name, goElem.Version)
				}

				for _, script := range scripts {
					command, err := extractShContent(script)
					if err != nil {
						log.Fatal(err)
					}
					goElem.Sections[i].FixCommands = append(goElem.Sections[i].FixCommands, command)
				}

				beforeFilename, err := beforeFilename(section.name)
				if err != nil {
					log.Fatal(err)
				}

				before, err := extractGoContent(beforeFilename)
				if err != nil {
					log.Fatal(err)
				}
				goElem.Sections[i].Before = before

				afterFilename, err := afterFilename(section.name)
				if err != nil {
					log.Fatal(err)
				}
				after, err := extractGoContent(afterFilename)
				if err != nil {
					log.Fatal(err)
				}
				goElem.Sections[i].After = template.HTML(after)
			}
		}()
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, gos)
	if err != nil {
		log.Fatal(err)
	}

	if debug {
		fmt.Println(buf.String())
	}

	// workaround
	b := bytes.Replace(buf.Bytes(), []byte("&lt;!--more--&gt;"), []byte("<!--more-->"), 1)

	err = os.WriteFile(filepath.Join("..", articleFilename+".md"), b, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func beforeFilename(name string) (string, error) {
	files, err := filepath.Glob(name + "_before*.go")
	if err != nil {
		return "", err
	}
	if n := len(files); n != 1 {
		return "", fmt.Errorf("should be only one before file, but found %d", n)
	}
	return files[0], err
}

func afterFilename(name string) (string, error) {
	files, err := filepath.Glob(name + "_after*.go")
	if err != nil {
		return "", err
	}
	if n := len(files); n != 1 {
		return "", fmt.Errorf("should be only one after file, but found %d", n)
	}
	return files[0], err
}

const (
	snippetBegin = "<< snippet begin >>"
	snippetEnd   = "<< snippet end >>"
)

func extractShContent(filename string) (template.HTML, error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	begin := "# " + snippetBegin + "\n"
	_, after, ok := strings.Cut(string(contents), begin)
	if !ok {
		return "", fmt.Errorf("missed %q in file %q", begin, filename)
	}
	end := "# " + snippetEnd
	before, _, ok := strings.Cut(after, end)
	if !ok {
		return "", fmt.Errorf("missed %q in file %q", end, filename)
	}
	return template.HTML(strings.TrimSpace(before)), nil
}

func extractGoContent(filename string) (template.HTML, error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	begin := "// " + snippetBegin + "\n"
	_, after, ok := strings.Cut(string(contents), begin)
	if !ok {
		return "", fmt.Errorf("missed %q in file %q", begin, filename)
	}
	end := "// " + snippetEnd
	before, _, ok := strings.Cut(after, end)
	if !ok {
		return "", fmt.Errorf("missed %q in file %q", end, filename)
	}
	return template.HTML(before), nil
}

func formatRepoLink(url string) template.HTML {
	// Extract repository name from GitHub or GitLab URL
	// Pattern: https://github.com/owner/repo/... or https://gitlab.com/owner/repo/...
	re := regexp.MustCompile(`https://(github|gitlab)\.com/([^/]+/[^/]+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 3 {
		return template.HTML(fmt.Sprintf("[%s](%s)", url, url))
	}

	repoName := matches[2]
	return template.HTML(fmt.Sprintf("[%s](%s)", repoName, url))
}
