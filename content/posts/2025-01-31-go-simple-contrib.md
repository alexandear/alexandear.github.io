---
title: How to contribute to the Go language
date: 2025-01-31
tags: ["go", "opensource", "git"]
---

In this post, I show that contributing to Go repositories is as simple as contributing to any other open-source repo.
This is a step-by-step guide to contributing to Go repositories.
As an example, I use the main Go repository [go.googlesource.com/go](https://go.googlesource.com/go).

{{< figure src="/img/2025-01-31-go-simple-contrib/merged-change.png" width="80%" alt="Merged change" >}}

<!--more-->

## What we will contribute

I subscribed to [Boldly Go: Daily](https://boldlygo.tech/) for a year and recently received the [following email](https://boldlygo.tech/archive/2025-01-08-determining-the-size-of-a-variable/):

{{< figure src="/img/2025-01-31-go-simple-contrib/boldly-email.png" width="80%" caption="Email from the Boldly Go subscription" >}}

It says nobody cares about the grammatical mistake in the sentence
`The functions Alignof and Sizeof take an expression x of any type and return the alignment or size, respectively, of a hypothetical variable v as if v was declared via var v = x`
in the [Go language specification](https://tip.golang.org/doc/go1.17_spec#was-declared-via:~:text=variable%20v%20as-,if%20v%20was%20declared,-via%20var%20v).

But I care. Let me fix it.

TL;DR: This [CL](https://go-review.googlesource.com/c/go/+/642037) fixes the grammatical error.

*"Was" should be changed to "were" because the sentence uses the [subjunctive mood](https://www.englishoxford.com/2022/04/20/subjunctive-english/), which describes hypothetical or non-real situations. In English, the subjunctive mood often uses "were" instead of "was" for all subjects (I, you, he, she, it, we, they).*

## Contribution steps

There is an official [Contribution Guide](https://go.dev/doc/contribute), but it is a bit long to read.

*I assume you know how to contribute to a project on GitHub.
If not, read [this guide](https://opensource.guide/how-to-contribute/).
I cover only the specifics that apply to most Google software projects.*

### 1. Sign the Google CLA

First, sign the Google Contributor License Agreement with a valid [Google Account](https://go.dev/doc/contribute#google_account)
at the [Google CLA](https://cla.developers.google.com/clas).

{{< figure src="/img/2025-01-31-go-simple-contrib/google-cla.png" width="80%" caption="Signed in to Google CLA" >}}

### 2. Configure Git authentication

Next, configure authentication so you can push your changes to the Go repo.
Log in to the [Go Git repositories](https://go.googlesource.com/) with the Google account you used to sign the CLA and click Generate password.

{{< figure src="/img/2025-01-31-go-simple-contrib/go-git-repos.png" width="80%" caption="Home of Go Git repositories" >}}

Paste the generated script into a `bash` or `zsh` shell.

{{< figure src="/img/2025-01-31-go-simple-contrib/configure-git.png" width="80%" caption="Configure Git script" >}}

Now you can push to any repository in [Go Git](https://go.googlesource.com).

### 3. Register with Gerrit

Code review is handled on the [Gerrit](https://www.gerritcodereview.com/) platform,
so [sign in](https://go-review.googlesource.com/login/) with your Google Account.
Gerrit differs from GitHub pull requests and can look odd at first,
but it is powerful—you may come to like it.

{{< figure src="/img/2025-01-31-go-simple-contrib/gerrit-panel.png" width="80%" caption="Typical Gerrit home page" >}}

### 4. Clone the Go repo

We know the spec is located in the [go repository](https://go.googlesource.com/go).

{{< figure src="/img/2025-01-31-go-simple-contrib/go-repo-spec.png" width="80%" caption="Go specification in the Go repository" >}}

Clone it locally to make changes.

Gerrit requires commits to include a line like `Change-Id: If4d3b3965762c8979d304a82493c9eb1068ee13c`.
Install the `git-codereview` add-on to insert this line automatically.

On my machine, it takes around two minutes to clone the repo and install `git-codereview` with hooks:

```sh
$ git clone https://go.googlesource.com/go && (cd go && go install golang.org/x/review/git-codereview@latest && git-codereview hooks)
Cloning into 'go'...
remote: Sending approximately 431.23 MiB ...
remote: Counting objects: 22, done
remote: Total 639024 (delta 496556), reused 639024 (delta 496556)
Receiving objects: 100% (639024/639024), 431.04 MiB | 6.22 MiB/s, done.
Resolving deltas: 100% (496556/496556), done.
Updating files: 100% (14130/14130), done.
go: downloading golang.org/x/review v1.13.0
```

### 5. Fix a typo

Open `doc/go_spec.html` in your editor and fix the typo.

{{< figure src="/img/2025-01-31-go-simple-contrib/fix-spec-typo.png" width="80%" caption="Go spec in VS Code" >}}

Create a commit on the branch `spec-fix-typo` with the following command:

```sh
$ git add . && git codereview change spec-fix-typo
git-codereview: created branch spec-fix-typo tracking origin/master.
git-codereview: change updated.
```

Write a meaningful commit message.

Follow the Git commit message conventions described [here](https://go.dev/doc/contribute#commit_messages).
In short, a good commit message looks like:

```txt
prefix: summary of changes

Optional multi-line description. Leave a blank line before it.
Can be empty.
```

Here, `prefix` is the name of the changed file or directory.

In our simple case, a good one-liner is `spec: fix grammar issue`.

Verify that the `Change-Id` line was created:

```sh
$ git log -1
commit c53307c3fdf1126eb6cdb1f09f4f9b83759be705 (HEAD -> spec-fix-typo)
Author: Oleksandr Redko <oleksandr.red+github@gmail.com>
Date:   Fri Jan 10 17:00:24 2025 +0200

    spec: fix grammar issue
    
    Change-Id: If4d3b3965762c8979d304a82493c9eb1068ee13c
```

### 6. Wait for review

Next, push your changes to Gerrit so someone from the Go team can review them.
Instead of `git push`, use:

```sh
$ git codereview mail
Counting objects: 5, done.
Delta compression using up to 4 threads.
Compressing objects: 100% (3/3), done.
Writing objects: 100% (3/3), 1.23 KiB | 1.23 MiB/s, done.
Total 3 (delta 2), reused 0 (delta 0)
remote: Resolving deltas: 100% (2/2)
remote: Processing changes: new: 1, done    
remote: 
remote: New Changes:
remote:   https://go-review.googlesource.com/c/go/+/642037 spec: fix grammar issue
remote: 
To https://go.googlesource.com/go
 * [new branch]      HEAD -> refs/for/master
```

Open the link [https://go-review.googlesource.com/c/go/+/642037](https://go-review.googlesource.com/c/go/+/642037) to see the change.

After some time—usually from a few hours to a couple of weeks—someone will review your change and approve it with `+2`.

{{< figure src="/img/2025-01-31-go-simple-contrib/change-log.png" width="80%" caption="Change log in Gerrit" >}}

You can read more about the [review process](https://go.dev/doc/contribute#review).

Hooray! Now the Gopher Robot can merge it into the `master` branch.

{{< figure src="/img/2025-01-31-go-simple-contrib/merged-change.png" width="80%" caption="Merged change" >}}

You can see the fixed typo on the Go website.

[Before](https://tip.golang.org/doc/go1.17_spec#was-declared-via:~:text=variable%20v%20as-,if%20v%20was%20declared,-via%20var%20v):

{{< figure src="/img/2025-01-31-go-simple-contrib/before.png" width="80%" caption="Spec with typo" >}}

[After](https://tip.golang.org/ref/spec#:~:text=variable%20v%20as-,if%20v%20were%20declared,-via%20var%20v):

{{< figure src="/img/2025-01-31-go-simple-contrib/after.png" width="80%" caption="Spec without typo" >}}

## Conclusion

Contributing to the Go repository is straightforward but has some differences compared to contributing on GitHub.
By following the steps above, you can successfully contribute to the Go project.
While the process may seem complex at first, it ensures contributions are well managed and maintain high standards.
Anyone can contribute to Go—it becomes simple once you get the hang of it.
Happy coding!
