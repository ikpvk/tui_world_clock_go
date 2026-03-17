# TUI World Clock - Implementation Plan

## Project Overview

A cross-platform Terminal User Interface (TUI) world clock application built in Go. Displays multiple timezones in a table format with real-time updates.

## Technology Stack

- **Framework**: Bubble Tea (https://github.com/charmbracelet/bubbletea)
- **Styling**: Lipgloss (https://github.com/charmbracelet/lipgloss)
- **Components**: Bubbles (https://github.com/charmbracelet/bubbles) - TextInput for fuzzy search
- **Embedded Timezone Data**: Go `time/tzdata` package for Windows cross-compilation
- **Cross-platform**: Supports Windows, macOS, Linux

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
│       └── timezone.go      # Timezone utilities + country mappings
├── go.mod
├── go.sum
├── PLAN.md                  # This file (developer reference)
└── README.md                # User documentation
```

## Dependencies

```bash
github.com/charmbracelet/bubbles   # TextInput, viewport components
github.com/charmbracelet/bubbletea # TUI framework
github.com/charmbracelet/lipgloss  # Styling
time/tzdata                        # Embedded timezone database (Go stdlib)
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
| `j` / `k` | Alternative navigation (vim-style) |
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
GOOS=linux GOARCH=amd64 go build -tags timetzdata -o worldclock-linux-amd64 ./cmd

# macOS
GOOS=darwin GOARCH=amd64 go build -tags timetzdata -o worldclock-darwin-amd64 ./cmd

# Windows
GOOS=windows GOARCH=amd64 go build -tags timetzdata -o worldclock.exe ./cmd
```

**Note:** The `-tags timetzdata` flag embeds timezone data into the binary, ensuring it works on Windows without requiring system timezone database.

## Known Issues & Solutions

### Windows: "unknown time zone" error

**Problem:** `LoadLocation error for Asia/Kolkata: unknown time zone Asia/Kolkata`

**Solution:** Build with embedded timezone data:
```bash
GOOS=windows GOARCH=amd64 go build -tags timetzdata -o worldclock.exe ./cmd
```

This embeds ~450KB of timezone data into the binary, making it work on any Windows system.

## Performance

- **Binary size**: ~5 MB (includes embedded timezone data)
- **Memory usage**: ~17 MB working set
- **CPU usage**: Minimal (~0%), updates 2x per second

## Implementation Steps Completed

1. Initialize project - `go mod init`, add dependencies
2. Config handling - Read/write JSON config file with cross-platform path
3. Timezone list - Render list of user timezones with current time
4. Real-time updates - Auto-refresh every second
5. Add timezone - Fuzzy search using TextInput + filter
6. Remove timezone - Delete from list and config
7. Styling - Apply lipgloss for clean appearance
8. Country display - Added timezone to country mapping
9. Country search - Search by timezone OR country name
10. Table format - Main screen displays in table format
11. Windows support - Embed timezone data for cross-platform compatibility

## Known Timezones

The app includes ~90 timezone entries with country mappings for major cities worldwide. See `internal/timezone/timezone.go` for the complete list.

## For Future Developers

- To add new timezones: Update `timezoneCountryMap` and `GetAllTimezones()` in `internal/timezone/timezone.go`
- To change refresh interval: Update `time.Second` in `internal/model/model.go` (lines 41, 52, 90)
- To add features: Bubble Tea uses Model-View-Update pattern; main logic is in `internal/model/model.go`
