package main

import (
	"fmt"
	"os"

	"alambarkani.com/task_tracker/internal/dir_helper"
	"alambarkani.com/task_tracker/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	taskPath := dir_helper.TaskDir()

	p := tea.NewProgram(tui.New(taskPath), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "task_tracker: ", err)
		os.Exit(1)
	}
}
