package tui

import (
	"time"

	"alambarkani.com/task_tracker/internal/task_tracker"
	tea "github.com/charmbracelet/bubbletea"
)

// tasksLoadedMsg carries the result of reading tasks.json from disk.
type tasksLoadedMsg struct {
	tasks []task_tracker.Task
	err   error
}

// actionDoneMsg carries the result of a mutating action (add/edit/delete/mark)
// along with a human-readable status line to show once the reload completes.
type actionDoneMsg struct {
	err       error
	statusMsg string
}

// splashDoneMsg fires once the splash screen's minimum display time has elapsed.
type splashDoneMsg struct{}

// loadTasksCmd reads the task file and reports the result.
func loadTasksCmd(path string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := task_tracker.TaskList(path)
		if err != nil {
			// A missing/empty file is not an error condition for the TUI —
			// it just means there are no tasks yet.
			return tasksLoadedMsg{tasks: nil, err: nil}
		}
		return tasksLoadedMsg{tasks: tasks, err: nil}
	}
}

// splashTickCmd waits briefly so the splash screen is visible for a beat
// before handing off to the task list.
func splashTickCmd() tea.Cmd {
	return tea.Tick(900*time.Millisecond, func(time.Time) tea.Msg {
		return splashDoneMsg{}
	})
}

// addTaskCmd creates a task then reports completion.
func addTaskCmd(path, description string) tea.Cmd {
	return func() tea.Msg {
		err := task_tracker.AddTask(path, description)
		if err != nil {
			return actionDoneMsg{err: err}
		}
		return actionDoneMsg{statusMsg: "Task added"}
	}
}

// editTaskCmd updates a task's description then reports completion.
func editTaskCmd(path string, id int, description string) tea.Cmd {
	return func() tea.Msg {
		err := task_tracker.EditTask(path, id, description)
		if err != nil {
			return actionDoneMsg{err: err}
		}
		return actionDoneMsg{statusMsg: "Task updated"}
	}
}

// deleteTaskCmd removes a task then reports completion.
func deleteTaskCmd(path string, id int) tea.Cmd {
	return func() tea.Msg {
		err := task_tracker.DeleteTask(path, id)
		if err != nil {
			return actionDoneMsg{err: err}
		}
		return actionDoneMsg{statusMsg: "Task deleted"}
	}
}

// markTaskCmd sets a task's status then reports completion.
func markTaskCmd(path string, id int, status task_tracker.Status) tea.Cmd {
	return func() tea.Msg {
		err := task_tracker.MarkTask(path, id, status)
		if err != nil {
			return actionDoneMsg{err: err}
		}
		return actionDoneMsg{statusMsg: "Marked as " + statusLabel(string(status))}
	}
}
