package main

import (
	"fmt"
	"os"

	_ "time/tzdata"

	tea "github.com/charmbracelet/bubbletea"
	"worldclock/internal/config"
	"worldclock/internal/model"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	m := model.New()
	m.SetTimezones(cfg.Timezones)
	m.SetTheme(cfg.Theme)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting app: %v\n", err)
		os.Exit(1)
	}
}
