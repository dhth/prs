# prs

✨ Overview
---

`prs` lets you stay updated on the PRs you care about without leaving the
terminal.

*`prs` is not a replacement of [gh](https://github.com/cli/cli), or the Github
web UI itself, it simply allows you to get to the updates you care about in
fewer key presses.*

<p align="center">
  <img src="https://tools.dhruvs.space/images/prs/prs.gif" alt="Usage" />
</p>

🤔 Motivation
---

For my day job as a tech lead, I need to stay updated on several PRs, and my
hope is that `prs` will let me do that faster than the Github web UI (or other
tools for that matter).

💾 Installation
---

**homebrew**:

```sh
brew install dhth/tap/prs
```

**go**:

```sh
go install github.com/dhth/prs@latest
```

🛠️ Pre-requisites
---

- [gh](https://github.com/cli/cli)


🛠️ Configuration
---

Create a configuration file that looks like the following. By default `prs` will
look for this file at `~/.config/prs/prs.yml`.

```yaml
diff-pager: delta
pr-count: 20
sources:
- owner: owner
  repos:
    - name: repo-a
    - name: repo-b
    - name: repo-c
```

⚡️ Usage
---

```bash
prs
prs -config-file /path/to/config.yml
```

Screenshots
---

![Screen 1](https://tools.dhruvs.space/images/prs/prs-1.png)

![Screen 2](https://tools.dhruvs.space/images/prs/prs-timeline-1.png)

![Screen 3](https://tools.dhruvs.space/images/prs/prs-repos-1.png)

Acknowledgements
---

`prs` is built using [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea
