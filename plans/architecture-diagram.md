# Marvin CLI Architecture - Markdownlint & Dashboard

## System Architecture Overview

```mermaid
graph TB
    subgraph "CLI Commands"
        A[marvin vale]
        B[marvin markdownlint]
        C[marvin dashboard]
    end
    
    subgraph "Dependency Detection"
        D[MultiDetector]
        E[BrewDetector]
        F[NpmDetector]
        G[SystemDetector]
    end
    
    subgraph "Checkers"
        H[ValeChecker]
        I[MarkdownlintChecker]
    end
    
    subgraph "Output Layer"
        J[JSONWriter]
        K[PlainTextFormatter]
    end
    
    subgraph "TUI Layer"
        L[ResultsViewer]
        M[DashboardViewer]
    end
    
    subgraph "Storage"
        N[.marvin/results/]
    end
    
    A --> D
    B --> D
    D --> E
    D --> F
    D --> G
    
    A --> H
    B --> I
    
    H --> J
    I --> J
    
    J --> N
    
    A --> L
    B --> L
    A --> K
    B --> K
    
    C --> N
    N --> M
```

## Command Flow - Markdownlint

```mermaid
sequenceDiagram
    participant User
    participant CLI as markdownlint cmd
    participant Detector as Dependency Detector
    participant Checker as MarkdownlintChecker
    participant Writer as JSON Writer
    participant TUI as TUI Viewer
    participant FS as File System
    
    User->>CLI: marvin markdownlint docs/
    CLI->>Detector: IsInstalled("markdownlint")
    Detector->>Detector: Check brew
    Detector->>Detector: Check npm
    Detector->>Detector: Check system PATH
    Detector-->>CLI: Found at /path/to/markdownlint
    
    CLI->>Checker: NewMarkdownlintChecker()
    CLI->>Checker: Check(docs/)
    Checker->>FS: Execute markdownlint --json docs/
    FS-->>Checker: JSON output
    Checker->>Checker: Parse & transform to Result
    Checker-->>CLI: Result object
    
    CLI->>Writer: Write(result)
    Writer->>FS: Save to .marvin/results/markdownlint-{timestamp}.json
    Writer-->>CLI: Output path
    
    CLI->>TUI: ShowResults(result)
    TUI-->>User: Display interactive TUI
```

## Dashboard Data Flow

```mermaid
sequenceDiagram
    participant User
    participant CLI as dashboard cmd
    participant Agg as Aggregator
    participant FS as File System
    participant TUI as Dashboard TUI
    
    User->>CLI: marvin dashboard
    CLI->>Agg: LoadDashboardData(".marvin/results/")
    Agg->>FS: Scan directory for *.json
    FS-->>Agg: List of result files
    
    loop For each result file
        Agg->>FS: Read JSON file
        FS-->>Agg: Result data
        Agg->>Agg: Parse into Result object
    end
    
    Agg->>Agg: Group by checker
    Agg->>Agg: Calculate summaries
    Agg-->>CLI: DashboardData
    
    CLI->>TUI: ShowDashboard(data)
    TUI-->>User: Display dashboard with tabs
    
    User->>TUI: Navigate between checkers
    TUI->>TUI: Update view
    TUI-->>User: Show selected checker details
```

## File Structure After Implementation

```
cli/
├── cmd/
│   ├── root.go              # Root command (existing)
│   ├── help.go              # Help command (update)
│   ├── vale.go              # Vale command (existing)
│   ├── markdownlint.go      # NEW: Markdownlint command
│   └── dashboard.go         # NEW: Dashboard command
│
├── internal/
│   ├── app/
│   │   ├── checker/
│   │   │   ├── checker.go           # Interface (existing)
│   │   │   ├── vale.go              # Vale implementation (existing)
│   │   │   ├── vale_test.go         # Vale tests (existing)
│   │   │   ├── markdownlint.go      # NEW: Markdownlint implementation
│   │   │   └── markdownlint_test.go # NEW: Markdownlint tests
│   │   │
│   │   ├── dependency/
│   │   │   ├── detector.go          # Multi-detector (update)
│   │   │   ├── brew.go              # Brew detection (existing)
│   │   │   ├── npm.go               # npm detection (existing)
│   │   │   └── system.go            # System PATH (existing)
│   │   │
│   │   ├── dashboard/               # NEW: Dashboard package
│   │   │   ├── aggregator.go        # Result aggregation
│   │   │   └── aggregator_test.go   # Aggregator tests
│   │   │
│   │   ├── output/
│   │   │   ├── writer.go            # JSON writer (existing)
│   │   │   └── formatter.go         # Plain text (existing)
│   │   │
│   │   └── tui/
│   │       ├── viewer.go            # Single result viewer (existing)
│   │       ├── dashboard.go         # NEW: Dashboard viewer
│   │       └── styles.go            # Shared styles (existing)
│   │
│   └── pkg/
│       └── models/
│           ├── result.go            # Result models (existing)
│           └── dashboard.go         # NEW: Dashboard models
│
└── test/
    └── fixtures/
        ├── vale/                    # Vale test data (existing)
        └── markdownlint/            # NEW: Markdownlint test data
            ├── .markdownlint.yaml
            ├── valid.md
            └── invalid.md
```

## TUI State Machine - Dashboard

```mermaid
stateDiagram-v2
    [*] --> Loading
    Loading --> Summary: Data loaded
    
    Summary --> CheckerTab1: Tab key
    Summary --> CheckerTab2: Tab key
    Summary --> CheckerTabN: Tab key
    
    CheckerTab1 --> Summary: Shift+Tab
    CheckerTab2 --> CheckerTab1: Shift+Tab
    CheckerTabN --> CheckerTab2: Shift+Tab
    
    CheckerTab1 --> Details: Enter
    CheckerTab2 --> Details: Enter
    CheckerTabN --> Details: Enter
    
    Details --> CheckerTab1: Esc
    Details --> CheckerTab2: Esc
    Details --> CheckerTabN: Esc
    
    Summary --> [*]: q
    CheckerTab1 --> [*]: q
    CheckerTab2 --> [*]: q
    CheckerTabN --> [*]: q
    Details --> [*]: q
```

## Data Model Relationships

```mermaid
classDiagram
    class Result {
        +string Checker
        +time.Time Timestamp
        +string Path
        +Summary Summary
        +[]Issue Issues
        +map Metadata
    }
    
    class Summary {
        +int TotalFiles
        +int FilesWithIssues
        +int TotalIssues
        +int ErrorCount
        +int WarningCount
        +int InfoCount
    }
    
    class Issue {
        +string File
        +int Line
        +int Column
        +string Severity
        +string Message
        +string Rule
        +string Context
    }
    
    class DashboardData {
        +[]CheckerSummary Checkers
        +int TotalChecks
        +map LatestResults
        +[]Result AllResults
    }
    
    class CheckerSummary {
        +string Name
        +int TotalRuns
        +time.Time LatestRun
        +int TotalIssues
        +int ErrorCount
        +int WarningCount
        +int InfoCount
    }
    
    Result "1" --> "1" Summary
    Result "1" --> "*" Issue
    DashboardData "1" --> "*" CheckerSummary
    DashboardData "1" --> "*" Result
```

## Dependency Detection Flow

```mermaid
flowchart TD
    A[Check for markdownlint] --> B{Try markdownlint-cli2}
    B -->|Found| C[Use markdownlint-cli2]
    B -->|Not found| D{Try markdownlint-cli}
    D -->|Found| E[Use markdownlint-cli]
    D -->|Not found| F{Try markdownlint binary}
    F -->|Found| G[Use markdownlint]
    F -->|Not found| H[Show installation instructions]
    
    C --> I[Return path and version]
    E --> I
    G --> I
    H --> J[Exit with error]
    
    style C fill:#90EE90
    style E fill:#90EE90
    style G fill:#90EE90
    style H fill:#FFB6C1
```

## Integration with Next.js Dashboard

```mermaid
graph LR
    subgraph "CLI"
        A[marvin vale]
        B[marvin markdownlint]
        C[JSON Writer]
    end
    
    subgraph "File System"
        D[.marvin/results/]
    end
    
    subgraph "Next.js Dashboard"
        E[API Routes]
        F[Results Page]
        G[Checker Detail Page]
    end
    
    A --> C
    B --> C
    C --> D
    D --> E
    E --> F
    E --> G
    
    style D fill:#FFE4B5
```

## Checker Interface Implementation

```mermaid
classDiagram
    class Checker {
        <<interface>>
        +Name() string
        +Check(ctx, opts) Result
        +Validate() error
    }
    
    class ValeChecker {
        -configFile string
        -minAlertLevel string
        -valePath string
        -glob string
        +Name() string
        +Check(ctx, opts) Result
        +Validate() error
        -transformResult() Result
    }
    
    class MarkdownlintChecker {
        -configFile string
        -fix bool
        -markdownlintPath string
        +Name() string
        +Check(ctx, opts) Result
        +Validate() error
        -transformResult() Result
    }
    
    Checker <|.. ValeChecker
    Checker <|.. MarkdownlintChecker
```

## Output Directory Structure

```
.marvin/
└── results/
    ├── vale-20260104-090000.json
    ├── vale-20260104-100000.json
    ├── vale-20260104-110000.json
    ├── markdownlint-20260104-090500.json
    ├── markdownlint-20260104-100500.json
    └── markdownlint-20260104-110500.json
```

Each JSON file contains a complete `Result` object:
```json
{
  "checker": "markdownlint",
  "timestamp": "2026-01-04T09:05:00Z",
  "path": "docs/",
  "summary": {
    "total_files": 38,
    "files_with_issues": 12,
    "total_issues": 24,
    "error_count": 5,
    "warning_count": 11,
    "info_count": 8
  },
  "issues": [...],
  "metadata": {
    "config_file": ".markdownlint.yaml",
    "version": "0.32.0"
  }
}
```
