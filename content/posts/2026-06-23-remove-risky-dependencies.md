---
title: "A risky dependency and how AI helped me remove it"
date: 2026-06-23
tags: ["golang", "security", "opensource", "claude"]
---

A single line in `go.mod` can pull in code from a vendor you would never trust on purpose.
Your project may depend on the russian-maintained `mailru/easyjson` without ever choosing it, buried several levels deep in the dependency graph of popular Go libraries.

This is a real risk, because the dependencies you did not pick are the ones you never check.

In this post, I show how I found one such indirect dependency in a chain of popular libraries, and how a single prompt to Claude helped me remove it.

{{< figure src="/img/2026-06-23-remove-risky-dependencies/oapi_codegen_remove_easyjson.webp" width="100%" alt="Renovate PR on oapi-codegen that removes the mailru/easyjson dependency" >}}

**TL;DR**: `oapi-codegen` depended on `mailru/easyjson` without asking for it directly.
The dependency came through `getkin/kin-openapi`, which used `perimeterx/marshmallow`, which used `easyjson`.
I asked Claude if `marshmallow` was really needed, it was not, and a few PRs later the dependency was gone from the whole chain.

<!--more-->

## The problem: a dependency you can't live without

About a year ago, Hunted Labs published research called
["The Russian open source project that we can't live without"](https://www.huntedlabs.com/research/the-russian-open-source-project-that-we-cant-live-without).
It is about [mailru/easyjson](https://github.com/mailru/easyjson), a fast JSON library made by Mail.Ru (VK), a russian company.
The library is used [everywhere in the Go ecosystem](https://github.com/search?q=path%3A%22go.mod%22+%22github.com%2Fmailru%2Feasyjson%22++NOT+is%3Aarchived+language%3A%22Go+Module%22&type=code&l=Go+Module), often as an indirect dependency that teams never chose on purpose.

{{< figure src="/img/2026-06-23-remove-risky-dependencies/huntedlabs_easyjson.webp" width="100%" alt="Hunted Labs research page titled 'The Russian open source project that we can't live without'" >}}

The concern is not a known bug, but trust and control.
The project is run by a company under russian jurisdiction, so the people with commit access could be pressured to add harmful code.
And because so many projects depend on `easyjson`, one bad release could reach a huge part of the ecosystem.

## Why this matters: every dependency is attack surface

Attacks through dependencies are real, and Go is not safe from them.

In 2025, [Socket researchers found](https://socket.dev/blog/malicious-package-exploits-go-module-proxy-caching-for-persistence) a backdoored copy of [`boltdb/bolt`](https://github.com/boltdb/bolt), a popular Go key-value store, published under a typosquat name to look like the real package.
The [Go module mirror](https://proxy.golang.org/) never changes a version once it is cached, and the attacker abused this: after the bad version was cached, the mirror kept serving it even after the GitHub tag was cleaned.
The backdoor stayed hidden for years.

I have seen this myself.
In a [previous post](/posts/2025-02-28-malicious-go-programs/), I wrote about a group of malicious Go projects that copy real tools and run a trojan the moment you build them.

Both cases teach the same lesson: every dependency you add is attack surface, and a dependency you never chose is the worst kind, because nobody on your team watches it.
You cannot check what you do not know is there.
The safest way to shrink the attack surface is to remove a dependency, because code that is not in your `go.mod` cannot be hacked.

## A little copying is better than a little dependency

There is a good rule for this from Rob Pike, one of the authors of Go, and it is one of the [Go Proverbs](https://go-proverbs.github.io/):

> A little copying is better than a little dependency.

In other words, copying a few lines of code is often safer than adding a whole dependency you must trust and keep up to date.
This is the idea I followed to remove `easyjson`.

## Tracing the dependency chain

We use [oapi-codegen/oapi-codegen](https://github.com/oapi-codegen/oapi-codegen), a popular OpenAPI code generator, in production, so I picked it as the first project to clean up.
There I saw [mailru/easyjson](https://github.com/mailru/easyjson) in its `go.mod`, where it appeared as an [indirect dependency](https://github.com/oapi-codegen/oapi-codegen/blob/1e1b2a251bc24450a5766af9d648e6601436e5aa/go.mod#L22).

```console
❯ git switch -d 1e1b2a251bc24450a5766af9d648e6601436e5aa
HEAD is now at 1e1b2a25 chore(deps): update dessant/label-actions action to v5.0.3 (.github/workflows)

❯ sed -n '16,25p' go.mod
require (
        github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
        github.com/dprotaso/go-yit v0.0.0-20220510233725-9ba8df137936 // indirect
        github.com/go-openapi/jsonpointer v0.23.1 // indirect
        github.com/go-openapi/swag/jsonname v0.26.0 // indirect
        github.com/josharian/intern v1.0.0 // indirect
        github.com/mailru/easyjson v0.9.2 // indirect
        github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
        github.com/oasdiff/yaml v0.1.0 // indirect
        github.com/oasdiff/yaml3 v0.0.13 // indirect
```

The command [`go mod why`](https://pkg.go.dev/cmd/go@go1.26.4#hdr-Explain_why_packages_or_modules_are_needed) shows the exact path that pulls it in:

```console
❯ go mod why -m github.com/mailru/easyjson
# github.com/mailru/easyjson
github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen
github.com/getkin/kin-openapi/openapi3
github.com/perimeterx/marshmallow
github.com/mailru/easyjson/jlexer
```

Read from top to bottom, the chain is:

```text
oapi-codegen
  └── getkin/kin-openapi
        └── perimeterx/marshmallow   (now HumanSecurity/marshmallow)
              └── mailru/easyjson
```

[`getkin/kin-openapi`](https://github.com/getkin/kin-openapi/blob/v0.139.0/go.mod) used `perimeterx/marshmallow`, and `marshmallow` used `easyjson`.

## Fixing it at the source

My first idea was to fix `marshmallow` itself, so I opened [issue #32](https://github.com/HumanSecurity/marshmallow/issues/32).

{{< figure src="/img/2026-06-23-remove-risky-dependencies/marshmallow_issue.webp" width="100%" alt="GitHub issue on the marshmallow asking to drop the easyjson" >}}

The maintainer [agreed to accept a patch](https://github.com/HumanSecurity/marshmallow/issues/32#issuecomment-3513551886), so I sent a [fix in PR #33](https://github.com/HumanSecurity/marshmallow/pull/33).

But the PR was not merged yet, and waiting for a maintainer is slow, so I went one level up the chain instead.
I opened [issue #1192](https://github.com/getkin/kin-openapi/issues/1192) on `kin-openapi` to raise the problem.

## Asking Claude: do we even need this dependency?

Instead of fixing `marshmallow`, I asked if `kin-openapi` needed it at all.
I gave [Claude](https://claude.ai/) a simple, open prompt:

```text
Explore if this project requires using perimeterx/marshmallow or we can use some other more trusted dependency.
```

The answer was clear and useful:

{{< figure src="/img/2026-06-23-remove-risky-dependencies/claude_prompt_marshmallow.webp" width="100%" alt="Claude answering the prompt that marshmallow can be replaced with the encoding/json" >}}

The `Ref` struct has only one JSON field, `$ref`, so two `json.Unmarshal` calls do the same job.
The first call fills the struct, and the second fills a `map[string]any` from which you delete `"$ref"`.
This gives the same result as `marshmallow.Unmarshal(..., WithExcludeKnownFieldsFromMap(true))`, but uses only the standard library.

So I followed Claude's advice and sent [PR #1196](https://github.com/getkin/kin-openapi/pull/1196) to `kin-openapi`, replacing `marshmallow` with `encoding/json`.
A few lines of the standard library replaced a whole external package and broke the chain that pulled in `easyjson`.

{{< figure src="/img/2026-06-23-remove-risky-dependencies/kin_openapi_remove_easyjson.webp" width="100%" alt="My PR #1196 on kin-openapi replacing marshmallow with the standard library encoding/json" >}}

## The fix propagates through the ecosystem

After that, the change spread out on its own, and I never touched `oapi-codegen`:

1. `kin-openapi` merged my PR and released it as [v0.140.0](https://github.com/getkin/kin-openapi/releases/tag/v0.140.0).
2. Renovate then opened [oapi-codegen PR #2395](https://github.com/oapi-codegen/oapi-codegen/pull/2395) on its own to upgrade to v0.140.0.
   This upgrade dropped `easyjson` from the [dependency graph](https://github.com/oapi-codegen/oapi-codegen/blob/5d02a03ac0e92e7c3abe9f185800ea26b3a47fed/go.mod#L16).

```console
❯ git switch -d 5d02a03ac0e92e7c3abe9f185800ea26b3a47fed
HEAD is now at 5d02a03a refactor petstore example (#2277)

❯ go mod why -m github.com/mailru/easyjson
# github.com/mailru/easyjson
(main module does not need module github.com/mailru/easyjson)
```

## Conclusion

You can use AI for more than writing code: you can also use it to ask whether a dependency should exist at all.
One prompt helped me see that `marshmallow` was easy to replace with the standard library.

Fixing the problem upstream helped everyone, not just my project.
Once `kin-openapi` dropped `marshmallow`, the change reached `oapi-codegen` and every project that uses it.

Removing a dependency is the best way to remove its risk, because code that is not in your `go.mod` cannot be hacked, and you do not have to trust its vendor.

---

I help teams clean up their dependencies and use AI to write safer Go code.
See my [GitHub profile](https://github.com/alexandear).
If you want help auditing your `go.mod` or setting up CI/CD for this, [reach out](mailto:oleksandr.red+website@gmail.com?subject=Dependencies%20and%20security).
I would be happy to help.
