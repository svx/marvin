---
title: Command Line Interface Reference
description: Complete reference for Marvin's CLI commands, flags, and exit codes
permalink: /reference/cli/
outline: deep
---

# Command Line Interface Reference

Marvin provides a command-line interface for running documentation quality assurance checks. The CLI features an interactive TUI (Terminal User Interface) for viewing results and supports multiple output formats.

## Installation

```bash
# Install from source (requires Go 1.21+)
git clone https://github.com/svx/marvin
cd marvin/cli
go build -o marvin .

# Move to your PATH
sudo mv marvin /usr/local/bin/
```

## Global Usage

```bash
marvin [command] [flags]
```

### Global Flags

All commands support these global flags:

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--output-dir` | string | `.marvin/results/` | Output directory for JSON results |
| `--no-tui` | boolean | `false` | Disable TUI, output plain text to stdout |
| `--json` | boolean | `false` | Output raw JSON to stdout (implies `--no-tui`) |
| `--verbose` | boolean | `false` | Enable verbose logging |
| `--config` | string | `.marvin.yaml` | Path to config file |
| `-h, --help` | boolean | `false` | Display help information |
| `-v, --version` | boolean | `false` | Display version information |

## Commands

### `help` - Display Help Information

Displays help information for Marvin or a specific command.

#### Usage

```bash
marvin help [command]
```

#### Examples

```bash
# Show general help
marvin help

# Show help for a specific command
marvin help vale
```

#### Output

The help command displays:

- Command description and usage
- Available subcommands
- Flags and their descriptions
- Usage examples
- Additional resources

The output is styled using lipgloss for better readability in the terminal.

---

### `vale` - Run Vale Prose Linting

Runs Vale prose linting on documentation files to check for style guide violations, grammar issues, and other prose problems.

#### Usage

```bash
marvin vale [path] [flags]
```

#### Arguments

| Argument | Required | Default | Description |
|----------|----------|---------|-------------|
| `path` | No | `docs/` | Path to scan for documentation files |

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--config` | string | auto-detect | Vale config file path (`.vale.ini`) |
| `--min-alert-level` | string | `suggestion` | Minimum alert level: `suggestion`, `warning`, or `error` |
| `--glob` | string | - | Glob pattern to filter files (for example, `!node_modules` or `!{dir1/*,dir2/*}`) |

Plus all [global flags](#global-flags).

#### Examples

```bash
# Scan default docs/ directory with TUI
marvin vale

# Scan specific directory
marvin vale ./content

# Use custom Vale config
marvin vale --config .vale.ini

# Output JSON only
marvin vale --json

# Disable TUI, show plain text
marvin vale --no-tui

# Scan with verbose logging
marvin vale --verbose

# Only show errors and warnings
marvin vale --min-alert-level warning

# Ignore node_modules directory
marvin vale --glob='!node_modules'

# Ignore multiple directories
marvin vale --glob='!{node_modules/*,.vitepress/*,dist/*}'

# Ignore specific files and directories
marvin vale --glob='!{README.md,legal/*,meetings/*}'
```

#### Behavior

1. **Dependency Check**: Verifies Vale is installed via Homebrew, npm, or system PATH
2. **Execution**: Runs Vale with `--output=JSON` flag
3. **Result Storage**: Saves JSON results to `.marvin/results/vale-{timestamp}.json`
4. **Display**: Shows results in TUI (default), plain text, or JSON format

#### Output Formats

##### TUI (Default)

Interactive terminal interface with:
- Summary statistics (files scanned, issues by severity)
- Scrollable issue list with file locations
- Color-coded severity levels
- Keyboard navigation (q to quit, arrows to scroll)

Example:
```
 Marvin - Vale Results 

Summary
  Path: docs/
  Files Scanned: 42
  Files with Issues: 8
  Total Issues: 23 (5 errors, 12 warnings, 6 suggestions)

Issues
  docs/getting-started.md:12:5
  [error] Vale.Spelling
  Did you really mean 'installtion'?

  docs/api-reference.md:45:10
  [warning] Vale.Terms
  Use 'API' instead of 'api'

Press q to quit
```

##### Plain Text (`--no-tui`)

Formatted text output suitable for logs or non-interactive environments:

```
Marvin - vale Results
═══════════════════════════════════════════════════════════

Summary:
  Path: docs/
  Files Scanned: 42
  Files with Issues: 8
  Total Issues: 23 (5 errors, 12 warnings, 6 suggestions)

Issues:
───────────────────────────────────────────────────────────

docs/getting-started.md:12:5
[error] Vale.Spelling
Did you really mean 'installtion'?

docs/api-reference.md:45:10
[warning] Vale.Terms
Use 'API' instead of 'api'

Results saved to: .marvin/results/vale-20260102-130000.json
```

##### JSON (`--json`)

Machine-readable JSON output for integration with other tools:

```json
{
  "checker": "vale",
  "timestamp": "2026-01-02T13:00:00Z",
  "path": "docs/",
  "summary": {
    "total_files": 42,
    "files_with_issues": 8,
    "total_issues": 23,
    "error_count": 5,
    "warning_count": 12,
    "info_count": 6
  },
  "issues": [
    {
      "file": "docs/getting-started.md",
      "line": 12,
      "column": 5,
      "severity": "error",
      "message": "Did you really mean 'installtion'?",
      "rule": "Vale.Spelling",
      "context": "installtion"
    }
  ],
  "metadata": {
    "config_file": "",
    "min_alert_level": "suggestion"
  }
}
```

#### Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success - no errors found |
| `1` | Errors found in documentation |
| `1` | Vale not installed or execution failed |

#### Error Handling

If Vale is not installed, Marvin displays installation instructions:

```
Error: vale is not installed

Marvin requires vale to run this check.

Installation options:

  Homebrew (recommended):
    brew install vale

  npm:
    npm install -g vale

  Manual:
    https://vale.sh/docs/vale-cli/installation/

After installation, run this command again.
```

#### JSON Output Files

Results are automatically saved to `.marvin/results/` with the naming pattern:
```
vale-{timestamp}.json
```

Example: `vale-20260102-130405.json`

You can customize the output directory:
```bash
marvin vale --output-dir ./qa-results
```

## Configuration File

Marvin supports a configuration file (`.marvin.yaml`) in the project root:

```yaml
# Default output directory
output_dir: .marvin/results/

# Default scan paths for each checker
defaults:
  vale:
    path: docs/
    config: .vale.ini
    min_alert_level: suggestion

# TUI settings
tui:
  enabled: true
  theme: default

# Dependency detection
dependencies:
  check_brew: true
  check_npm: true
  check_system: true
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Documentation QA

on: [push, pull_request]

jobs:
  vale:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Vale
        run: |
          brew install vale
      
      - name: Install Marvin
        run: |
          cd cli
          go build -o marvin .
      
      - name: Run Vale Check
        run: |
          ./cli/marvin vale --no-tui
```

### GitLab CI

```yaml
documentation-qa:
  image: golang:1.21
  script:
    - apt-get update && apt-get install -y vale
    - cd cli && go build -o marvin .
    - ./marvin vale --no-tui
  only:
    - merge_requests
    - main
```

## Troubleshooting

### Vale Not Found

**Problem**: `Error: vale is not installed`

**Solution**: Install Vale using one of these methods:
- Homebrew: `brew install vale`
- npm: `npm install -g vale`
- Manual: See [Vale installation docs](https://vale.sh/docs/vale-cli/installation/)

### Path Does Not Exist

**Problem**: `Error: path does not exist: docs/`

**Solution**: Specify the correct path to your documentation:
```bash
marvin vale ./path/to/docs
```

### Permission Denied

**Problem**: Cannot write to output directory

**Solution**: Ensure you have write permissions or specify a different output directory:
```bash
marvin vale --output-dir ~/marvin-results
```

## Future Commands

The following commands are planned for future releases:

- `marvin markdownlint` - Markdown linting
- `marvin linkcheck` - Broken link detection
- `marvin spellcheck` - Spell checking
- `marvin all` - Run all checks sequentially

## See Also

- [Vale Documentation](https://vale.sh/docs/)
- [Contributing to Marvin CLI](../contribute/cli/)
- [Marvin Architecture](../../cli/README.md)
