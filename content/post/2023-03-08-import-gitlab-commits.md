---
title: "Importing GitLab Commit History to GitHub"
date: 2023-03-08
tags: ["go", "gitlab", "github", "git"]
draft: true
---

On an interview sometimes people judge developers by their GitHub statistic. Recently, I saw the tweet:

> "Please don't apply for a Senior dev position if your GitHub looks like this..."

![Twitter Senior with GitHub](/img/2023-03-08-import-gitlab-commits/twitter-senior-github.jpg)

There are a lot of tools to fake GitHub history: [1](https://github.com/Shpota/github-activity-generator),
[2](https://github.com/aljazst/github-contributions-generator), [3](https://github.com/artiebits/fake-git-history).
But they create fake commits and it's cheating. We need to find another way of enriching GitHub activity.

You recall, that your company uses GitLab for day-to-day commits. So, let's export your GitLab statistics to GitHub in anonymized way to comply with NDA.

I will use the tool called [import-gitlab-commits](https://github.com/alexandear/import-gitlab-commits). It's written in Go and very handy to run.

Let's import my company commits from the Clarity GitLab VSC server [clarity.gitlab.com](https://clarity.gitlab.com).

### Step 1: Install import-gitlab-commits

First, let's install [Go](https://go.dev/dl).
Next, run the command
```shell
go install github.com/alexandear/import-gitlab-commits@latest
```

### Step 2: Execute import-gitlab-commits

Set required environment variables:
```shell
export GITLAB_BASE_URL=https://clarity.gitlab.com
export GITLAB_TOKEN=secure_token
export COMMITTER_NAME="Oleksandr Redko"
export COMMITTER_EMAIL=oleksandr.red+github@gmail.com
```
where
- `GITLAB_BASE_URL` is a GitLab [instance URL](https://stackoverflow.com/questions/58236175/what-is-a-gitlab-instance-url-and-how-can-i-get-it).
- `GITLAB_TOKEN` is a personal [access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token). Will be used to fetch your commits from `GITLAB_BASE_URL`;
- `COMMITTER_NAME`, `COMMITTER_EMAIL` are my GitHub name with surname and email.

Run the command:
```
import-gitlab-commits
```

The tool do the following operations:

1. Connects to `GITLAB_BASE_URL` using `GITLAB_TOKEN` and gets my `oredko` user info.
2. Fetches all projects that `oredko` contributed to.
3. For all projects get all commits where author is `oredko`.
4. Creates new repo `repo.clarity.gitlab.com.oredko` on a disk. It adds new commits for all fetched info with the message `Project: GITLAB_PROJECT_ID commit: GITLAB_COMMIT_HASH` and committer `Oleksandr Redko <oleksandr.red+github@gmail.com>`.

### Step 3: Create a new repository on GitHub

The first thing you need to do is create a new repository on GitHub. Log in to your GitHub account and click on the "+" sign in the top right corner of the page. Select "New repository" from the dropdown menu. Give your repository a name, a description, and choose whether you want it to be public or private. Finally, click on the "Create repository" button.

After the import is complete, you should verify that all your GitLab commits have been imported to GitHub. Go to your new repository on GitHub and click on the "Commits" button. You should see all your GitLab commits listed there, along with their metadata and commit messages.

In conclusion, importing your GitLab commits to GitHub is a simple process that requires just a few steps. By following this guide, you can easily migrate your code to GitHub and take advantage of its features and integrations. Good luck with your migration!

### Step 4. Push to GitHub
