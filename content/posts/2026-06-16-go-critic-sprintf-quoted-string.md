---
title: "Teaching go-critic a new trick with %#q"
date: 2026-06-16
tags: ["golang", "linter", "opensource", "golangci-lint", "lima", "revive"]
---

Recently, I learned that `%#q` in the `fmt` package wraps a string in backquotes.
Unlike plain `%q`, which adds double quotes `"`, the `#` flag switches to backquotes `` ` ``.

This little discovery led me to [extend](https://github.com/go-critic/go-critic/pull/1529) a [go-critic](https://go-critic.com/overview.html#sprintfquotedstring) check that didn't know about it.

{{< figure src="/img/2026-06-16-go-critic-sprintf-quoted-string/go-critic-pr-1529.webp" width="100%" alt="Screenshot of the PR on go-critic to extend the sprintfQuotedString checker" >}}

<!--more-->

## The discovery

I came across `%#q` while reviewing a [pull request in Lima](https://github.com/lima-vm/lima/pull/5029), where I'm a reviewer.

It turns out this isn't just a cosmetic preference.
The PR solves a real readability problem on Windows, described in [issue #4898](https://github.com/lima-vm/lima/issues/4898).
Windows file paths use backslashes, and `%q` escapes each one, so a path renders with doubled backslashes:

```text
msg="failed to resolve vm for \"C:\\Users\\john\\.lima\\default\\lima.yaml\""
```

With `%#q`, the path is wrapped in backquotes and no escaping is needed, which is far easier to read:

```text
msg="failed to resolve vm for `C:\Users\john\.lima\default\lima.yaml`"
```

That's why the PR replaces `%q` with `%#q` in logs and errors so that paths and identifiers are wrapped in backquotes instead of double quotes:

```go
// before
return fmt.Errorf("failed to read meta data file %q: %w", metaDataPath, err)
// after
return fmt.Errorf("failed to read meta data file %#q: %w", metaDataPath, err)
```

The verbs `%q` and `%#q` exist exactly so you don't have to wrap a string by hand:

- Wrapping a string with `"%s"` can be simplified to `%q`, which also safely escapes special characters.
- Wrapping a string with `` "`%s`" `` can be simplified to `%#q`, which uses raw backtick quoting when the string contains no backtick characters; otherwise it falls back to the `%q` form.

{{< note title="Go string literals" >}}

Double-quoted strings (`"..."`) process escape sequences such as `\\` and `\n`, while backtick strings (`` `...` ``) are raw without escaping.

{{< /note >}}

It's all clearly documented in the [`fmt` package docs](https://pkg.go.dev/fmt); I just had never run into the `#` flag before.

## The gap in go-critic


{{< note title="What is go-critic?" >}}

[go-critic](https://go-critic.com/) is a Go linter that bundles a large collection of checks on style, performance issues, and common errors.
It is supported by [golangci-lint](https://golangci-lint.run/) via the `gocritic` linter.

{{< /note >}}

go-critic has a [`sprintfQuotedString`](https://go-critic.com/overview.html#sprintfquotedstring) check that recommends using `%q` instead of manually wrapping a string with `"%s"`:

```go
fmt.Printf(`"%s"`, v) // suggestion: use %q instead of "%s" for quoted strings
```

But the check didn't detect the backquote case.
This code passed without any warning:

```go
fmt.Printf("`%s`", v) // not detected, but could be simplified to %#q
```

So the linter caught one half of the pattern but not the other.

## Extending the check

I created a [PR that extends the `sprintfQuotedString` rule](https://github.com/go-critic/go-critic/pull/1529) to also flag manually backquoted strings and suggest `%#q`.

Now both forms are covered:

- ``fmt.Sprintf(`"%s"`, s)`` transforms to ``fmt.Sprintf("%q", s)``
- ``fmt.Sprintf("`%s`", s)`` transforms to ``fmt.Sprintf("%#q", s)``

A small change, but it closes a gap and helps everyone use the idiomatic `fmt` verbs.

## The payoff: simpler code

These verbs aren't just shorter, they let you delete code.
While working on [revive](https://github.com/mgechev/revive/pull/1375), I used the same insight to simplify code by replacing a hand-escaped `"%s"` with `%q`:

```go
// before
fmt.Sprintf("Import alias \"%s\" is redundant", imp.Name.Name)
// after
fmt.Sprintf("Import alias %q is redundant", imp.Name.Name)
```

## Conclusion

Reading the docs of tools you use every day pays off.
A flag as small as `#` can simplify your format strings, and a linter is the perfect place to encode that knowledge so the whole community benefits.

---

By the way, I'm a [revive](https://github.com/mgechev/revive) maintainer, a linter expert, and a [golangci-lint contributor](https://github.com/golangci/golangci-lint/commits?author=alexandear).
If you want to write a custom linter, or set up CI/CD that helps AI write better code, [reach out](mailto:oleksandr.red+website@gmail.com?subject=Linters%20and%20CI%2FCD) — I'd be happy to help.
