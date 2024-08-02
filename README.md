# prs

✨ Overview
---

`prs` lets you stay updated on the PRs you care about from your terminal.

<p align="center">
  <img src="https://tools.dhruvs.space/images/prs/v0-3-0/prs.gif" alt="Usage" />
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

Or get the binaries directly from a [release][3]. Read more about verifying the
authenticity of released artifacts [here](#verifying-release-artifacts).

🛠️ Pre-requisites
---

- [gh](https://github.com/cli/cli)


🛠️ Configuration
---

Create a configuration file that looks like the following. By default `prs` will
look for this file at `~/.config/prs/prs.yml`.

```yaml
pr-count: 20
query: 'is: pr repo:neovim/neovim is:open label:lua'
```

⚡️ Usage
---

`prs` has four modes:

- "Query" mode (default): lets you search PRs based on a query you provide (based
  on github's [search
  syntax](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests))
- "Repos" mode: lets you see PRs for a repo from a predefined list
- "Author" mode: shows open PRs you've authored
- "Reviewer" mode: shows open PRs that are awaiting a review from you

### Query Mode

```bash
prs
prs --query='type:pr repo:neovim/neovim is:open label:treesitter'
```

Or configure `prs`' YAML configuration file as follows:

```yaml
query: 'is: pr repo:neovim/neovim is:open label:lua'
```

### Repos Mode

Configure `prs`' YAML configuration file as follows:

```yaml
sources:
- owner: owner
  repos:
    - name: repo-a
    - name: repo-b
    - name: repo-c
```

```bash
prs --mode=repos
```

### Author/Reviewer Mode

```bash
prs --mode=author
prs --mode=reviewer
```

### Other settings

`prs`' config file also accepts the following settings:

```yaml
diff-pager: delta
pr-count: 20
```

You can also set up multiple config files, and reference them via the
`--config-file` flag.

**[`^ back to top ^`](#prs)**

Screenshots
---

![Screen 1](https://tools.dhruvs.space/images/prs/v0-3-0/prs-1.png)

![Screen 2](https://tools.dhruvs.space/images/prs/v0-3-0/prs-2.png)

![Screen 3](https://tools.dhruvs.space/images/prs/v0-3-0/prs-3.png)

![Screen 4](https://tools.dhruvs.space/images/prs/v0-3-0/prs-4.png)

![Screen 5](https://tools.dhruvs.space/images/prs/v0-3-0/prs-5.png)

![Screen 6](https://tools.dhruvs.space/images/prs/v0-3-0/prs-6.png)

**[`^ back to top ^`](#prs)**

Keyboard Shortcuts
---

### General

```text
  q/esc/ctrl+c                      go back
  Q                                 quit from anywhere
  ?                                 Open Help View
  d                                 Open PR Details View
  ctrl+v                            Show PR details using gh
```

### PR List View

```text
  Indicators for current review decision:

  ±  implies                        CHANGES_REQUESTED
  🟡 implies                        REVIEW_REQUIRED
  ✅ implies                        APPROVED

  ⏎/tab/shift+tab/2                 Switch focus to PR Timeline View
  ctrl+s                            Switch focus to Repo List View (when --mode=repos)
  ctrl+d                            Show PR diff
  ctrl+r                            Reload PR list
  ctrl+b                            Open PR in browser
```

### PR Details View

```text
  h/N/←                             Go to previous section
  l/n/→                             Go to next section
  1/2/3...                          Go to specific section
  J/]                               Go to next PR
  K/[                               Go to previous PR
  d                                 Go back to last view
  ctrl+b                            Open PR in browser
```

### Timeline List View


```text
  tab/shift+tab/1                   Switch focus to PR List View
  ⏎/3                               Show details for PR timeline item (when applicable)
  ctrl+d                            Show PR diff
  ctrl+b                            Open timeline item in browser
  ctrl+r                            Reload PR timeline
```

### Timeline Item Detail View


```text
  1                                 Switch focus to PR List View
  2                                 Switch focus to PR Timeline List View
  ctrl+d                            Show PR diff
  ctrl+b                            Open timeline item in browser
  h/N/←                             Go to previous section
  l/n/→                             Go to next section
```

### 🔐 Verifying release artifacts

In case you get the `prs` binary directly from a [release][3], you may want to
verify its authenticity. Checksums are applied to all released artifacts, and
the resulting checksum file is signed using
[cosign](https://docs.sigstore.dev/cosign/installation/).

Steps to verify (replace the version in the commands listed with the one you
want):

1. Download the following files from the release:

   - prs_0.3.0_checksums.txt
   - prs_0.3.0_checksums.txt.pem
   - prs_0.3.0_checksums.txt.sig

2. Verify the signature:

   ```shell
   cosign verify-blob prs_0.3.0_checksums.txt \
   --certificate prs_0.3.0_checksums.txt.pem \
   --signature prs_0.3.0_checksums.txt.sig \
   --certificate-identity-regexp 'https://github\.com/dhth/prs/\.github/workflows/.+' \
   --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
   ```

3. Download the compressed archive you want, and validate its checksum:

   ```shell
   curl -sSLO https://github.com/dhth/prs/releases/download/v0.3.0/prs_0.3.0_linux_amd64.tar.gz
   sha256sum --ignore-missing -c prs_0.3.0_checksums.txt
   ```

3. If checksum validation goes through, uncompress the archive:

   ```shell
   tar -xzf prs_0.3.0_linux_amd64.tar.gz
   ./prs
   # profit!
   ```

Acknowledgements
---

`prs` is built using [bubbletea][1], and released via [goreleaser][2].

[1]: https://github.com/charmbracelet/bubbletea
[2]: https://github.com/goreleaser/goreleaser
[3]: https://github.com/dhth/prs/releases

**[`^ back to top ^`](#prs)**
