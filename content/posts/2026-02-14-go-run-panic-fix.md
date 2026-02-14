---
title: "Fixing a panic in 'go run' command"
date: 2026-02-14
tags: ["go", "panic", "opensource"]
---

Have you found bugs in tools that you are using?
Recently, while working on some [google/go-github](../2026-02-02-small-change-big-impact/) issues,
I discovered a panic in `go run`.
It happened when I forgot to provide an argument to the `-C` flag.

<!-- https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=seti&wt=none&l=application%2Fx-sh&width=800&ds=false&dsyoff=18px&dsblur=68px&wc=true&wa=false&pv=9px&ph=100px&ln=false&fl=1&fm=Hack&fs=14px&lh=142%25&si=false&es=2x&wm=false&code=%25E2%259D%25AF%2520go%2520run%2520-C%250Apanic%253A%2520runtime%2520error%253A%2520slice%2520bounds%2520out%2520of%2520range%2520%255B1%253A0%255D%250A%250Agoroutine%25201%2520%255Brunning%255D%253A%250Acmd%252Fgo%252Finternal%252Ftoolchain.maybeSwitchForGoInstallVersion%28%257B0x10273bbca%252C%25200x6%257D%29%250A%2520%2520%2520%2520%2520%2520%2520%2520cmd%252Fgo%252Finternal%252Ftoolchain%252Fselect.go%253A665%2520%252B0xa98%250Acmd%252Fgo%252Finternal%252Ftoolchain.Select%28%29%250A%2520%2520%2520%2520%2520%2520%2520%2520cmd%252Fgo%252Finternal%252Ftoolchain%252Fselect.go%253A234%2520%252B0xbc8%250Amain.main%28%29%250A%2520%2520%2520%2520%2520%2520%2520%2520cmd%252Fgo%252Fmain.go%253A107%2520%252B0x54%250AHello%252C%2520World%21 -->
{{< figure src="/img/2026-02-14-go-run-panic-fix/go-run-c-panic.webp" width="100%" alt="Screenshot of the panic in go run -C" >}}

This article explains how I found and fixed the bug in `go run -C` and `go install -C`, added tests to prevent future regressions, and how my fix was successfully accepted into the Go repository.

<!--more-->

## Reproducing the panic

The panic is easy to reproduce on Go 1.26 or 1.25:

```console
❯ go version
go version go1.26.0 darwin/arm64
❯ go run -C
panic: runtime error: slice bounds out of range [1:0]

goroutine 1 [running]:
cmd/go/internal/toolchain.maybeSwitchForGoInstallVersion({0x10273bbca, 0x6})
        cmd/go/internal/toolchain/select.go:665 +0xa98
cmd/go/internal/toolchain.Select()
        cmd/go/internal/toolchain/select.go:234 +0xbc8
main.main()
        cmd/go/main.go:107 +0x54
```

## Understanding the `-C` flag

The `go run` command compiles and runs a Go program.
`go run main.go` is like running `go build -o app main.go && ./app`, but it uses a temporary file.

The `-C` build flag changes the working directory before running the command.
From the Go 1.26 help:

```console
❯ go help build | head -39 | tail -8
The build flags are shared by the build, clean, get, install, list, run,
and test commands:

        -C dir
                Change to dir before running the command.
                Any files named on the command line are interpreted after
                changing directories.
                If used, this flag must be the first one in the command line.
```

The `-C` flag works like `cd`. These commands are roughly equivalent:

```console
❯ go run -C example main.go
❯ cd example && go run main.go
```

Now let's see how `go build` handles a missing argument:

```console
❯ go build -C
flag needs an argument: -C
usage: go build [-o output] [build flags] [packages]
Run 'go help build' for details.
```

But `go run -C` produces a panic instead of a helpful error message.
This is definitely a bug.

## Confirming the bug in gotip

First, I checked whether this bug also exists in the development version of Go.
Maybe it was already fixed.

The quickest way to test this is to install the [`gotip`](https://pkg.go.dev/golang.org/dl/gotip) binary:

```console
❯ go run golang.org/dl/gotip@latest download
Cloning into '/Users/alexandear/sdk/gotip'...
remote: Counting objects: 16315, done
remote: Finding sources: 100% (16315/16315)
...
Building packages and commands for darwin/arm64.
---
Installed Go for darwin/arm64 in /Users/alexandear/sdk/gotip
Installed commands in /Users/alexandear/sdk/gotip/bin
Success. You may now run 'gotip'!
❯ go install golang.org/dl/gotip@latest
❯ gotip version
go version go1.27-devel_d4febb4 Thu Feb 5 17:05:55 2026 -0800 darwin/arm64
```

The bug still exists in the development version:

```console
❯ gotip run -C
panic: runtime error: slice bounds out of range [1:0]

goroutine 1 [running]:
cmd/go/internal/toolchain.maybeSwitchForGoInstallVersion(0x29ffab1703c0, {0x1052c0782, 0x4})
...
```

## Searching for existing issues

Next, I checked the [Go issue tracker](https://go.dev/issue) to see if anyone had already reported this.

I filtered for **open** issues with the **GoCommand** label and searched for "panic" using [this query](https://go.dev/issue?q=is%3Aissue%20state%3Aopen%20label%3AGoCommand%20panic).

{{< figure src="/img/2026-02-14-go-run-panic-fix/go-open-panic-issues.webp" width="100%" alt="Screenshot of the open panic issues for the go command" >}}

I didn't find a similar issue, so I opened [issue #77483](https://go.dev/issue/77483).

{{< figure src="/img/2026-02-14-go-run-panic-fix/go-issue-77483.webp" width="100%" alt="Screenshot of issue 77483 reporting the panic in go" >}}

## Setting up the development environment

To fix the bug myself, I needed to set up a development environment.
I've written about this in detail in ["How to contribute to the Go language"](../2025-01-31-go-simple-contrib/).

First, I cloned the Go repository, navigated to the `src` directory, and built the `go` binary:

```console
❯ cd src
❯ ./make.bash
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

I confirmed the bug was still present:

```console
❯ ../bin/go run -C
panic: runtime error: slice bounds out of range [1:0]

goroutine 1 [running]:
cmd/go/internal/toolchain.maybeSwitchForGoInstallVersion(0x7dba035b8300, {0x10469c582, 0x4})
        /Users/alexandear/src/go.googlesource.com/go/src/cmd/go/internal/toolchain/select.go:662 +0xab8
cmd/go/internal/toolchain.Select()
        /Users/alexandear/src/go.googlesource.com/go/src/cmd/go/internal/toolchain/select.go:236 +0xcc8
main.main()
        /Users/alexandear/src/go.googlesource.com/go/src/cmd/go/main.go:106 +0x54
```

## Identifying the root cause

The [stack trace](https://pkg.go.dev/runtime) from the panic pointed me directly to the problem: the `maybeSwitchForGoInstallVersion` function in `select.go` at line 662.

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

The problematic line is `args = args[1:]`.
This causes a panic when `len(args)` is `0` because slicing beyond the bounds is invalid.

I also checked whether `go install -C` has the same problem:

```console
❯ ../bin/go install -C
panic: runtime error: slice bounds out of range [1:0]

goroutine 1 [running]:
...
```

Yes, it does. So the fix needs to solve both issues.

## Implementing the fix

The fix is straightforward: add a bounds check before slicing:

```go
func maybeSwitchForGoInstallVersion(loaderstate *modload.State, minVers string) {
		// ...
		if bf, ok := f.Value.(interface{ IsBoolFlag() bool }); !ok || !bf.IsBoolFlag() {
			if len(args) == 0 {
				return
			}
			// The next arg is the value for this flag. Skip it.
			args = args[1:]
			continue
		}
		// ...
}
```

After rebuilding the binary, both commands now produce the correct error:

```console
❯ ./make.bash
...
❯ ../bin/go run -C
flag needs an argument: -C
usage: go run [build flags] [-exec xprog] package [arguments...]
Run 'go help run' for details.
❯ ../bin/go install -C
flag needs an argument: -C
usage: go install [build flags] [packages]
Run 'go help install' for details.
```

Perfect! Now `go run` and `go install` behave the same as `go build`.

## Writing a regression test

According to the "Software Engineering at Google" book,
["addressing a bug with a revised test is often necessary"](https://abseil.io/resources/swe-book/html/ch09.html#bug_fixes_and_rollbacks).

The Go repository has an excellent testing system for the command in the [`src/cmd/go/testdata/script`](https://go.googlesource.com/go/+/refs/tags/go1.26.0/src/cmd/go/testdata/script) directory.
These test scripts use a simple DSL and are executed during `go test cmd/go`.
I recommend reading the [README](https://go.googlesource.com/go/+/refs/tags/go1.26.0/src/cmd/go/testdata/script/README) and watching Russ Cox's [great talk](https://research.swtch.com/testing) about testing.

I created the test file `mod_run_flags_issue77483.txt`:

```txt
# Regression test for https://go.dev/issue/77483: 'go run -C' should not panic.

! go run -C
stderr 'flag needs an argument: -C'
```

Here's what each line does:

* Line 1: A comment explaining the test's purpose.
* Line 3: The `!` prefix indicates the command should fail.
* Line 4: `stderr` verifies the error message contains `flag needs an argument: -C`.

{{< note >}}

I didn't write a test for `go install -C`. You can do it!

{{< /note >}}

To run this specific test:

```console
❯ ../bin/go test cmd/go -run=Script/^mod_run_flags_issue77483$
ok      cmd/go  0.455s
```

## Verifying the test caught the bug

To confirm the test actually catches the panic, I reverted the fix, rebuilt, and ran the test:

```console
❯ git stash push src/cmd/go/internal/toolchain/select.go
❯ ./make.bash
...
❯ ../bin/go test cmd/go -run=Script/^mod_run_flags_issue77483$
vcs-test.golang.org rerouted to http://127.0.0.1:62389
https://vcs-test.golang.org rerouted to https://127.0.0.1:62390
go test proxy running at GOPROXY=http://127.0.0.1:62391/mod
--- FAIL: TestScript (0.02s)
    --- FAIL: TestScript/mod_run_flags_issue77483 (0.01s)
        script_test.go:139: 2026-02-14T17:14:21Z
        script_test.go:141: $WORK=/var/folders/q0/5_z6pvw574z1c0zcrr2xvk0r0000gn/T/cmd-go-test-2788014133/tmpdir1655598675/mod_run_flags_issue77483614334795
        script_test.go:163: 
            # Regression test for https://go.dev/issue/77483: 'go run -C' should not panic. (0.010s)
            > ! go run -C
            [stderr]
            panic: runtime error: slice bounds out of range [1:0]

...

            > stderr 'flag needs an argument: -C'
        script_test.go:163: FAIL: testdata/script/mod_run_flags_issue77483.txt:4: stderr 'flag needs an argument: -C': no match for `(?m)flag needs an argument: -C` in stderr
        script_test.go:405: go was invoked but no counters were incremented
FAIL
FAIL    cmd/go  0.453s
FAIL
```

The test failed correctly.

## Submitting the fix

Finally, I created a [changelist](https://go.dev/cl/742860), submitted it for review, and it was merged.

{{< figure src="/img/2026-02-14-go-run-panic-fix/go-merged-changelist.webp" width="100%" alt="Screenshot of the merged changelist fixing the panic" >}}

The changelist was accepted without additional comments since the problem and fix were straightforward.

## Conclusion

This experience reinforced a simple idea: if you find bugs in tools you use, fix them.
Stack traces show you where the problem is, tests stop bugs from coming back, and your work helps everyone.
If you find similar bugs in Go, I encourage you to report and fix them too.
