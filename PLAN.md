# TUI World Clock - Implementation Plan

## Project Overview

A cross-platform Terminal User Interface (TUI) world clock application built in Go. Displays multiple timezones in a table format with real-time updates.

## Technology Stack

- **Framework**: Bubble Tea (https://github.com/charmbracelet/bubbletea)
- **Styling**: Lipgloss (https://github.com/charmbracelet/lipgloss)
- **Components**: Bubbles (https://github.com/charmbracelet/bubbles) - TextInput for fuzzy search
- **Cross-platform**: Pure Go (no CGO) - supports Windows, macOS, Linux

## Project Structure

```
worldclock/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── model/
│   │   └── model.go         # Bubble Tea model (state, Update, View)
│   ├── config/
│   │   └── config.go        # Load/save user timezones
│   └── timezone/
│       └── timezone.go      # Timezone utilities
├── go.mod
├── go.sum
└── PLAN.md
```

## Features

| Feature | Implementation |
|---------|----------------|
| Display timezone table | Table format with Timezone, Country, Date/Time columns |
| Real-time updates | Auto-refresh every second |
| Add timezone | Press `a` to open fuzzy search input |
| Remove timezone | Press `d` or `Delete` when selected |
| Quit | Press `q` or `Ctrl+C` |
| Country display | Shows timezone with country name (e.g., "America/New_York (United States)") |
| Country search | Can search by country name or timezone name |

## User Interactions

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate timezone list |
| `a` | Add new timezone (opens search) |
| `d` / `Delete` | Remove selected timezone |
| `Esc` | Cancel add / go back to list |
| `Enter` | Confirm selection in add mode |
| `q` | Quit application |

## Search Screen (AddView)

- Shows all available timezones with country names
- Fuzzy search filters by timezone name OR country name
- Last item selected by default when entering search
- Arrow keys navigate through filtered results

## Table Format

Main screen displays in table format:

```
+-------------------------+------------------+--------------------+
| Timezone                | Country          | Date/Time         |
+-------------------------+------------------+--------------------+
| America/New_York        | United States    | 2024-01-15 15:45 |
| Europe/London           | United Kingdom   | 2024-01-15 20:45 |
| Asia/Tokyo              | Japan            | 2024-01-16 04:45 |
+-------------------------+------------------+--------------------+
```

Column widths:
- Timezone: 25 chars
- Country: 18 chars
- Date/Time: 20 chars

## Data Model

```go
// Config stored in ~/.config/worldclock/config.json
type Config struct {
    Timezones []string  // e.g., ["America/New_York", "Asia/Tokyo", "Europe/London"]
}

// Bubble Tea Model
type Model struct {
    timezones []string
    selected  int
    view      ViewMode  // "list" | "add"
    input     textinput.Model
    ticker    *time.Ticker
    done      bool
}

type ViewMode string
const (
    ListView ViewMode = "list"
    AddView  ViewMode = "add"
)
```

## Time Format

- 24-hour format
- Date shown alongside time: `YYYY-MM-DD HH:MM:SS`

## Config Location

- Linux/macOS: `~/.config/worldclock/config.json`
- Windows: `%APPDATA%\worldclock\config.json`

## Build Commands

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o worldclock-linux-amd64 ./cmd
GOOS=linux GOARCH=arm64 go build -o worldclock-linux-arm64 ./cmd

# macOS
GOOS=darwin GOARCH=amd64 go build -o worldclock-darwin-amd64 ./cmd
GOOS=darwin GOARCH=arm64 go build -o worldclock-darwin-arm64 ./cmd

# Windows
GOOS=windows GOARCH=amd64 go build -o worldclock.exe ./cmd
```

## Implementation Steps Completed

1. ✅ Initialize project - `go mod init`, add dependencies
2. ✅ Config handling - Read/write JSON config file with cross-platform path
3. ✅ Timezone list - Render list of user timezones with current time
4. ✅ Real-time updates - Auto-refresh every second
5. ✅ Add timezone - Fuzzy search using TextInput + filter
6. ✅ Remove timezone - Delete from list and config
7. ✅ Styling - Apply lipgloss for clean appearance
8. ✅ Country display - Added timezone → country mapping
9. ✅ Country search - Search by timezone OR country name
10. ✅ Table format - Main screen displays in table format
11. ✅ Bug fixes - Various navigation and input handling fixes

## Known Timezones

The app includes ~100 timezone entries with country mappings for major cities worldwide.
