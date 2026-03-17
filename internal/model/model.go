package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"worldclock/internal/config"
	"worldclock/internal/timezone"
)

type Theme struct {
	Name       string
	Foreground string
	Background string
	Header     string
	SelectedBg string
	SelectedFg string
}

var themes = []Theme{
	{Name: "dracula", Foreground: "rgb(189,147,249)", Background: "rgb(40,42,54)", Header: "rgb(139,233,253)", SelectedBg: "rgb(97,175,239)", SelectedFg: "rgb(255,255,255)"},
	{Name: "nord", Foreground: "rgb(136,192,208)", Background: "rgb(46,52,64)", Header: "rgb(136,192,208)", SelectedBg: "rgb(136,192,208)", SelectedFg: "rgb(46,52,64)"},
	{Name: "one-dark", Foreground: "rgb(198,120,221)", Background: "rgb(40,44,52)", Header: "rgb(97,175,239)", SelectedBg: "rgb(97,175,239)", SelectedFg: "rgb(40,44,52)"},
	{Name: "gruvbox", Foreground: "rgb(204,36,29)", Background: "rgb(40,40,40)", Header: "rgb(215,153,33)", SelectedBg: "rgb(215,153,33)", SelectedFg: "rgb(40,40,40)"},
	{Name: "solarized-light", Foreground: "rgb(181,137,0)", Background: "rgb(253,246,227)", Header: "rgb(42,161,152)", SelectedBg: "rgb(42,161,152)", SelectedFg: "rgb(253,246,227)"},
	{Name: "github-light", Foreground: "rgb(3,102,214)", Background: "rgb(255,255,255)", Header: "rgb(28,92,105)", SelectedBg: "rgb(40,167,69)", SelectedFg: "rgb(255,255,255)"},
	{Name: "monokai-light", Foreground: "rgb(242,38,114)", Background: "rgb(250,248,245)", Header: "rgb(166,226,46)", SelectedBg: "rgb(166,226,46)", SelectedFg: "rgb(250,248,245)"},
	{Name: "paper", Foreground: "rgb(0,0,255)", Background: "rgb(255,255,255)", Header: "rgb(0,0,128)", SelectedBg: "rgb(0,0,128)", SelectedFg: "rgb(255,255,255)"},
}

func getThemeIndex(name string) int {
	for i, t := range themes {
		if t.Name == name {
			return i
		}
	}
	return 0
}

type ViewMode string

const (
	ListView ViewMode = "list"
	AddView  ViewMode = "add"
)

type Model struct {
	timezones  []string
	selected   int
	view       ViewMode
	input      textinput.Model
	ticker     *time.Ticker
	done       bool
	themeIndex int
}

func New() *Model {
	ti := textinput.New()
	ti.Placeholder = "Search timezone..."
	ti.Focus()

	return &Model{
		timezones:  []string{},
		selected:   0,
		view:       ListView,
		input:      ti,
		ticker:     time.NewTicker(time.Second),
		themeIndex: 0,
	}
}

func (m *Model) SetTimezones(tz []string) {
	m.timezones = tz
}

func (m *Model) SetTheme(themeName string) {
	m.themeIndex = getThemeIndex(themeName)
}

func (m *Model) GetThemeName() string {
	return themes[m.themeIndex].Name
}

type TickMsg time.Time

func (m *Model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.view == AddView {
			if msg.String() == "up" || msg.String() == "k" {
				m.selected = max(0, m.selected-1)
				return m, nil
			}
			if msg.String() == "down" || msg.String() == "j" {
				filtered := timezone.FilterTimezones(m.input.Value())
				if len(filtered) > 0 {
					m.selected = min(len(filtered)-1, m.selected+1)
				}
				return m, nil
			}
			if msg.String() == "enter" {
				return m.handleAddConfirm()
			}
			if msg.String() == "esc" {
				m.view = ListView
				m.input.SetValue("")
				return m, nil
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			filtered := timezone.FilterTimezones(m.input.Value())
			if len(filtered) > 0 {
				m.selected = len(filtered) - 1
			}
			return m, cmd
		}
		return m.handleKeyMsg(msg)
	case TickMsg:
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return TickMsg(t)
		})
	}

	return m, nil
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.done = true
		return m, tea.Quit
	case "a":
		m.view = AddView
		filtered := timezone.FilterTimezones("")
		if len(filtered) > 0 {
			m.selected = len(filtered) - 1
		} else {
			m.selected = 0
		}
		m.input.Focus()
		return m, nil
	case "t":
		m.themeIndex = (m.themeIndex + 1) % len(themes)
		m.saveConfig()
		return m, nil
	case "T":
		m.themeIndex = (m.themeIndex - 1 + len(themes)) % len(themes)
		m.saveConfig()
		return m, nil
	case "esc":
		if m.view == AddView {
			m.view = ListView
			m.input.SetValue("")
		}
		return m, nil
	case "enter":
		if m.view == AddView {
			return m.handleAddConfirm()
		}
	case "d", "delete":
		if m.view == ListView && len(m.timezones) > 0 {
			return m.handleDelete()
		}
	case "up", "k":
		m.selected = max(0, m.selected-1)
	case "down", "j":
		m.selected = min(len(m.timezones)-1, m.selected+1)
	}
	return m, nil
}

func (m *Model) handleAddConfirm() (tea.Model, tea.Cmd) {
	filtered := timezone.FilterTimezones(m.input.Value())
	if len(filtered) > 0 && m.selected >= 0 && m.selected < len(filtered) {
		tz := filtered[m.selected]
		if !contains(m.timezones, tz) {
			m.timezones = append(m.timezones, tz)
			m.saveConfig()
		}
	}
	m.view = ListView
	m.input.SetValue("")
	m.selected = 0
	return m, nil
}

func (m *Model) handleDelete() (tea.Model, tea.Cmd) {
	if m.selected >= 0 && m.selected < len(m.timezones) {
		m.timezones = append(m.timezones[:m.selected], m.timezones[m.selected+1:]...)
		if m.selected >= len(m.timezones) && len(m.timezones) > 0 {
			m.selected = len(m.timezones) - 1
		}
		m.saveConfig()
	}
	return m, nil
}

func (m *Model) saveConfig() {
	cfg := &config.Config{
		Timezones: m.timezones,
		Theme:     themes[m.themeIndex].Name,
	}
	config.Save(cfg)
}

func (m *Model) View() string {
	if m.done {
		return ""
	}

	var content string

	if m.view == AddView {
		content = m.renderAddView()
	} else {
		content = m.renderListView()
	}

	theme := themes[m.themeIndex]
	bgStyle := lipgloss.NewStyle().Background(lipgloss.Color(theme.Background))
	return bgStyle.Render(content)
}

const (
	colTimezone = 25
	colCountry  = 18
	colDateTime = 20
)

func (m *Model) renderListView() string {
	var lines []string

	theme := themes[m.themeIndex]

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Header)).
		Bold(true)

	lines = append(lines, headerStyle.Render("World Clock"))
	lines = append(lines, "")

	if len(m.timezones) == 0 {
		lines = append(lines, "No timezones added.")
		lines = append(lines, "Press 'a' to add a timezone")
	} else {
		headerTimezone := fmt.Sprintf("%-*s", colTimezone, "Timezone")
		headerCountry := fmt.Sprintf("%-*s", colCountry, "Country")
		headerDateTime := fmt.Sprintf("%-*s", colDateTime, "Date/Time")
		lines = append(lines, fmt.Sprintf("| %s | %s | %s |", headerTimezone, headerCountry, headerDateTime))

		sepTimezone := strings.Repeat("-", colTimezone)
		sepCountry := strings.Repeat("-", colCountry)
		sepDateTime := strings.Repeat("-", colDateTime)
		lines = append(lines, fmt.Sprintf("| %s | %s | %s |", sepTimezone, sepCountry, sepDateTime))

		for i, tz := range m.timezones {
			timeStr := timezone.GetCurrentTimeForZone(tz)
			displayName := timezone.GetDisplayName(tz)

			tzName, country := splitDisplayName(displayName)

			colTz := fmt.Sprintf("%-*s", colTimezone, tzName)
			colCountry := fmt.Sprintf("%-*s", colCountry, country)
			colDt := fmt.Sprintf("%-*s", colDateTime, timeStr)

			row := fmt.Sprintf("| %s | %s | %s |", colTz, colCountry, colDt)

			if i == m.selected {
				selectedStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color(theme.SelectedFg)).
					Background(lipgloss.Color(theme.SelectedBg))
				lines = append(lines, selectedStyle.Render(row))
			} else {
				lines = append(lines, row)
			}
		}
	}

	lines = append(lines, "")
	lines = append(lines, "Navigate: j/k or up/down | a: Add | d: Delete | t/T: Theme | q: Quit")

	return lipgloss.NewStyle().Render("\n" + joinLines(lines))
}

func splitDisplayName(displayName string) (string, string) {
	if idx := strings.Index(displayName, " ("); idx != -1 {
		tz := displayName[:idx]
		country := displayName[idx+2 : len(displayName)-1]
		return tz, country
	}
	return displayName, ""
}

func (m *Model) renderAddView() string {
	var lines []string

	theme := themes[m.themeIndex]

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Header)).
		Bold(true)

	lines = append(lines, headerStyle.Render("Add Timezone"))
	lines = append(lines, "")

	filtered := timezone.FilterTimezones(m.input.Value())

	if len(filtered) == 0 {
		lines = append(lines, "No matching timezones")
	} else {
		for i, tz := range filtered {
			displayName := timezone.GetDisplayName(tz)
			if i == m.selected {
				selectedStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color(theme.SelectedFg)).
					Background(lipgloss.Color(theme.SelectedBg))
				lines = append(lines, selectedStyle.Render("> "+displayName))
			} else {
				lines = append(lines, "  "+displayName)
			}
		}
	}

	lines = append(lines, "")
	lines = append(lines, m.input.View())
	lines = append(lines, "")
	lines = append(lines, "Navigate: j/k or up/down | Enter: Select | Esc: Cancel")

	return lipgloss.NewStyle().Render("\n" + joinLines(lines))
}

func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		result += line
		if i < len(lines)-1 {
			result += "\n"
		}
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
