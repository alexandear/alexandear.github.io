---
title: How to fix unreadable images in GitHub dark theme
date: 2025-11-10
tags: ["github", "markdown", "opensource"]
---

I use GitHub's dark theme daily.
Sometimes I see README images that are unreadable on the dark theme even though they look fine on the light theme.
For example, the [etcd](https://github.com/etcd-io/etcd/tree/19aa0dbe8fd6317a237bae9b6ea52a4f1b445b19) repo logo:

{{< figure src="/img/2025-11-10-github-unreadable-dark-logo/etcd-unreadable-logo-dark-theme.jpg" width="80%" alt="Unreadable etcd logo on GitHub dark theme" >}}

In this article, I show how to fix this visibility issue using GitHub's Markdown syntax.

<!-- more -->

While the etcd logo is clearly visible on GitHub’s light theme, it’s barely readable on the dark theme.

{{< figure src="/img/2025-11-10-github-unreadable-dark-logo/etcd-logo-white-theme.jpg" width="80%" caption="etcd logo on GitHub light theme" >}}

## Actual Fix

Use the browser’s `prefers-color-scheme` media query with the HTML `<picture>` element so visitors get a light or dark image depending on their GitHub theme.

This solution is mentioned in the [GitHub docs](https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/quickstart-for-writing-on-github#adding-an-image-to-suit-your-visitors):

{{< figure src="/img/2025-11-10-github-unreadable-dark-logo/github-adding-an-image-to-suit-your-visitors.png" width="80%" caption="GitHub docs on adding an image to suit your visitors" >}}

Given this, we can fix the etcd logo by replacing:

```md
![etcd Logo](logos/etcd-horizontal-color.svg)
```

with

```md
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/cncf/artwork/9870640f123303a355611065195c43ac3f27aa19/projects/etcd/horizontal/white/etcd-horizontal-white.png">
  <source media="(prefers-color-scheme: light)" srcset="logos/etcd-horizontal-color.svg">
  <img alt="etcd logo" src="logos/etcd-horizontal-color.svg" width=269 />
</picture>
```

I created a [PR](https://github.com/etcd-io/etcd/pull/18891) to apply this fix, and now the logo renders correctly on GitHub’s dark theme:

{{< figure src="/img/2025-11-10-github-unreadable-dark-logo/etcd-fixed-logo-dark-theme.jpg" width="80%" caption="Fixed etcd logo on GitHub dark theme" >}}

## Repositories with unreadable images

Other repositories with unreadable README images on GitHub's dark theme. Some of them already fixed, but mostly not:

- [containers/.github](https://github.com/containers/.github/tree/145756951b6fe9a25915be833b68a7ff79f5a7de/profile) - [issue](https://github.com/containers/container-libs/issues/70)

    {{< figure src="/img/2025-11-10-github-unreadable-dark-logo/containers-unreadable-logo-dark-theme.png" width="80%" alt="Containers: unreadable logo on GitHub dark theme" >}}

- [reviewdog/reviewdog](https://github.com/reviewdog/reviewdog/tree/54d508bedf6587359eaa38beb523012c30b51c7a) - [issue](https://github.com/reviewdog/reviewdog/issues/1688)

    {{< figure src="/img/2025-11-10-github-unreadable-dark-logo/reviewdog-unreadable-logo-dark-theme.png" width="80%" alt="Reviewdog: unreadable logo on GitHub dark theme" >}}

- [invopop/gobl](https://github.com/invopop/gobl/tree/88769e830d5c6808ea2b710dd4ab4e5ff278aa1b) - [issue](https://github.com/invopop/gobl/issues/642)

    {{< figure src="/img/2025-11-10-github-unreadable-dark-logo/gobl-unreadable-logo-dark-theme.png" width="80%" alt="GOBL: unreadable logo on GitHub dark theme" >}}

- [go-co-op/gocron](https://github.com/go-co-op/gocron/tree/63f3701d571c1bc0c46eea6f7fa238ac16bde3e1) - [PR](https://github.com/go-co-op/gocron/pull/844)

    {{< figure src="/img/2025-11-10-github-unreadable-dark-logo/gocron-unreadable-images-dark-theme.png" width="80%" alt="Go cron: unreadable images on GitHub dark theme" >}}

- [lima-vm/lima](https://github.com/lima-vm/lima/tree/9d4eccb4490920ee62665847161dda740dd7443b) - [PR 2085](https://github.com/lima-vm/lima/pull/2085), [PR 2380](https://github.com/lima-vm/lima/pull/2380)

    {{< figure src="/img/2025-11-10-github-unreadable-dark-logo/lima-unreadable-image-dark-theme.png" width="80%" alt="Lima: unreadable CNCF logo on GitHub dark theme" >}}

- [golangci/golangci-lint](https://github.com/golangci/golangci-lint/tree/1f032fbc4b117e4247b19ff606cc847ab5383bc9) - [PR 5598](https://github.com/golangci/golangci-lint/pull/5598), [PR 5613](https://github.com/golangci/golangci-lint/pull/5613)

    {{< figure src="/img/2025-11-10-github-unreadable-dark-logo/golangci-lint-unreadable-image-dark-theme.png" width="80%" alt="Golangci-lint: unreadable GoLand logo on GitHub dark theme" >}}

- [rqlite/rqlite](https://github.com/rqlite/rqlite/tree/de422a5497d6cde343616287e6a63cdc3c10765a) - [issue](https://github.com/rqlite/rqlite/issues/2023), [commit](https://github.com/rqlite/rqlite/commit/22f20504fde112ca5664752cf2698f54824c3803)

    {{< figure src="/img/2025-11-10-github-unreadable-dark-logo/rqlite-unreadable-logo-dark-theme.png" width="80%" alt="rqlite: unreadable logo on GitHub dark theme" >}}

Feel free to open PRs to fix them!
