# World Clock

A simple terminal-based world clock that shows multiple timezones in a table format. Built with Go.

## Screenshots

### Main View (List)
```
+-------------------------+------------------+--------------------+
| Timezone                | Country          | Date/Time         |
+-------------------------+------------------+--------------------+
| America/New_York        | United States    | 2026-03-17 08:30 |
| Europe/London           | United Kingdom   | 2026-03-17 13:00 |
| Asia/Tokyo              | Japan            | 2026-03-17 21:30 |
+-------------------------+------------------+--------------------+

↑↓ Navigate | a: Add | d: Delete | q: Quit
```

### Add Timezone (Search)
```
Add Timezone

  Africa/Cairo (Egypt)
  Africa/Johannesburg (South Africa)
> America/Los_Angeles (United States)
  America/New_York (United States)

Search: los_
↑↓ Navigate | Enter: Select | Esc: Cancel
```

## Features

- View multiple timezones side by side
- Real-time updates (refreshes every second)
- Add timezones with fuzzy search
- Remove timezones you don't need
- Search by city or country name

## Quick Start

```bash
# Clone and run
git clone https://github.com/ikpvk/tui_world_clock_go.git
cd tui_world_clock_go
go run ./cmd
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate the list |
| `j` / `k` | Alternative navigation |
| `a` | Add a new timezone |
| `d` or `Delete` | Remove selected timezone |
| `Enter` | Confirm selection |
| `Esc` | Cancel / Go back |
| `q` | Quit |

## How to Use

1. **First launch**: The app starts empty
2. **Add timezone**: Press `a`, then type to search (e.g., "Tokyo" or "Japan")
3. **Select**: Use arrow keys to choose, press `Enter`
4. **Remove**: Select a timezone and press `d`
5. **Quit**: Press `q`

Your timezone list is saved automatically and will be there next time you open the app.

## Building from Source

### Linux / macOS

```bash
# Clone
git clone https://github.com/ikpvk/tui_world_clock_go.git
cd tui_world_clock_go

# Run directly
go run ./cmd

# Build binary
go build -o worldclock ./cmd

# Run the built binary
./worldclock
```

### Windows

```bash
# Clone
git clone https://github.com/ikpvk/tui_world_clock_go.git
cd tui_world_clock_go

# Build for Windows
go build -tags timetzdata -o worldclock.exe ./cmd

# Run
worldclock.exe
```

Or download the pre-built `worldclock.exe` from the Releases page.

## Supported Platforms

- Linux (x64, ARM64)
- macOS (x64, ARM64)
- Windows (x64)

## Technical Details

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss)
- Timezone data embedded for Windows compatibility

## License

MIT
