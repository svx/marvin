# CLI

## Tech Stack

- Written in Golang using [Cobra](https://github.com/spf13/cobra), [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [lipgloss](https://github.com/charmbracelet/lipgloss)
- [GoReleaser](https://github.com/goreleaser/goreleaser) for publishing
- Uses [golangci-lint](https://github.com/golangci/golangci-lint) for linting (must pass)

## Design

- Letting the TUI be a pure viewer
- Designing data models with a future frontend (Next.JS) in mind
- We will add more commands and QA checks in the future, the structure must be unified and effortless to extend
- JSON output files must be placed in a dedicated directory
- The CLI must check if the application that is called for checks, for example `vale` or `markdownlint` is installed via `brew` or as part of the project via the package.json file
