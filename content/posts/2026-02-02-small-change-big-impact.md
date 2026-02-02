---
title: "How I Simplified go-github for Millions"
date: 2026-02-02
tags: ["opensource", "github", "go", "api"]
---

How can a simple contribution to a library impact thousands of programs around the world?
And why will this change become obsolete after the Go 1.26 release?
These and a few other questions I'll answer in detail in this post.

**TL;DR**: I replaced four helper functions (`String`, `Bool`, `Int`, `Int64`) with a single generic `Ptr[T]` function in [google/go-github][], affecting 10K+ projects.
Ironically, Go 1.26's enhanced `new` builtin will make this pattern unnecessary.

{{< figure src="/img/2026-02-02-small-change-big-impact/go-github-ptr.webp" width="100%" alt="Screenshot of the Ptr function" >}}

<!--more-->

## The contribution

Replacing four helper routines with one generic function `github.Ptr` in [PR][] in the [google/go-github][] project is my largest and most impactful contribution ever.
And here's why I think so.

## google/go-github

The `google/go-github`, the Go library for accessing the GitHub REST API, is a widely popular package with over [11K stars](https://github.com/google/go-github/stargazers) on GitHub.
[GitHub shows](https://github.com/google/go-github/network/dependents) that it's used by more than 10K programs across the GitHub ecosystem.
The [Go Reference](https://pkg.go.dev/github.com/google/go-github/github) shows that this package is imported more than 7.5K times.
However, this understates the real usage since the library has over 80 major versions, each tracked separately.
SourceGraph [lists 10K+ go-github imports](https://sourcegraph.com/search?q=context:global+github.com/google/go-github+lang:go+&patternType=standard&sm=1) and it's not a limit.
Also, we don't know how many closed-source libraries import go-github.
The [google/go-github][] is a [backbone](https://github.com/github/github-mcp-server/blob/1820a0ff11cc2ff91476066ab3c6088e9f6dfa1e/go.mod#L6) for the [GitHub MCP server](https://github.com/github/github-mcp-server) (26K stars).
The MCP server is used by every AI agent that is doing something on GitHub.

## Technical details

From a technical perspective, replacing a few helper functions with one generic function is not a big deal.
But these helper functions are used by everyone who interacts with the GitHub API.
This means millions of developers benefit from this small change.

### Taking an address of a constant

In Go, there is no simple way to create pointers to non-zero bool/numeric/string values.
For example, the following code does [not compile:](https://go.dev/play/p/0DZNC09cYkk)

```go
var a = &1234 // invalid operation: cannot take address of 1234 (untyped int constant)
var b = &true // invalid operation: cannot take address of true (untyped bool constant)
var c = &"hi" // invalid operation: cannot take address of "hi" (untyped string constant)
```

To overcome this, we should use an additional variable which leads to more [verbose code:](https://go.dev/play/p/2FXBEIsxHGR)

```go
var num = 1234
var a = &num
var tr = true
var b = &tr
var str = "hi"
var c = &str
```

The [Go101](https://go101.org/) tells, that we can also use these one-liners via [these constructs](https://go.dev/play/p/JgJAMab2Tfl):

```go
var a = &(&[1]int{1234})[0]
var b = &(&[1]bool{true})[0]
var c = &(&[1]string{"hi"})[0]
```

Or slightly [less efficient](https://go.dev/play/p/fzKy68ZEQ5G):

```go
var a = &([]int{1234})[0]
var b = &([]bool{true})[0]
var c = &([]string{"hi"})[0]
```

But using an [anonymous functions](https://go.dev/play/p/We491YvxknU) is the most popular way to achieve that:

```go
var a = func(v int) *int { return &v }(1234)
var b = func(v bool) *bool { return &v }(true)
var c  = func(v string) *string { return &v }("hi")
```

### Helpers for optional values in go-github

Optional values are heavily used in [google/go-github][] functions.

The authors of `google/go-github` implemented the third way but moved to named functions to achieve that:

```go
func Int(v int) *int { return &v }

func Bool(v bool) *bool { return &v }

func Int64(v int64) *int64 { return &v }

func String(v string) *string { return &v }
```

This example is from the [README](https://github.com/google/go-github/tree/v67.0.0?tab=readme-ov-file#creating-and-updating-resources):

```go
repo := &github.Repository{
	Name:    github.String("foo"),
	Private: github.Bool(true),
}
client.Repositories.Create(ctx, "", repo)
```

### My pull request to replace helper functions with `github.Ptr`

My [PR][] replaces these four helper functions with one that uses Go 1.18 generics:

```go
func Ptr[T any](v T) *T { return &v }
```

And we can use it as:

```go
var a = Ptr(1234)
var b = Ptr(true)
var c = Ptr("hi")
```

In [google/go-github][]:

```go
repo := &github.Repository{
	Name:    github.Ptr("foo"),
	Private: github.Ptr(true),
}
client.Repositories.Create(ctx, "", repo)
```

To be more precise, it should restrict also `T` to only four types used in `google/go-github`:

```go
func Ptr[T bool | int | int64 | string](v T) *T { return &v }
```

But we decided to use `any` (`interface{}`).

This change is only adding one line and 20 lines of tests. But [PR][]'s changes are 15K additions and 15K deletions.
Why is that big?

It's because we need to replace `github.String`/`github.Int`/`github.Int64`/`github.Bool` to `github.Ptr` in all existing examples and tests.

## Why `Ptr` will become obsolete in Go 1.26

Interestingly, just as my contribution unified the helper functions into a single generic `Ptr`, the Go team is about to make it obsolete.

The `Ptr` function is [highly popular](https://pkg.go.dev/search?q=ptr) and exists in every major project:

- [kubernetes](https://github.com/kubernetes/utils/blob/914a6e750570/ptr/ptr.go#L50) (called `ptr.To`)
- [tailscale](https://github.com/tailscale/tailscale/blob/abdbca47af098469fba238c408dd1f4b373d254c/types/ptr/ptr.go#L8)
- [gitlab-org/api/client-go](https://gitlab.com/gitlab-org/api/client-go/-/blob/1f279e473b9fd7a145db0dd099527c2b1f81d881/types.go#L32)
- [getkin/kin-openapi](https://github.com/getkin/kin-openapi/blob/45db2adb0102203579276f24456ce494b59e751a/openapi3/helpers.go#L35) (my [PR to add it](https://github.com/getkin/kin-openapi/pull/1033))

But Go [1.26](https://go.dev/doc/go1.26#language) allows passing an expression to the `new` builtin function.
It is going to [be released](https://groups.google.com/g/golang-dev/c/KnvIaNZU8XQ) in two weeks, but we can try it in the [Gotip Playground](https://gotipplay.golang.org/).

So, we can simplify our initial example to the [following](https://gotipplay.golang.org/p/9ATIHserid7):

```go
var a = new(1234)
var b = new(true)
var c = new("hi")
```

And in [google/go-github][]:

```go
repo := &github.Repository{
	Name:    new("foo"),
	Private: new(true),
}
client.Repositories.Create(ctx, "", repo)
```

## Conclusion

Open source contributions can be small and simple yet have a huge impact on the world.
My one-line generic function is now used by 10K+ projects and powers every AI agent interacting with GitHub.

Yes, Go 1.26 will make `Ptr` obsolete — but that's the nature of software.
The real value was not just the code. It was finding an opportunity to simplify something used by millions.

Keep reading release notes, keep contributing, and don't be afraid to make small changes.
Sometimes they matter more than you think.

[google/go-github]: https://github.com/google/go-github
[PR]: https://github.com/google/go-github/pull/3355
