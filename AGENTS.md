This file provides guidance to AI coding agents like Claude Code (claude.ai/code), Cursor AI, Codex, Gemini CLI, GitHub Copilot, and other AI coding assistants when working with code in this repository.

# Marvin - Documentation QA Tool

Marvin is a documentation quality assurance tool with three main components: a Go CLI with TUI, a Next.js web dashboard, and VitePress documentation.

## Development Environment

This project uses [Devbox](https://www.jetify.com/docs/devbox) for environment management. Always start by running:

```shell
devbox shell
```

This provides:

- Bun (latest)
- Go 1.25.4
- go-task (latest)

## Common Commands

### Task Runner

The project uses [Task](https://taskfile.dev) as the primary task runner:

```shell
task              # Show available tasks
task docs:serve   # Run documentation locally (from docs/ directory)
```

### Documentation (VitePress)

From the `docs/` directory:

```shell
bun install           # Install dependencies
bun run docs:dev      # Start development server
bun run docs:build    # Build for production
bun run docs:preview  # Preview production build
```

### CLI (Go)

The CLI uses:

- `golangci-lint` for linting (must pass)
- Standard Go testing with `*_test.go` files

### Quality Checks

The project performs documentation QA using:

- [Vale](https://vale.sh) for prose linting
- [markdownlint](https://github.com/DavidAnson/markdownlint) for Markdown linting

Markdownlint configuration (`.markdownlint.yaml`):

- Line length: 120 characters (prose), 80 (headings/code blocks)
- Multiple H1 headings allowed (MD025: false)

## Architecture Overview

### Monorepo Structure

```
marvin/
├── cli/          # Go CLI + TUI
├── web/          # Next.js dashboard (planned)
└── docs/         # VitePress documentation
```

### CLI Architecture (`cli/`)

**Tech Stack:**

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [GoReleaser](https://github.com/goreleaser/goreleaser) - Release automation

**Directory Structure:**

- `/cmd` - Main application commands (one file per command, e.g., `vale.go`, `help.go`)
- `/internal/app` - Application-specific code
- `/internal/pkg` - Shared internal libraries
- `/test` - External test apps and test data
- `/tools` - Supporting tools (can import from `/internal`)

**Design Principles:**

1. **TUI as Pure Viewer**: The TUI displays results but doesn't perform checks
2. **Data Models for Frontend**: Design with Next.js dashboard integration in mind
3. **Extensible Command Structure**: Unified pattern for adding new QA checks
4. **JSON Output**: All check results must be placed in a dedicated directory
5. **Dependency Detection**: CLI must verify that external tools (vale, markdownlint) are installed via brew or npm before running checks
6. **One Command Per File**: Each new command requires a new file under `cmd/`

### Web Dashboard (`web/`)

Next.js-based dashboard (in development) where users can:
- View documentation QA results
- Run checks through a web interface
- Monitor documentation quality metrics

### Documentation (`docs/`)

**Tech Stack:**
- [VitePress](https://vitepress.dev) - Static site generator
- [vitepress-plugin-llms](https://github.com/okineadev/vitepress-plugin-llms) - LLM-friendly documentation

**Structure:**

- `.vitepress/` - Configuration and theme
- `contribute/` - Contribution guides (cli, web, docs)
- `reference/` - API and CLI reference documentation

## Key Constraints

1. **CLI Testing**: Unit tests must be in `*_test.go` files in the same package as the code being tested
2. **Go Conventions**: Files/directories starting with `.` or `_` are ignored by Go
3. **Linting**: All Go code must pass `golangci-lint` checks
4. **Future Extensibility**: The structure must support adding more QA checks and commands easily
