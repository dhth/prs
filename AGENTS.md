# AGENTS.md

This file provides guidance to AI coding agents when working with code in this repository.

## Project Overview

`prs` is a Go terminal UI application for browsing GitHub pull requests. It uses the Bubble Tea TUI framework with GitHub GraphQL API integration via `gh` CLI authentication.

## Common Commands

```bash
just build      # go build -ldflags='-s -w' .
just run        # go run .
just install    # go install -ldflags='-s -w' .
just lint       # golangci-lint run
just fmt        # gofumpt -w .
```

Note: always use `just` to run commands.

## Architecture

The app follows the standard **Bubble Tea** (Elm-architecture) pattern: `Model` → `Update(msg)` → `View()` cycle.

**Entry flow**: `main.go` → `cmd.Execute()` (Cobra CLI) → `ui.RenderUI()` → `tea.NewProgram(InitialModel())`

### Key packages

- **`cmd/`** — Cobra CLI setup, config loading (Viper). Config priority: flags > env vars (`PRS_*`) > config file.
- **`ui/`** — All TUI logic, split by concern:
  - `model.go` — Central Bubble Tea `Model` struct (state)
  - `update.go` — Event handling (largest file)
  - `view.go` — View rendering
  - `cmds.go` — Async `tea.Cmd` functions (GitHub API calls)
  - `msgs.go` — Message types passed through Bubble Tea
  - `gh.go` — GitHub GraphQL queries
  - `types.go` — GraphQL query/response type definitions
  - `navigation.go` — View/section navigation with back-stack (`activePane`, `lastPane`, `secondLastActivePane`)
  - `styles.go` / `colors.go` — Lipgloss terminal styling
  - `render_helpers.go` — Display formatting helpers
- **`internal/utils/`** — Markdown rendering with Glamour (embedded Gruvbox theme)

### View system

Six panes managed by `activePane` enum: `repoListView`, `prListView`, `prDetailsView`, `prTLListView`, `prTLItemDetailView`, `helpView`.

PR details has sub-sections: metadata, description, checks, references, files, commits, comments.

### Caching

PR details and timeline data are cached in maps keyed by `"owner/repo:number"` to avoid redundant API calls.

## Linting

Uses golangci-lint v2 with `gofumpt` formatter and `revive` linter (among others). See `.golangci.yml` for full config.
