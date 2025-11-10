---
title: Upgrading Golangci-lint to v2
date: 2025-11-11
tags: ["go", "lint", "golangci-lint", "opensource"]
draft: true
---

Golangci-lint is the most popular third-party linter runner for Go projects.
It has been around since 2019.
In March 2025, [version 2 was released](https://ldez.github.io/blog/2025/03/23/golangci-lint-v2/), with improved configuration, a new `fmt` command, removal of legacy workarounds, and a nice website.
Adoption of v2 is still low, despite the release being half a year old.
This article shows how to migrate Golangci-lint from v1 to v2.
I have contributed to its development since 2020, so I know many of the corner cases and caveats.

<!-- more -->

## Why migrate from v1 to v2
