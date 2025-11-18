---
title: Upgrading Golangci-lint to v2
date: 2025-11-11
tags: ["go", "lint", "golangci-lint", "opensource"]
draft: true
---

Golangci-lint is the most popular third-party linter runner for Go projects.
It has been around since 2019.
In March 2025, [version 2 was released](https://ldez.github.io/blog/2025/03/23/golangci-lint-v2/).
Adoption of v2 is still low, despite the release being half a year old.
This article shows how to migrate Golangci-lint from v1 to v2.

{{< figure src="/img/2025-11-11-upgrade-golangci-lint-v2/golangci-lint-website-v2-dark.png" width="100%" alt="Screenshot of Golangci-lint v2 website (dark theme)" >}}

<!--more-->

<br/>

[I have contributed](https://github.com/golangci/golangci-lint/commits?author=alexandear) to Golangci-lint development [since 2020](https://github.com/golangci/golangci-lint/graphs/contributors),
reviewed pull requests for v2 changes, and [wrote](https://github.com/golangci/golangci-lint/pull/5439) the migration guide.
Therefore, I know many of the corner cases and caveats.

## Why migrate from v1 to v2

### v1 is deprecated

The latest v1 release is [`v1.64`](https://golangci-lint.run/docs/product/changelog/#v1648) (Feb 12 2025).
It supports Go 1.24 but will not receive further updates or support for newer Go versions.
Development of v2 continues; the latest version is [`v2.6`](https://golangci-lint.run/docs/product/changelog/#v261) with Go 1.25 support.

### New `fmt` command

v2 adds the `fmt` command. It is like `gofmt` but runs multiple formatters:

- gci
- gofmt
- gofumpt
- goimports
- golines
- swaggo

Get the full list with `golangci-lint formatters`.

### Revamped configuration

The v2 configuration is more consistent and logical:

- Changed linter settings (see this [proposal](https://github.com/golangci/golangci-lint/issues/5299)).
- Updated exclusion rules (see this [proposal](https://github.com/golangci/golangci-lint/issues/5298)).
- Changed some default values.
- Removed obsolete flags.

These changes make v1 and v2 configurations incompatible, but migration is straightforward.

### v2 website is better

The new v2 website looks better, is easier to navigate, and has search.
See [this PR](https://github.com/golangci/golangci-lint/pull/5965) for details.
The link structure changed, but old links still work.

| [v1](https://golangci.github.io/legacy-v1-doc/) | [v2](https://golangci-lint.run/) |
|:-----------------------------------------------:|:--------------------------------:|
| {{< figure src="/img/2025-11-11-upgrade-golangci-lint-v2/golangci-lint-website-v1.png" width="100%" alt="Screenshot of Golangci-lint v1 legacy website" >}} | {{< figure src="/img/2025-11-11-upgrade-golangci-lint-v2/golangci-lint-website-v2-light.png" width="100%" alt="Screenshot of Golangci-lint v2 website (light theme)" >}} |

<br/>

<details>
<summary>Screenshot of Golangci-lint v2 website (dark theme)</summary>

<img src="/img/2025-11-11-upgrade-golangci-lint-v2/golangci-lint-website-v2-dark.png"
     alt="Golangci-lint v2 website (dark theme)" width="100%" />
</details>

### Migration to v2 is easy

To help users migrate, Golangci-lint includes the `migration` command and a [migration guide](https://golangci-lint.run/docs/product/migration-guide/) written primarily [by me](https://github.com/golangci/golangci-lint/pull/5439).

## Migration caveats

Although the migration is as simple as running `golangci-lint migrate`, it has some pitfalls.

### Comments in the configuration are not migrated

Due to a technical constraint, the Golangci-lint team can’t adapt the migration command to preserve comments from the v1 configuration.
See this PR for detailed [explanations](https://github.com/golangci/golangci-lint/pull/5506).
You need to manually migrate all comments from the v1 configuration file to v2.

{{< comparison
v1title="[v1 .golangci.yml](https://golangci.github.io/legacy-v1-doc/usage/configuration/#config-file)"
v2title="[v2 .golangci.yml](https://golangci-lint.run/docs/configuration/file/)"
>}}

```yaml
linters:
  # Enable specific linter.
  enable:
    - govet
```
<!-- SPLIT -->
```yaml
version: "2"
linters:
  enable:
    - govet
```

{{< /comparison >}}

### Deprecated options from v1 or unknown fields are not migrated

Deprecated linters and linter settings are removed in v2.

{{< comparison
v1title="[v1 .golangci.yml](https://golangci.github.io/legacy-v1-doc/usage/configuration/#config-file)"
v2title="[v2 .golangci.yml](https://golangci-lint.run/docs/configuration/file/)"
>}}

```yaml
run:
linters:
  enable:
    - cyclop
    - deadcode
linters-settings:
  cyclop:
    skip-tests: true
```
<!-- SPLIT -->
```yaml
version: "2"
linters:
  enable:
    - cyclop
```

{{< /comparison >}}
