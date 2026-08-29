package tui

import (
	"fmt"
	"strings"

	"alambarkani.com/task_tracker/internal/task_tracker"
	"github.com/charmbracelet/lipgloss"
)

const logo = `
#####  ###   #### #   #
  #   #   # #     #  #
  #   #####  ###  ###
  #   #   #     # #  #
  #   #   # ####  #   #

##### ####   ###   #### #   # ##### ####
  #   #   # #   # #     #  #  #     #   #
  #   ####  ##### #     ###   ###   ####
  #   #  #  #   # #     #  #  #     #  #
  #   #   # #   #  #### #   # ##### #   #
`

// View renders the current screen.
func (m Model) View() string {
	if m.Quitting {
		return ""
	}
	if m.Width == 0 {
		return "Loading…"
	}

	var body string
	switch m.Screen {
	case ScreenSplash:
		body = m.viewSplash()
	case ScreenTaskList:
		body = m.viewTaskList()
	case ScreenAddTask:
		body = m.viewForm("New Task", "Description", m.Input.View(), "enter confirm • esc cancel")
	case ScreenEditTask:
		body = m.viewForm(fmt.Sprintf("Edit Task #%d", m.EditingID), "Description", m.Input.View(), "enter save • esc cancel")
	case ScreenDeleteConfirm:
		body = m.viewDeleteConfirm()
	case ScreenHelp:
		body = m.viewHelp()
	default:
		body = ""
	}

	return appStyle.Render(body)
}

func (m Model) viewSplash() string {
	art := logoStyle.Render(strings.TrimRight(logo, "\n"))
	sub := subtitleStyle.Render("track what matters, right from your terminal")
	credit := descKeyStyle.Render("by Alam Barkani")
	return lipgloss.JoinVertical(lipgloss.Center, art, "", sub, "", credit)
}

func (m Model) viewTaskList() string {
	title := titleBarStyle.Render(" ✓ Task Tracker ")

	tabs := m.viewTabs()

	tasks := m.FilteredTasks()
	var list string
	if len(tasks) == 0 {
		list = panelStyle.Width(60).Render(subtitleStyle.Render("No tasks here yet. Press " + keyStyle.Render("a") + subtitleStyle.Render(" to add one.")))
	} else {
		list = m.viewTaskRows(tasks)
	}

	status := m.viewStatusLine()
	footer := footerStyle.Render(
		keyStyle.Render("a") + descKeyStyle.Render(" add  ") +
			keyStyle.Render("e") + descKeyStyle.Render(" edit  ") +
			keyStyle.Render("d") + descKeyStyle.Render(" delete  ") +
			keyStyle.Render("space") + descKeyStyle.Render(" cycle status  ") +
			keyStyle.Render("←/→") + descKeyStyle.Render(" filter  ") +
			keyStyle.Render("?") + descKeyStyle.Render(" help  ") +
			keyStyle.Render("q") + descKeyStyle.Render(" quit"),
	)

	sections := []string{title, "", tabs, "", list}
	if status != "" {
		sections = append(sections, status)
	}
	sections = append(sections, footer)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) viewTabs() string {
	labels := map[string]string{
		FilterAll:        "All",
		FilterTodo:       "Todo",
		FilterInProgress: "In Progress",
		FilterDone:       "Done",
	}
	var parts []string
	for i, f := range filterOrder {
		label := labels[f]
		if i == m.FilterIndex {
			parts = append(parts, tabActiveStyle.Render(label))
		} else {
			parts = append(parts, tabInactiveStyle.Render(label))
		}
	}
	return strings.Join(parts, " ")
}

func (m Model) viewTaskRows(tasks []task_tracker.Task) string {
	visible := m.visibleRows()
	start := m.Scroll
	end := min(start+visible, len(tasks))
	if start > end {
		start = end
	}

	descWidth := 46
	var rows []string
	for i := start; i < end; i++ {
		t := tasks[i]
		icon := statusStyle(string(t.Status)).Render(statusIcon(string(t.Status)))
		id := idStyle.Render(fmt.Sprintf("#%-3d", t.ID))
		desc := truncate(t.Description, descWidth)

		line := fmt.Sprintf("%s %s %-*s", icon, id, descWidth, desc)

		if i == m.Cursor {
			line = cursorStyle.Render("▶ ") + rowSelectedStyle.Render(line)
		} else {
			line = "  " + rowStyle.Render(line)
		}
		rows = append(rows, line)
	}

	body := strings.Join(rows, "\n")

	scrollHint := ""
	if len(tasks) > visible {
		scrollHint = "\n" + descKeyStyle.Render(fmt.Sprintf("  showing %d-%d of %d", start+1, end, len(tasks)))
	}

	return panelStyle.Width(70).Render(body + scrollHint)
}

func truncate(s string, width int) string {
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	if width <= 1 {
		return string(r[:width])
	}
	return string(r[:width-1]) + "…"
}

func (m Model) viewStatusLine() string {
	if m.ErrMsg != "" {
		return errorStyle.Render("✗ " + m.ErrMsg)
	}
	if m.StatusMsg != "" {
		return successStyle.Render("✓ " + m.StatusMsg)
	}
	return ""
}

func (m Model) viewForm(title, label, input, hint string) string {
	heading := titleBarStyle.Render(" " + title + " ")
	box := inputBoxStyle.Width(56).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			descKeyStyle.Render(label+":"),
			"",
			input,
		),
	)
	var errLine string
	if m.ErrMsg != "" {
		errLine = errorStyle.Render("✗ " + m.ErrMsg)
	}
	footer := footerStyle.Render(hint)

	sections := []string{heading, "", box}
	if errLine != "" {
		sections = append(sections, errLine)
	}
	sections = append(sections, footer)
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) viewDeleteConfirm() string {
	heading := titleBarStyle.Render(" Delete Task ")
	msg := fmt.Sprintf("Delete task #%d?\n\n%s", m.DeleteID, rowStyle.Render(`"`+m.DeleteDescription+`"`))
	box := dangerBoxStyle.Width(56).Render(msg)
	footer := footerStyle.Render(keyStyle.Render("y") + descKeyStyle.Render(" confirm  ") + keyStyle.Render("n/esc") + descKeyStyle.Render(" cancel"))
	return lipgloss.JoinVertical(lipgloss.Left, heading, "", box, footer)
}

func (m Model) viewHelp() string {
	heading := titleBarStyle.Render(" Help ")

	rows := [][2]string{
		{"↑/k ↓/j", "move selection"},
		{"←/h →/l / tab", "change filter tab"},
		{"a", "add a new task"},
		{"e", "edit selected task"},
		{"d", "delete selected task (confirm)"},
		{"space / n", "cycle status: todo → in progress → done"},
		{"1 / 2 / 3", "set status directly (todo/in progress/done)"},
		{"r", "reload from disk"},
		{"?", "toggle this help"},
		{"q / ctrl+c", "quit"},
	}

	var b strings.Builder
	for _, r := range rows {
		b.WriteString(keyStyle.Width(16).Render(r[0]))
		b.WriteString(descKeyStyle.Render(r[1]))
		b.WriteString("\n")
	}

	box := panelStyle.Width(56).Render(b.String())
	footer := footerStyle.Render("esc / ? / enter to go back")
	credit := descKeyStyle.Render("Task Tracker — by Alam Barkani")
	return lipgloss.JoinVertical(lipgloss.Left, heading, "", box, footer, credit)
}
