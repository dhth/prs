# prs

✨ Overview
---

`prs` lets you stay updated on the PRs you care about without leaving the
terminal.

<p align="center">
  <img src="./static/prs.gif?raw=true" alt="Usage" />
</p>


💾 Installation
---

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

Acknowledgements
---

`prs` is built using [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea
