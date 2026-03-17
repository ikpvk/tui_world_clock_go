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

type ViewMode string

const (
	ListView ViewMode = "list"
	AddView  ViewMode = "add"
)

type Model struct {
	timezones []string
	selected  int
	view      ViewMode
	input     textinput.Model
	ticker    *time.Ticker
	done      bool
}

func New() *Model {
	ti := textinput.New()
	ti.Placeholder = "Search timezone..."
	ti.Focus()

	return &Model{
		timezones: []string{},
		selected:  0,
		view:      ListView,
		input:     ti,
		ticker:    time.NewTicker(time.Second),
	}
}

func (m *Model) SetTimezones(tz []string) {
	m.timezones = tz
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
			cfg := &config.Config{Timezones: m.timezones}
			config.Save(cfg)
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
		cfg := &config.Config{Timezones: m.timezones}
		config.Save(cfg)
	}
	return m, nil
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

	return content
}

const (
	colTimezone = 25
	colCountry  = 18
	colDateTime = 20
)

func (m *Model) renderListView() string {
	var lines []string

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
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
					Foreground(lipgloss.Color("15")).
					Background(lipgloss.Color("57"))
				lines = append(lines, selectedStyle.Render(row))
			} else {
				lines = append(lines, row)
			}
		}
	}

	lines = append(lines, "")
	lines = append(lines, "↑↓ Navigate | a: Add | d: Delete | q: Quit")

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

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
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
					Foreground(lipgloss.Color("15")).
					Background(lipgloss.Color("57"))
				lines = append(lines, selectedStyle.Render("> "+displayName))
			} else {
				lines = append(lines, "  "+displayName)
			}
		}
	}

	lines = append(lines, "")
	lines = append(lines, m.input.View())
	lines = append(lines, "")
	lines = append(lines, "↑↓ Navigate | Enter: Select | Esc: Cancel")

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
