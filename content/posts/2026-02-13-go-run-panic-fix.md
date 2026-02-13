---
title: "How to fix panic in go command"
date: 2026-02-13
draft: true
---

TL;DR: The story on how I find and fixed the bug in `go run -C`.

Note: I'm using Go 1.26:

```console
$ go version
go version go1.26.0 darwin/arm64
```

Recently I work with `go run` and `-C` flag and accidentally spot the panic occurred when I didn't provide an argument to the flag:

```console
$ go run -C
panic: runtime error: slice bounds out of range [1:0]

goroutine 1 [running]:
cmd/go/internal/toolchain.maybeSwitchForGoInstallVersion({0x10273bbca, 0x6})
        cmd/go/internal/toolchain/select.go:665 +0xa98
cmd/go/internal/toolchain.Select()
        cmd/go/internal/toolchain/select.go:234 +0xbc8
main.main()
        cmd/go/main.go:107 +0x54
```

`go run` compiles and runs the named main Go package. `go run main.go` is similar to `go build -o app main.go && ./app`, but the `go run` uses a temporary binary.

`-C` is a build flag to change the dir before running the command. From the Go 1.26 help:

```console
$ go help build | head -39 | tail -8
The build flags are shared by the build, clean, get, install, list, run,
and test commands:

        -C dir
                Change to dir before running the command.
                Any files named on the command line are interpreted after
                changing directories.
                If used, this flag must be the first one in the command line.
```

It's like `cd`. Here are roughly the same lines:

```console
$ go run -C example main.go
$ cd example && go run main.go
```

Let's check the `build`:

```console
$ go build -C
flag needs an argument: -C
usage: go build [-o output] [build flags] [packages]
Run 'go help build' for details.
```

And the same for `go run` produces the panic "runtime error: slice bounds out of range [1:0]".
It is definitely a bug.

First of all, let's try if the bug reproduced with the development `go` version.
Maybe the panic has already been addressed.

The fastest way to try it is to install [`gotip`](https://pkg.go.dev/golang.org/dl@v0.0.0-20260210192738-6a105f684182/gotip).
I use the following commands:

```console
$ go run golang.org/dl/gotip@latest download
Cloning into '/Users/alexandear/sdk/gotip'...
remote: Counting objects: 16315, done
remote: Finding sources: 100% (16315/16315)
remote: Total 16315 (delta 2076), reused 10885 (delta 2076)
Receiving objects: 100% (16315/16315), 34.54 MiB | 7.68 MiB/s, done.
Resolving deltas: 100% (2076/2076), done.
Updating files: 100% (14980/14980), done.
Updating the go development tree...
From https://go.googlesource.com/go
 * branch            master     -> FETCH_HEAD
HEAD is now at d4febb4 crypto/tls: avoid data race when canceling a QUICConn''s Context
Building Go cmd/dist using /opt/homebrew/Cellar/go/1.25.7_1/libexec. (go1.25.7 darwin/arm64)
Building Go toolchain1 using /opt/homebrew/Cellar/go/1.25.7_1/libexec.
Building Go bootstrap cmd/go (go_bootstrap) using Go toolchain1.
Building Go toolchain2 using go_bootstrap and Go toolchain1.
Building Go toolchain3 using go_bootstrap and Go toolchain2.
Building packages and commands for darwin/arm64.
---
Installed Go for darwin/arm64 in /Users/alexandear/sdk/gotip
Installed commands in /Users/alexandear/sdk/gotip/bin
Success. You may now run 'gotip'!
$ go install golang.org/dl/gotip@latest
$ gotip version
go version go1.27-devel_d4febb4 Thu Feb 5 17:05:55 2026 -0800 darwin/arm64
```

And it reproduced again:

```console
$ gotip run -C
panic: runtime error: slice bounds out of range [1:0]

goroutine 1 [running]:
cmd/go/internal/toolchain.maybeSwitchForGoInstallVersion(0x29ffab1703c0, {0x1052c0782, 0x4})
        /Users/alexandear/sdk/gotip/src/cmd/go/internal/toolchain/select.go:662 +0xab8
cmd/go/internal/toolchain.Select()
        /Users/alexandear/sdk/gotip/src/cmd/go/internal/toolchain/select.go:236 +0xcc8
main.main()
        /Users/alexandear/sdk/gotip/src/cmd/go/main.go:106 +0x54
```

So, let's open to the [Go issue tracker](https://go.dev/issue) to search if someone has been reported this.

We should filter by the **open** issues with **label: GoCommand** and the **panic** keyword: I used [this query](https://go.dev/issue?q=is%3Aissue%20state%3Aopen%20label%3AGoCommand%20panic).

{{< figure src="/img/2026-02-13-go-run-panic-fix/go-open-panic-issues.png" width="100%" alt="Screenshot of the open panic issues for the go command" >}}

But no one is similar to what I'm looking.
So I opened [go.dev/issue/77483](https://go.dev/issue/77483).

{{< figure src="/img/2026-02-13-go-run-panic-fix/go-issue-77483.png" width="100%" alt="Screenshot of the issue 77483 reporting about the panic in go" >}}

After that I tried to fix it by myself.
First of all, we need to clone the Go repo and setup necessary tools for contribution.
I already wrote about this in detail in the article ["How to contribute to the Go language"](../2025-01-31-go-simple-contrib/).

Open the cloned repo and look at the directory structure:

```console
$ cd go
$ ls | col
CONTRIBUTING.md
LICENSE
PATENTS
README.md
SECURITY.md
VERSION.cache
api
bin
codereview.cfg
doc
go.env
lib
misc
pkg
src
test
```

Next, we should navigate to the `src` directory and build the `go`:

```console
$ cd src
$ ./make.bash
Building Go cmd/dist using /opt/homebrew/Cellar/go/1.26.0/libexec. (go1.26.0 darwin/arm64)
Building Go toolchain1 using /opt/homebrew/Cellar/go/1.26.0/libexec.
Building Go bootstrap cmd/go (go_bootstrap) using Go toolchain1.
Building Go toolchain2 using go_bootstrap and Go toolchain1.
Building Go toolchain3 using go_bootstrap and Go toolchain2.
Building packages and commands for darwin/arm64.
---
Installed Go for darwin/arm64 in /Users/alexandear/src/go.googlesource.com/go
Installed commands in /Users/alexandear/src/go.googlesource.com/go/bin
*** You need to add /Users/alexandear/src/go.googlesource.com/go/bin to your PATH.
```

Let's check again:

```console
$ ../bin/go run -C
panic: runtime error: slice bounds out of range [1:0]

goroutine 1 [running]:
cmd/go/internal/toolchain.maybeSwitchForGoInstallVersion(0x7dba035b8300, {0x10469c582, 0x4})
        /Users/alexandear/src/go.googlesource.com/go/src/cmd/go/internal/toolchain/select.go:662 +0xab8
cmd/go/internal/toolchain.Select()
        /Users/alexandear/src/go.googlesource.com/go/src/cmd/go/internal/toolchain/select.go:236 +0xcc8
main.main()
        /Users/alexandear/src/go.googlesource.com/go/src/cmd/go/main.go:106 +0x54
```

And it's panic again.

The fix should be easy because the panic log show where the problem lives: it's in the `maybeSwitchForGoInstallVersion` on the line 662:

```go
// maybeSwitchForGoInstallVersion reports whether the command line is go install m@v or go run m@v.
// If so, switch to the go version required to build m@v if it's higher than minVers.
func maybeSwitchForGoInstallVersion(loaderstate *modload.State, minVers string) {
		// ...
		if bf, ok := f.Value.(interface{ IsBoolFlag() bool }); !ok || !bf.IsBoolFlag() {
			// The next arg is the value for this flag. Skip it.
			args = args[1:]
			continue
		}
		// ...
}
```

The problematic line is `args = args[1:]`, which panics when `len(args)` is `0`.

Before we actually fixing, we spot that `maybeSwitchForGoInstallVersion` responsible for `go run` and `go install`. Let's check this:

```console
$ ../bin/go install -C
panic: runtime error: slice bounds out of range [1:0]

goroutine 1 [running]:
cmd/go/internal/toolchain.maybeSwitchForGoInstallVersion(0x79255db90360, {0x10131c582, 0x4})
        /Users/alexandear/src/go.googlesource.com/go/src/cmd/go/internal/toolchain/select.go:662 +0xab8
cmd/go/internal/toolchain.Select()
        /Users/alexandear/src/go.googlesource.com/go/src/cmd/go/internal/toolchain/select.go:236 +0xcc8
main.main()
        /Users/alexandear/src/go.googlesource.com/go/src/cmd/go/main.go:106 +0x54
```

And here we go, panic again. The fix should fix two panics: one in `go run -C` and one in `go install -C`.

After some digging I found that the fix is simple as:

```

```

https://go.dev/cl/742860
