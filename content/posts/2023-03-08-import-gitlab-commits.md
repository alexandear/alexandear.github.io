---
title: "Importing commit history from GitLab to GitHub"
date: 2023-03-08
tags: ["gitlab", "github", "git", "go", "opensource"]
bigimg: [{src: "/img/2023-03-08-import-gitlab-commits/contribs-after-run.png", desc: "GitHub Contribution Graph"}]
---

In interviews, people judge developers by their GitHub. Recently, I saw a tweet with a picture showing GitHub contributions
with one commit of activity and the caption:

> "Please don't apply for a Senior dev position if your GitHub looks like this..."

{{< figure src="/img/2023-03-08-import-gitlab-commits/twitter-senior-github.jpg" width="50%" alt="Tweet Senior GitHub" >}}

In this blog post, I will answer the question of how to enrich GitHub statistics to improve your job prospects.

<!--more-->

There are many tools to fake GitHub history: [1](https://github.com/Shpota/github-activity-generator),
[2](https://github.com/aljazst/github-contributions-generator), [3](https://github.com/artiebits/fake-git-history).
But they create unreal commits, which is cheating. We need to find another way to enrich GitHub activity.

## Resolution

You may remember that your company uses GitLab for day-to-day commits. Let's export these commits to GitHub.

I use the tool [`import-gitlab-commits`](https://github.com/alexandear/import-gitlab-commits).
It's written in Go, very handy to run, and exports commits in an anonymized way to comply with NDA.

I will import my Clarity commits from their internal GitLab VCS [clarity.gitlab.com](https://clarity.gitlab.com).
My GitHub statistics for 2020 before running `import-gitlab-commits` look like:

{{< figure src="/img/2023-03-08-import-gitlab-commits/contribs-before-run.png" width="100%" caption="GitHub Before import-gitlab-commits" >}}

### Step 1: Install import-gitlab-commits

First, let's install [Go](https://go.dev/dl) and run the command:

```sh
go install github.com/alexandear/import-gitlab-commits@latest
```

### Step 2: Execute import-gitlab-commits

Next, set the required environment variables:

```sh
export GITLAB_BASE_URL=https://clarity.gitlab.com
export GITLAB_TOKEN=your_secure_token
export COMMITTER_NAME="Oleksandr Redko"
export COMMITTER_EMAIL=oleksandr.red+github@gmail.com
```

where:
- `GITLAB_BASE_URL` is a GitLab [instance URL](https://stackoverflow.com/questions/58236175/what-is-a-gitlab-instance-url-and-how-can-i-get-it).
- `GITLAB_TOKEN` is a personal [access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token).
It will be used to fetch your commits from `GITLAB_BASE_URL`.
- `COMMITTER_NAME`, `COMMITTER_EMAIL` are my GitHub name and email.

And run the command:

```sh
import-gitlab-commits
```

The tool will perform the following operations:

1. Connect to `GITLAB_BASE_URL` using `GITLAB_TOKEN` and get my `oredko` user info.
2. Fetch all GitLab projects that `oredko` contributed to.
3. For all projects, retrieve the commits where the author is `oredko`.
4. Create a new repo `repo.clarity.gitlab.com.oredko` on disk. Add new commits for all fetched info with
the message `Project: <PROJECT_ID> commit: <COMMIT_HASH>` and committer `Oleksandr Redko <oleksandr.red+github@gmail.com>`:

{{< figure src="/img/2023-03-08-import-gitlab-commits/anonymized-commits.png" width="50%" caption="Git Commit Anonymized Log" >}}

### Step 3: Create a GitHub repository and push to it

Finally, follow [the guide](https://docs.github.com/en/get-started/quickstart/create-a-repo)
and create a new GitHub repository called `clarity-contributions`.

Open the repo created by `import-gitlab-commits` and push to GitHub:

```sh
cd repo.clarity.gitlab.com.oredko
git remote add origin git@github.com:alexandear/clarity-contributions.git
git push
```

That's it. My empty GitHub contribution graph from 2020 became full of commits:

{{< figure src="/img/2023-03-08-import-gitlab-commits/contribs-after-run.png" width="100%" caption="GitHub After import-gitlab-commits" >}}

## Summary

In this article, I suggest a way to enrich GitHub activity by exporting real GitLab statistics.
The [`import-gitlab-commits`](https://github.com/alexandear/import-gitlab-commits) tool is a good solution for this purpose.
