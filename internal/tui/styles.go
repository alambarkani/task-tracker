package tui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	colorAccent    = lipgloss.Color("#8B5CF6") // violet
	colorAccentDim = lipgloss.Color("#5B21B6")
	colorTodo      = lipgloss.Color("#60A5FA") // blue
	colorProgress  = lipgloss.Color("#FBBF24") // amber
	colorDone      = lipgloss.Color("#34D399") // green
	colorError     = lipgloss.Color("#F87171") // red
	colorMuted     = lipgloss.Color("#6B7280")
	colorText      = lipgloss.Color("#E5E7EB")
	colorSubtle    = lipgloss.Color("#9CA3AF")
	colorBorder    = lipgloss.Color("#3F3F55")
	colorHighlight = lipgloss.Color("#2D2B44")
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleBarStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#0B0B14")).
			Background(colorAccent).
			Padding(0, 2)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorSubtle).
			Italic(true)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)

	footerStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			MarginTop(1)

	keyStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	descKeyStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	cursorStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	rowSelectedStyle = lipgloss.NewStyle().
				Background(colorHighlight).
				Foreground(colorText).
				Bold(true)

	rowStyle = lipgloss.NewStyle().
			Foreground(colorText)

	idStyle = lipgloss.NewStyle().
		Foreground(colorMuted).
		Width(4)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(colorDone).
			Bold(true)

	inputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(1, 2)

	dangerBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorError).
			Padding(1, 2)

	tabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#0B0B14")).
			Background(colorAccent).
			Padding(0, 1)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(colorSubtle).
				Padding(0, 1)

	logoStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)
)

// statusStyle returns the lipgloss style used to render a given status badge.
func statusStyle(s string) lipgloss.Style {
	switch s {
	case "todo":
		return lipgloss.NewStyle().Foreground(colorTodo).Bold(true)
	case "in_progress":
		return lipgloss.NewStyle().Foreground(colorProgress).Bold(true)
	case "done":
		return lipgloss.NewStyle().Foreground(colorDone).Bold(true)
	default:
		return lipgloss.NewStyle().Foreground(colorMuted)
	}
}

// statusIcon returns a small glyph representing a task status.
func statusIcon(s string) string {
	switch s {
	case "todo":
		return "○"
	case "in_progress":
		return "◐"
	case "done":
		return "●"
	default:
		return "?"
	}
}

// statusLabel returns a human friendly label for a status.
func statusLabel(s string) string {
	switch s {
	case "todo":
		return "Todo"
	case "in_progress":
		return "In Progress"
	case "done":
		return "Done"
	default:
		return s
	}
}
