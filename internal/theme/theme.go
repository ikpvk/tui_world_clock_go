package theme

type Theme struct {
	Name       string
	Foreground string
	Background string
	Header     string
	SelectedBg string
	SelectedFg string
}

var Themes = []Theme{
	{Name: "dracula", Foreground: "#bd93f9", Background: "#282a36", Header: "#ff79c6", SelectedBg: "#61afef", SelectedFg: "#ffffff"},
	{Name: "nord", Foreground: "#88c0d0", Background: "#2e3440", Header: "#81a1c1", SelectedBg: "#88c0d0", SelectedFg: "#2e3440"},
	{Name: "one-dark", Foreground: "#c678dd", Background: "#282c34", Header: "#98c379", SelectedBg: "#61afef", SelectedFg: "#282c34"},
	{Name: "gruvbox", Foreground: "#cc241d", Background: "#282828", Header: "#d79921", SelectedBg: "#d79921", SelectedFg: "#282828"},
	{Name: "solarized-light", Foreground: "#b58900", Background: "#fdf6e3", Header: "#cb4b16", SelectedBg: "#2aa198", SelectedFg: "#fdf6e3"},
	{Name: "github-light", Foreground: "#0366d6", Background: "#ffffff", Header: "#28a745", SelectedBg: "#28a745", SelectedFg: "#ffffff"},
	{Name: "monokai-light", Foreground: "#f92672", Background: "#faf8f5", Header: "#f4bf75", SelectedBg: "#a6e22e", SelectedFg: "#faf8f5"},
	{Name: "paper", Foreground: "#0000ff", Background: "#ffffff", Header: "#000080", SelectedBg: "#000080", SelectedFg: "#ffffff"},
}

func GetThemeIndex(name string) int {
	for i, t := range Themes {
		if t.Name == name {
			return i
		}
	}
	return 0
}

func GetThemeByIndex(index int) Theme {
	if index < 0 || index >= len(Themes) {
		return Themes[0]
	}
	return Themes[index]
}
