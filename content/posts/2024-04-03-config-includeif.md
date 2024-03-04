---
title: Streamlining Project Management in Git Using IncludeIf
date: 2024-03-04T15:05:32+02:00
tags: ["git"]
---

I work simultaneously on different companies' projects, open-source projects, and personal projects. I should therefore use the same username but distinct email addresses for each of them. Although we can set a `git config user.email <email@example.com>` for each project repository, this might be a tedious manual task and perhaps unsuitable.

The Git config manual describes the use of configuration files to better manage such a task through  ["Conditional includes"](https://git-scm.com/docs/git-config#_conditional_includes).

# Situation

All my projects, both work and non-work, are located under the the `~/src/` directory.
The directories associated with each project, and the emails that I wish to use for them, are as follows:

1. Personal: `~/src/github.com/alexandear/` with the email alexandear@users.noreply.github.com.
2. GitHub Open Source: `~/src/github.com/` with the email oleksandr.red+github@gmail.com.
4. EPAM work projects: `~/src/github.com/epam/` with the email oredko@epam.com.
4. GitLab Open Source: `~/src/gitlab.com/` with the email oleksandr.red+gitlab@gmail.com.

# Resolution

First, set a global username and email in the `~/.gitconfig` file. These will be used by default:

```ini
[user]
    name = Oleksandr Redko
    email = oleksandr.red@gmail.com
```

Second, create different git config files for the different types of work:

- `~/.gitconfig-alexandear` file:

    ```ini
    [user]
        email = alexandear@users.noreply.github.com
    ```

- `~/.gitconfig-github` file:

    ```ini
    [user]
        email = oleksandr.red+github@gmail.com
    ```

- `~/.gitconfig-epam` file:

    ```ini
    [user]
        email = oredko@epam.com
    ```

- `~/.gitconfig-gitlab` file:

    ```ini
    [user]
        email = oleksandr.red+gitlab@gmail.com
    ```

Third, set conditional includes in the `~/.gitconfig` by providing path patterns:

```ini
[includeIf "gitdir:~/src/github.com/"]
    path = ~/.gitconfig-github
[includeIf "gitdir:~/src/github.com/alexandear/**/.git"]
    path = ~/.gitconfig-alexandear
[includeIf "gitdir:~/src/github.com/epam/**/.git"]
    path = ~/.gitconfig-epam
[includeIf "gitdir:~/src/gitlab.com/"]
    path = ~/.gitconfig-gitlab
```

Explanation:

- `"gitdir:~/src/github.com/"` matches all projects under `~/src/github.com/` directories
with specific excaptions: `~/src/github.com/alexandear/**` and `~/src/github.com/epam/**`.
These are excluded because `"gitdir:~/src/github.com/alexandear/**/.git"` and `"gitdir:~/src/github.com/epam/**/.git"`
are placed after `"gitdir:~/src/github.com/"`.
- `"gitdir:~/src/gitlab.com/"` matches all projects located within `~/src/gitlab.com/` directories.

Now, let's verify if the correct emails are set:

1. For personal projects under the `github.com/alexandear` directory:
    ```sh
    ❯ cd ~/src/github.com/alexandear/import-gitlab-commits
    ❯ git config user.email
    alexandear@users.noreply.github.com
    ```

2. For open source GitHub projects:
    ```sh
    ❯ cd ~/github.com/golangci/golangci-lint
    ❯ git config user.email
    oleksandr.red+github@gmail.com
    ```

3. For EPAM work projects under the `github.com/epam` directory:
    ```sh
    ❯ cd ~/src/github.com/epam/OSCI
    ❯ git config user.email
    oredko@epam.com
    ```

4. For open source GitLab projects:
    ```sh
    ❯ cd ~/src/gitlab.com/gitlab-org/gitlab
    ❯ git config user.email
    oleksandr.red+gitlab@gmail.com
    ```

## Conclusion

In conclusion, using `includeIf` and setting up related `.gitconfig` files allows you
to specify your preferred settings for each project or company individually.
With this setup, you can enjoy hassle-free project management in Git.
