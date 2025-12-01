---
title: Upgrading Golangci-lint to v2
date: 2025-12-01
tags: ["go", "lint", "golangci-lint", "opensource"]
draft: true
---

Golangci-lint is the most popular third-party linter runner for Go projects.
It has been around since 2019.
In March 2025, [version 2 was released](https://ldez.github.io/blog/2025/03/23/golangci-lint-v2/).
Adoption of v2 is still low, despite the release being half a year old.
This article shows how to migrate Golangci-lint from v1 to v2.

{{< figure src="/img/2025-12-01-upgrade-golangci-lint-v2/golangci-lint-website-v2-dark.png" width="100%" alt="Screenshot of Golangci-lint v2 website (dark theme)" >}}

<!--more-->

{{< toc >}}

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
| {{< figure src="/img/2025-12-01-upgrade-golangci-lint-v2/golangci-lint-website-v1.png" width="100%" alt="Screenshot of Golangci-lint v1 legacy website" >}} | {{< figure src="/img/2025-12-01-upgrade-golangci-lint-v2/golangci-lint-website-v2-light.png" width="100%" alt="Screenshot of Golangci-lint v2 website (light theme)" >}} |

<br/>

<details>
<summary>Screenshot of Golangci-lint v2 website (dark theme)</summary>

<img src="/img/2025-12-01-upgrade-golangci-lint-v2/golangci-lint-website-v2-dark.png"
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

## Example: upgrading Golangci-lint v1 configuration in the Lima project

As an example of migrating a Golangci-lint v1 configuration, I chose Lima, which is hosted on [GitHub](https://github.com/lima-vm/lima).
[Lima](https://lima-vm.io/docs/) is a tool that launches Linux virtual machines with automatic file sharing and port forwarding.
It's a popular project with 19K stars, written in Go, that uses Golangci-lint for automatic checking.

Every migration of Golangci-lint to v2 consists of the following steps:

1. Install Golangci-lint v2.
2. Run `golangci-lint migrate`.
3. Manually migrate comments from the v1 to v2 configuration file.
4. Run Golangci-lint and deal with new lint issues.
5. Upgrade the Golangci-lint version in CI.

The PR with the Golangci-lint migration in Lima, contributed by me, can be [found here](https://github.com/lima-vm/lima/pull/3330).
Below is a step-by-step guide showing how I did it.

Let's clone the project and switch to commit `0625d0b0` with the Golangci-lint v1 [configuration](https://github.com/lima-vm/lima/blob/0625d0b084450e874869dcbc9f63d4312797c3fe/.golangci.yml):

```console
$ git clone https://github.com/lima-vm/lima.git
$ git switch 0625d0b0 -c chore/migrate-golangci-lint-v2
Switched to a new branch 'chore/migrate-golangci-lint-v2'
```

The Golangci-lint v1 configuration file in Lima has non-default linter configurations, many comments, deprecated linters, and several settings that changed in v2.

<a href="/file/2025-12-01-upgrade-golangci-lint-v2/.golangci.yml-v1-before-migrate.txt" target="_blank" rel="noopener noreferrer">View .golangci.yml (v1) before migration</a>

### Install Golangci-lint v2

The installation manual on the official site is [comprehensive and understandable](https://golangci-lint.run/docs/welcome/install/#local-installation):

```console
$ brew install golangci-lint
$ golangci-lint version
golangci-lint has version 2.6.2 built with go1.25.4 from dc16cf4 on 2025-11-14T02:47:46Z
```

### Run `golangci-lint migrate`

The following command automatically detects the [`.golangci.yml`](https://github.com/lima-vm/lima/blob/0625d0b084450e874869dcbc9f63d4312797c3fe/.golangci.yml) configuration and migrates to v2 in-place:

```console
$ golangci-lint migrate
WARN The configuration comments are not migrated. 
WARN Details about the migration: https://golangci-lint.run/docs/product/migration-guide/ 
WARN The configuration `run.timeout` is ignored. By default, in v2, the timeout is disabled. 
╭───────────────────────────────────────────────────────────────────────────╮
│                                                                           │
│                               We need you!                                │
│                                                                           │
│ Donations help fund the ongoing development and maintenance of this tool. │
│  If golangci-lint has been useful to you, please consider contributing.   │
│                                                                           │
│                  Donate now: https://donate.golangci.org                  │
│                                                                           │
╰───────────────────────────────────────────────────────────────────────────╯
```

<a href="/file/2025-12-01-upgrade-golangci-lint-v2/.golangci.yml-after-migrate.txt" target="_blank" rel="noopener noreferrer">View .golangci.yml (v2) after `golangci-lint migrate`</a>

#### Migration changes

The main difference in the v2 configuration from v1 is `version: "2"` at the beginning:

{{< comparison
v1title="[v1 .golangci.yml](https://github.com/lima-vm/lima/blob/0625d0b084450e874869dcbc9f63d4312797c3fe/.golangci.yml)"
v2title="[v2 .golangci.yml](https://github.com/lima-vm/lima/blob/cad433e25e450dab7cd0396ad6802f8c402ea407/.golangci.yml#L3)"
>}}

```yaml
# no version information means v1
```
<!-- SPLIT -->
```yaml
version: "2"
```

{{< /comparison >}}

The `run.timeout` setting is removed, which means no execution time limit by default:

{{< comparison
v1title="[v1 .golangci.yml](https://github.com/lima-vm/lima/blob/0625d0b084450e874869dcbc9f63d4312797c3fe/.golangci.yml#L21)"
v2title="[v2 .golangci.yml](https://github.com/lima-vm/lima/blob/cad433e25e450dab7cd0396ad6802f8c402ea407/.golangci.yml#L4)"
>}}

```yaml
run:
  timeout: 2m
```
<!-- SPLIT -->
```yaml
run:
  # no timeout
```

{{< /comparison >}}

The next change is the setting for disabling all default linters:

{{< comparison
v1title="[v1 .golangci.yml](https://github.com/lima-vm/lima/blob/0625d0b084450e874869dcbc9f63d4312797c3fe/.golangci.yml#L23)"
v2title="[v2 .golangci.yml](https://github.com/lima-vm/lima/blob/cad433e25e450dab7cd0396ad6802f8c402ea407/.golangci.yml#L7)"
>}}

```yaml
linters:
  disable-all: true
```
<!-- SPLIT -->
```yaml
linters:
  default: none
```

{{< /comparison >}}

Updated linters in the `enable` setting.
The list is sorted alphabetically.
`gofmt`, `gofumpt`, and `goimports` are moved to `formatters`.
`typecheck` is [not a linter](https://golangci-lint.run/docs/product/migration-guide/#typecheck) and was removed.
`gosimple` and `staticcheck` are combined into `staticcheck`.

{{< comparison
v1title="[v1 .golangci.yml](https://github.com/lima-vm/lima/blob/0625d0b084450e874869dcbc9f63d4312797c3fe/.golangci.yml#L25-L79)"
v2title="[v2 .golangci.yml](https://github.com/lima-vm/lima/blob/cad433e25e450dab7cd0396ad6802f8c402ea407/.golangci.yml#L9-L28)"
>}}

```yaml
linters:
  enable:
    - staticcheck
    - gofmt
    - gofumpt
    - gosimple
    - revive
    - goimports
    - govet
    - typecheck
```
<!-- SPLIT -->
```yaml
linters:
  enable:
    - govet
    - revive
    - staticcheck
formatters:
  enable:
    - gofmt
    - gofumpt
    - goimports
```

{{< /comparison >}}

`linters-settings` are moved to `linters.settings`:

{{< comparison
v1title="[v1 .golangci.yml](https://github.com/lima-vm/lima/blob/0625d0b084450e874869dcbc9f63d4312797c3fe/.golangci.yml#L80-L165)"
v2title="[v2 .golangci.yml](https://github.com/lima-vm/lima/blob/cad433e25e450dab7cd0396ad6802f8c402ea407/.golangci.yml#L29-L120)"
>}}

```yaml
linters-settings:
  errorlint:
    asserts: false
```
<!-- SPLIT -->
```yaml
linters:
  settings:
    errorlint:
      asserts: false
```

{{< /comparison >}}

The `issues.exclude-rules` settings are moved to `linters.exclusions.rules`.
The `issues.include` settings are moved to `linters.exclusions.presets`.
Note that in v2, `exclusions.paths` are added that were always excluded by v1.

{{< comparison
v1title="[v1 .golangci.yml](https://github.com/lima-vm/lima/blob/0625d0b084450e874869dcbc9f63d4312797c3fe/.golangci.yml#L171-L189)"
v2title="[v2 .golangci.yml](https://github.com/lima-vm/lima/blob/cad433e25e450dab7cd0396ad6802f8c402ea407/.golangci.yml#L121-L137)"
>}}

```yaml
issues:
  include:
    - EXC0013
    - EXC0014
  exclude-rules:
    - path: "pkg/osutil/"
        text: "uid"
    - path: _test\.go
        linters:
          - godot
    - text: "exported: comment on exported const"
        linters:
          - revive
    - text: "fmt.Sprint.* can be replaced with faster"
        linters:
          - perfsprint
```
<!-- SPLIT -->
```yaml
linters:
  settings:
  exclusions:
    presets:
      - common-false-positives
      - legacy
      - std-error-handling
    rules:
      - path: pkg/osutil/
        text: '(?i)(uid)|(gid)'
      - linters:
          - godot
        path: _test\.go
      - linters:
          - perfsprint
        text: fmt.Sprint.* can be replaced with faster
  paths:
    - third_party$
    - builtin$
    - examples$
```

{{< /comparison >}}

### Manually migrate comments from v1 to v2 configuration file

The old v1 configuration file is kept in `.golangci.bck.yml`, so we can compare changes and add comments to the v2 configuration manually:

```console
$ git status -s
 M .golangci.yml
?? .golangci.bck.yml
```

<a href="/file/2025-12-01-upgrade-golangci-lint-v2/.golangci.yml-with-comments.txt" target="_blank" rel="noopener noreferrer">View .golangci.yml after copying comments from .golangci.bck.yml</a>

### Run Golangci-lint and deal with new lint issues

Run `golangci-lint run`:

```console
$ golangci-lint run > golangci-lint-run-after-migrate.txt
$ tail -6 golangci-lint-run-after-migrate.txt
577 issues:
* noctx: 48
* nolintlint: 1
* perfsprint: 5
* revive: 449
* staticcheck: 74
```

<a href="/file/2025-12-01-upgrade-golangci-lint-v2/golangci-lint-run-after-migrate.txt" target="_blank" rel="noopener noreferrer">View the full `golangci-lint run` log</a>

A lot of issues, and you might feel confused, right?
But it's not so bad. Most of them can be easily excluded and fixed later.

First, enable [`comments`](https://golangci-lint.run/docs/linters/false-positives/#preset-comments) exclusion preset to suppress comment-related issues:

```yaml
  exclusions:
    generated: lax
    presets:
      - comments # <-- this line added
      - common-false-positives
      - legacy
      - std-error-handling
```

This reduces the number of issues from 577 to 72:

```console
$ golangci-lint run
...
72 issues:
* noctx: 48
* nolintlint: 1
* perfsprint: 5
* revive: 5
* staticcheck: 13
```

Next, apply these changes:

- Exclude [`QF`](https://staticcheck.dev/docs/checks/#QF) and [`ST1001`](https://staticcheck.dev/docs/checks/#ST1001) checks from `staticcheck`.
- Exclude new `noctx` issues for `net.Dial`, `net.Listen`, and `exec.Command`.
- Disable the `concat-loop` check for `perfsprint`.
- Allow using `Uid` and `Gid` in `pkg/osutil`.
- Rename `loggerWithoutTs` to `loggerWithoutTS` to satisfy `staticcheck`.
- Disable `staticcheck` for `isColimaWrapper__useThisFunctionOnlyForPrintingHints__` (generated code).
- Remove the `nolint` comment to fix the `nolintlint` issue.

```yaml
linters:
  settings:
    perfsprint:
      int-conversion: false
      err-error: false
      errorf: true
      sprintf1: false
      strconcat: false
      concat-loop: false # <-- this disables concat-loop
    staticcheck:
      checks:
        - all
        - "-SA3000"
        - "-ST1001" # <-- this disables warn about using dot imports
        - "-QF*" # <-- this disables QF checks

  exclusions:
    generated: lax
    rules:
      - linters:
          - noctx
        text: "os/exec.Command must not be called."
      - linters:
          - noctx
        text: "net.* must not be called."
    rules:
      # Allow using Uid, Gid in pkg/osutil.
      - path: pkg/osutil/
        text: '(?i)(uid)|(gid)'
```

These changes eliminate the remaining issues:

```console
$ golangci-lint run
0 issues.
```

Additionally, you can remove [`generated`](https://golangci-lint.run/docs/configuration/file/#linters-configuration) and the default exclusion paths because they are not used in Lima:

```yaml
linters:
  exclusions:
    generated: lax
    paths:
      - third_party$
      - builtin$
      - examples$
```

<a href="/file/2025-12-01-upgrade-golangci-lint-v2/.golangci.yml-final.txt" target="_blank" rel="noopener noreferrer">View the final migrated Golangci-lint configuration</a>

Now you can remove `.golangci.bck.yml`, as it's no longer needed.

### Upgrade the Golangci-lint version in CI

Lima uses the following GitHub Actions workflow [configuration](https://github.com/lima-vm/lima/blob/0625d0b084450e874869dcbc9f63d4312797c3fe/.github/workflows/test.yml#L41-L45) to run Golangci-lint:

```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@55c2c1448f86e01eaae002a5a3a9624417608d84  # v6.5.2
  version: v1.64.2
  args: --verbose --timeout=10m
```

All you need to do is update the golangci-lint-action and the version to the latest:

```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@e7fa5ac41e1cf5b7d48e45e42232ce7ada589601  # v9.1.0
  version: v2.6
  args: --verbose
```

You can also remove the `--timeout` flag since this option is managed via the `.golangci.yml` configuration file.

That's all. You have now migrated Golangci-lint from v1 to v2.

Example PR: https://github.com/lima-vm/lima/pull/3330.

Support the Golangci-lint team: https://donate.golangci.org.
