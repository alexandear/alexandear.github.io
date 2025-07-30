---
title: "Modernizing Go programs"
date: 2025-07-28
tags: ["go", "opensource"]
draft: true
---

I have been working with Go programs since Go 1.10 (released in Feb 2018) and have vast experience in refactoring programs.
This article is a manual on how you can improve your Go program with the help of some knowledge and tools.
Some of the refactoring is possible with the help of modern AI agents, but as of July 2025, they can't fully do this work.
I show real-world examples from Open Source Go projects which I refactored.

<!--more-->

The article consists of chapters named by Go versions: "Go 1.25", "Go 1.24", "Go 1.23" etc.
Refactorings for Go 1.22 can be only done for this version and above (Go 1.22, 1.23, 1.24 etc.) due [to version compatibility rules](https://go.dev/doc/go1compat).

TODO: explain why do we need refactoring

## Tools that can be used for modernizing

### Golangci-lint

[Golangci-lint](https://golangci-lint.run/) is a linters runner with the ability to fix lint issues.

```sh
golangci-lint run --fix
```

TODO: installation manual

### gopls modernize

TODO: description

```sh
go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...
```

### go fix

TODO: explain

### sed

TODO: explain

### grep

TODO: explain

## Go 1.24


## Go 1.23


## Go 1.22

### Remove redundant loop variables

A construct like `tc := tc` in `for` loops is not needed anymore and we can remove it.
See [Fixing For Loops in Go 1.22](https://go.dev/blog/loopvar-preview) for details.

#### Benefit

Saves one line of code. Avoid weird for Go newbies constructions like `v := v`.

#### Before

```go
	values := []string{"a", "b", "c"}
	for _, v := range values {
		v := v
		go func() {
			fmt.Println(v)
			done <- true
		}()
	}
```

#### After

```go
	values := []string{"a", "b", "c"}
	for _, v := range values {
		go func() {
			fmt.Println(v)
			done <- true
		}()
	}
```

#### Can be fixed with tools


```sh
golangci-lint run --no-config --enable-only copyloopvar --fix ./...
```

```sh
go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -category forvar --fix ./...
```


#### Examples

- [goreleaser/goreleaser](https://github.com/goreleaser/goreleaser/pull/4856/files#diff-3756619488c8c0f0c0300fc0cdcfecbb39c2a7bcb4fe4b3ac5305c6057512986L486)
- [kubernetes-sigs/kueue](https://github.com/kubernetes-sigs/kueue/pull/1946/files#diff-22ad2263a86a607fd28df7741c704614d0f34e208b5270153aa39427e4325fb3L203)
- [IBM/sarama](https://github.com/IBM/sarama/pull/3214/files#diff-cb488ad8239edeaaf8b0c1f469cc15c03fde53cbf22ee996e2f3922b3cc6a0c9L426)
- [google/go-github](https://github.com/google/go-github/pull/3537/files#diff-0f446fb8e4e16b655368f9f1c774d667d5528c9b3103f35481f704e2e33a925fL292)
- [go-critic/go-critic](https://github.com/go-critic/go-critic/pull/1459/files#diff-c2dfb8c940e1232344ce37c2a5942712765d9acf23d43c89345feb81fdbeeb13L43)
- [99designs/gqlgen](https://github.com/99designs/gqlgen/pull/3387/files#diff-fa4826c514673a47321901386ae757f00b2faa73d1433d8dacfc836f4928829aL44)
- [air-verse/air](https://github.com/air-verse/air/pull/682/files#diff-0c22297be1ae696feec687c4dc3d1f425a6ff6c7dfd47d1d2a2275c32d3da14aL96)

### Simplify `for` range loops

"For" loops may now range over integers.
See [For statements with range clause](https://go.dev/ref/spec#For_range) for details.

#### Benefit

Improves readability and less symbols to type.

#### Before

```go
	for i := 0; i < 3; i++ {
		fmt.Println(i)
	}
```

#### After

```go
	for i := range 3 {
		fmt.Println(i)
	}
```

#### Can be fixed with tools


```sh
golangci-lint run --no-config --enable-only intrange --fix ./...
```

```sh
go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -category rangeint --fix ./...
```


#### Examples

- [kubernetes-sigs/kueue](https://github.com/kubernetes-sigs/kueue/pull/5914/files#diff-539f3fc7450aa4c1e6682c00a20c862a4d603225852fdd26bce2fbe6d60ed044R148)
- [lima-vm/lima](https://github.com/lima-vm/lima/pull/3399/files#diff-4fe57274e3aa074c4ccca2967546e5ad77ec58165d477f30560bef494c637e4dR180)
- [mgechev/revive](https://github.com/mgechev/revive/pull/1282/files#diff-75fa8cea7543dbb0e07700624e2760869a23cc2004dcb834e3e5a84739d25519L157)


## Go 1.21


## Go 1.20


## Go 1.19


## Go 1.18


## Go 1.17


## Go 1.16


## Go 1.15


## Go 1.14


## Go 1.13


## Go 1.12


## Go 1.11


## Go 1.10


## Go 1.9


## Go 1.8


## Go 1.7


## Go 1.6


## Go 1.5


## Go 1.4


## Go 1.3


## Go 1.2


## Go 1.1


