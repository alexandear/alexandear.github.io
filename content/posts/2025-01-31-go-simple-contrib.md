---
title: How to contribute to Go language
date: 2025-01-31
tags: ["go", "opensource", "git"]
---

In this blog post, I will show that contributing to Go language repositories is as simple as contributing to any other open-source repo.
This will be a step-by-step guide on how to contribute to Go language repositories.
For the example, I will use the main Go repo [go.googlesource.com/go](https://go.googlesource.com/go).

## What we will contribute

I subscribed to [Boldly Go: Daily](https://boldlygo.tech/) for a year and recently got the [following email](https://boldlygo.tech/archive/2025-01-08-determining-the-size-of-a-variable/):

{{< figure src="/img/2025-01-31-go-simple-contrib/boldly-email.png" width="80%" caption="Email from the Boldly Go Subscription" >}}

It says that nobody cares about the grammar mistake in the sentence
`The functions Alignof and Sizeof take an expression x of any type and return the alignment or size, respectively, of a hypothetical variable v as if v was declared via var v = x`
in the [Go language specification](https://tip.golang.org/doc/go1.17_spec#was-declared-via:~:text=variable%20v%20as-,if%20v%20was%20declared,-via%20var%20v).

But I care. Let me fix it.

TL;DR. This [CL](https://go-review.googlesource.com/c/go/+/642037) fixes the grammar nit.

*"was" should be changed to "were" because the sentence is using the [subjunctive mood](https://www.englishoxford.com/2022/04/20/subjunctive-english/), which is used to describe hypothetical or non-real situations. In English, the subjunctive mood often uses "were" instead of "was" for all subjects (I, you, he, she, it, we, they).*

## Contribution steps

There is an official [Contribution Guide](https://go.dev/doc/contribute), but it's too long and boring to read.

*I am assuming you know how to contribute to a project on GitHub.
If not, read [this manual](https://opensource.guide/how-to-contribute/).
I will only tell you about the specifics that apply to most Google software projects.*

### 1. Sign in to Google CLA

First of all, we need to sign the Google Contributor License Agreements with a valid [Google account](https://go.dev/doc/contribute#google_account)
at the [Google CLA](https://cla.developers.google.com/clas).

{{< figure src="/img/2025-01-31-go-simple-contrib/google-cla.png" width="80%" caption="Signed in Google CLA" >}}

### 2. Configure Git authentication

Second, let's configure authentication so we can push our changes to the Go repo.
Log in to [go Git repositories](https://go.googlesource.com/) with a signed Google CLA account and click generate password.

{{< figure src="/img/2025-01-31-go-simple-contrib/go-git-repos.png" width="80%" caption="Home for go Git repositories" >}}

Paste the generated script into a `bash` or `zsh` shell.

{{< figure src="/img/2025-01-31-go-simple-contrib/configure-git.png" width="80%" caption="Configure Git script" >}}

Now we are able to push to any of the repositories in [Go Git](https://go.googlesource.com).

### 3. Register to Gerrit

The communication during a review between people is handled on the [Gerrit](https://www.gerritcodereview.com/) platform,
so [sign in](https://go-review.googlesource.com/login/) with your Google Account.
Gerrit is a little bit different from GitHub Pull Requests and weird at first look.
But it's more powerful and you will definitely like it in the future.

{{< figure src="/img/2025-01-31-go-simple-contrib/gerrit-panel.png" width="80%" caption="Typical Gerrit home page" >}}

### 4. Clone the Go repo

We know that our spec is located in the [go repository](https://go.googlesource.com/go).

{{< figure src="/img/2025-01-31-go-simple-contrib/go-repo-spec.png" width="80%" caption="Go specification in the go git repo" >}}

We should [clone it locally](https://go-review.googlesource.com/admin/repos/go,general) to be able to make changes.

Gerrit requires git commits to have a special line like `Change-Id: If4d3b3965762c8979d304a82493c9eb1068ee13c` present.
We must install the `git-codereview` addon to automatically include this string.

On my machine, it takes around two minutes to fully clone the repo and install `git-codereview` with hooks:

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

Open `doc/go_spec.html` in your favorite editor and finally fix the typo.

{{< figure src="/img/2025-01-31-go-simple-contrib/fix-spec-typo.png" width="80%" caption="Go spec in the VS Code" >}}

Create a commit in the branch `spec-fix-typo` with the following command:

```sh
$ git add . && git codereview change spec-fix-typo
git-codereview: created branch spec-fix-typo tracking origin/master.
git-codereview: change updated.
```

Think of the meaningful commit message.

You must follow git commit message conventions described [here](https://go.dev/doc/contribute#commit_messages).
In short, the template for a good commit message would be:

```txt
prefix: summary of changes

Optional multiline detailed description. There must be a blank line before.
Can be empty.
```

Where `prefix` is the name of the changed file or directory.

In our simple case, a good commit message would be a one-liner `spec: fix grammar issue`.

We can verify the commit with the special line `Change-Id` that gets created:

```sh
$ git log -1
commit c53307c3fdf1126eb6cdb1f09f4f9b83759be705 (HEAD -> spec-fix-typo)
Author: Oleksandr Redko <oleksandr.red+github@gmail.com>
Date:   Fri Jan 10 17:00:24 2025 +0200

    spec: fix grammar issue
    
    Change-Id: If4d3b3965762c8979d304a82493c9eb1068ee13c
```

### 6. Wait for a review

Next, we should push our changes to Gerrit, so someone from the Go team can review them.
Instead of `git push` we use:

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

We can open the link [https://go-review.googlesource.com/c/go/+/642037](https://go-review.googlesource.com/c/go/+/642037) from the output in a browser to see the change.

After some time, usually from a couple of hours to a couple of weeks, someone will review our changes and approve them with `+2`.

{{< figure src="/img/2025-01-31-go-simple-contrib/change-log.png" width="80%" caption="Change log in Gerrit" >}}

You can read more about [the review process](https://go.dev/doc/contribute#review).

Hurray! Now the Gopher Robot can merge it into the `master`.

{{< figure src="/img/2025-01-31-go-simple-contrib/merged-change.png" width="80%" caption="Merged change" >}}

We can see the fixed typo in the wild on the Go website.

[Before](https://tip.golang.org/doc/go1.17_spec#was-declared-via:~:text=variable%20v%20as-,if%20v%20was%20declared,-via%20var%20v):

{{< figure src="/img/2025-01-31-go-simple-contrib/before.png" width="80%" caption="Spec with typo" >}}

[After](https://tip.golang.org/ref/spec#:~:text=variable%20v%20as-,if%20v%20were%20declared,-via%20var%20v):

{{< figure src="/img/2025-01-31-go-simple-contrib/after.png" width="80%" caption="Spec without typo" >}}

## Conclusion

Contributing to the Go repository is straightforward but does have some differences compared to contributing on GitHub.
By following the outlined steps, you can successfully contribute to the Go project.
While the process may seem a bit complex initially, it ensures that contributions are well-managed and maintain high standards.
Everyone can contribute to Go, and it's a simple process once you get the hang of it.
Happy coding!
