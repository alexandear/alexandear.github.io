---
title: "Apple M1 Chip, go1.15: No Binary Release of Go for Darwin/Arm64"
date: 2024-07-12T13:52:10+02:00
tags: ["go", "arm64", "macos"]
---

If you are using an Apple Silicon M1 chip, then you may encounter an issue when downloading Go versions 1.15 or earlier via [`dl`](https://go.googlesource.com/dl).

Here are the solution:

```sh
$ GOARCH=amd64; go run golang.org/dl/go1.15@latest download
$ go install golang.org/dl/go1.15@latest
$ go1.15 version
go version go1.15 darwin/amd64
```

*Note: You need the latest Go 1.22 installed on your machine.*

Below are the steps detailing how I found it.

### Try to Install In the Usual Way

I have a MacBook Pro with an M1 chip and `arm64` architecture:

```sh
$ uname -m
arm64
```

{{< figure src="/img/2024-07-12-old-go-darwin-arm64/macbook-pro-info.png" width="40%" caption="Info about My MacBook Pro" >}}

The topic ["Managing Go installations"](https://go.dev/doc/manage-install) describes how to install multiple versions of Go on the same machine.
The main steps to install Go 1.15 are as follows:

1. Install the `go1.15` program into the `$GOBIN` directory:

```sh
$ go install golang.org/dl/go1.15@latest
```

```sh
$ ls $(go env GOBIN)
...
go1.15
...
```

2. Download the binary:

```sh
$ go1.15 download
go1.15: download failed: no binary release of go1.15 for darwin/arm64 at https://dl.google.com/go/go1.15.darwin-arm64.tar.gz
```

When I first encountered this error, I created an issue for Go team: ["dl: failed to install Go 1.15 on darwin/arm64"](https://github.com/golang/go/issues/63626).
Shortly after that, the maintainer [wrote](https://github.com/golang/go/issues/63626#issuecomment-1770752650) that this won't be fixed as Go 1.15 is no longer supported.
I need to think of the work-around.

The Go team [states](https://go.dev/wiki/GoArm) that any Go program you can compile for x86_64 should work on Arm. So, let's build go1.15 for x86_64 and use it.

### Work-around by Modifying the dl Binary

First, I investigated the `dl` code by cloning it locally from [the repo](https://go.googlesource.com/dl).
I found that the function [`versionArchiveURL`](https://go.googlesource.com/dl/+/889c5db0dd1df202ef86c7d8a7ed78778309b73f/internal/version/version.go#420) is responsible for downloading the archive of the given Go version.

Through trial and error, I found that adding this code leads to a successful result:

```go
if goos == "darwin" && runtime.GOARCH == "arm64" {
	if strings.HasPrefix(version, "go1.15") {
		arch = "amd64"
	}
}
```

The corrected and extended function code:

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

We can build and run our custom binary:

```sh
$ cd go1.15
$ go build
$ ./go1.15
go1.15: not downloaded. Run 'go1.15 download' to install to /Users/Oleksandr_Redko/sdk/go1.15
```

And try to download:

```sh
$ ./go1.15 download
Downloaded   0.0% (    16384 / 122458233 bytes) ...
Downloaded   4.8% (  5914592 / 122458233 bytes) ...
Downloaded  12.0% ( 14696336 / 122458233 bytes) ...
Downloaded  29.1% ( 35684080 / 122458233 bytes) ...
Downloaded  47.6% ( 58326608 / 122458233 bytes) ...
Downloaded  70.7% ( 86605168 / 122458233 bytes) ...
Downloaded  83.9% (102792432 / 122458233 bytes) ...
Downloaded  97.6% (119569520 / 122458233 bytes) ...
Downloaded 100.0% (122458233 / 122458233 bytes)
Unpacking /Users/Oleksandr_Redko/sdk/go1.15/go1.15.darwin-amd64.tar.gz ...
Success. You may now run 'go1.15'
```

We can check the version:

```sh
$ ./go1.15 version
go version go1.15 darwin/amd64
```

Now, we can use the custom binary, `./go1.15`, to compile our old code.

However, modifying the dl's code is too cumbersome. Can we find a better solution? Absolutely!

### Better Approach by Setting GOARCH

We can try setting `GOARCH=amd64`.

```sh
$ GOARCH=amd64; go install golang.org/dl/go1.15@latest
go: cannot install cross-compiled binaries when GOBIN is set
```

Unfortunately, another error occurred. 
This issue, ["cmd/go: allow to install cross-compiled binaries when GOBIN is set"](https://github.com/golang/go/issues/57485)
was already registered on the Go issue tracker in 2022 and seems not to have been resolved yet.

### Final Solution

We know that we can [compile and run](https://pkg.go.dev/cmd/go#hdr-Compile_and_run_Go_program) a Go program in one command.

Let's try:

```sh
$ GOARCH=amd64; go run golang.org/dl/go1.15@latest download

Downloaded   0.0% (    16384 / 122458233 bytes) ...
Downloaded  18.7% ( 22888288 / 122458233 bytes) ...
Downloaded  37.7% ( 46153376 / 122458233 bytes) ...
Downloaded  64.5% ( 78953904 / 122458233 bytes) ...
Downloaded  87.3% (106921184 / 122458233 bytes) ...
Downloaded 100.0% (122458233 / 122458233 bytes)
Unpacking /Users/Oleksandr_Redko/sdk/go1.15/go1.15.darwin-amd64.tar.gz ...
Success. You may now run 'go1.15'
```

It should work:

```sh
$ go1.15 version
zsh: command not found: go1.15
```

The last thing is making available `go1.15` command:

```sh
$ go install golang.org/dl/go1.15@latest
$ go1.15 version
go version go1.15 darwin/amd64
```

Success.

In the same way, we can download other versions: Go 1.14, 1.13, 1.12, and earlier,
for which binary arm64 releases are not available:

```sh
$ GOARCH=amd64; go run golang.org/dl/go1.14@latest download
$ go install golang.org/dl/go1.14@latest
$ go1.14 version
```
