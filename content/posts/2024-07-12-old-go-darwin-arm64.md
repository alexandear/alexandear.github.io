---
title: "Fix: no Go 1.15 binary for darwin/arm64 on Apple M1"
date: 2024-07-12T13:52:10+02:00
tags: ["go", "arm64", "macos"]
---

If you use an Apple Silicon (M1) Mac, you will hit an issue when downloading Go 1.15 or earlier via the [`dl`](https://go.googlesource.com/dl) tool.

Here is the solution:

```console
$ GOARCH=amd64; go run golang.org/dl/go1.15@latest download
$ go install golang.org/dl/go1.15@latest
$ go1.15 version
go version go1.15 darwin/amd64
```

<!--more-->

Note: You need a recent Go (1.22+) installed first.

Below are the steps showing how this was found.

### Try to install in the usual way

I have a MacBook Pro with an M1 chip (arm64):

```console
$ uname -m
arm64
```

{{< figure src="/img/2024-07-12-old-go-darwin-arm64/macbook-pro-info.webp" width="40%" caption="About this Mac" >}}

The page “[Managing Go installations](https://go.dev/doc/manage-install)” describes installing multiple Go versions.
The normal steps for Go 1.15:

1. Install the `go1.15` helper into `$GOBIN`:

    ```console
    $ go install golang.org/dl/go1.15@latest
    $ ls "$(go env GOBIN)" | grep go1.15
    go1.15
    ```

2. Download the binary:

    ```console
    $ go1.15 download
    go1.15: download failed: no binary release of go1.15 for darwin/arm64 at https://dl.google.com/go/go1.15.darwin-arm64.tar.gz
    ```

I opened an issue: “[dl: failed to install Go 1.15 on darwin/arm64](https://github.com/golang/go/issues/63626)”.
A maintainer replied it would not be fixed because Go 1.15 is unsupported.
I needed a workaround.

The Go team states that any Go program you can compile for `amd64` should also work on Arm.
So we can use the `amd64` binary.

### Workaround by modifying the dl source

I cloned the [`dl`](https://go.googlesource.com/dl) repo and found that `versionArchiveURL` builds the download URL.

Adding:

```go
if goos == "darwin" && runtime.GOARCH == "arm64" {
if strings.HasPrefix(version, "go1.15") {
		arch = "amd64"
	}
}
```

makes it fall back to the `amd64` archive.

Full adjusted function:

```go
// versionArchiveURL returns the zip or tar.gz URL of the given Go version.
func versionArchiveURL(version string) string {
	goos := getOS()
	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}
	arch := runtime.GOARCH
	if goos == "linux" && runtime.GOARCH == "arm" {
		arch = "armv6l"
	}
	if goos == "darwin" && runtime.GOARCH == "arm64" {
	    if strings.HasPrefix(version, "go1.15") {
			arch = "amd64"
	    }
	}
	return "https://dl.google.com/go/" + version + "." + goos + "-" + arch + ext
}
```

Build and run:

```console
$ cd go1.15
$ go build
$ ./go1.15
go1.15: not downloaded. Run 'go1.15 download' to install ...
$ ./go1.15 download
... downloads ...
$ ./go1.15 version
go version go1.15 darwin/amd64
```

This works, but patching source is cumbersome. A simpler approach exists.

### Better approach: set GOARCH

Try:

```console
$ GOARCH=amd64; go install golang.org/dl/go1.15@latest
go: cannot install cross-compiled binaries when GOBIN is set
```

That fails due to a known issue: “[cmd/go: allow installing cross-compiled binaries when GOBIN is set](https://github.com/golang/go/issues/57485)”.

### Final solution

Use `go run` with `GOARCH` to download, then `go install` normally:

```console
$ GOARCH=amd64; go run golang.org/dl/go1.15@latest download
$ go install golang.org/dl/go1.15@latest
$ go1.15 version
go version go1.15 darwin/amd64
```

Success.

You can repeat this for earlier versions (1.14, 1.13, 1.12, …) that lack native `arm64` binaries:

```sh
GOARCH=amd64; go run golang.org/dl/go1.14@latest download
go install golang.org/dl/go1.14@latest
go1.14 version
```
