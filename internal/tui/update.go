package tui

import (
	"alambarkani.com/task_tracker/internal/task_tracker"
	tea "github.com/charmbracelet/bubbletea"
)

// Update is the single entry point for all state transitions.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case splashDoneMsg:
		if m.Screen == ScreenSplash {
			m.Screen = ScreenTaskList
		}
		return m, nil

	case tasksLoadedMsg:
		m.Tasks = msg.tasks
		if m.Cursor >= len(m.FilteredTasks()) {
			m.Cursor = max(len(m.FilteredTasks())-1, 0)
		}
		return m, nil

	case actionDoneMsg:
		if msg.err != nil {
			m.ErrMsg = msg.err.Error()
			m.StatusMsg = ""
			return m, nil
		}
		m.ErrMsg = ""
		m.StatusMsg = msg.statusMsg
		return m, loadTasksCmd(m.TaskPath)
	}

	return m, nil
}

// handleKeyPress dispatches a key press to the handler for the current screen.
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global quit, but not while typing in a text field.
	if key == "ctrl+c" {
		m.Quitting = true
		return m, tea.Quit
	}

	switch m.Screen {
	case ScreenSplash:
		m.Screen = ScreenTaskList
		return m, nil
	case ScreenTaskList:
		return m.handleTaskListKeys(key)
	case ScreenAddTask:
		return m.handleAddTaskKeys(msg)
	case ScreenEditTask:
		return m.handleEditTaskKeys(msg)
	case ScreenDeleteConfirm:
		return m.handleDeleteConfirmKeys(key)
	case ScreenHelp:
		return m.handleHelpKeys(key)
	}

	return m, nil
}

func (m Model) handleTaskListKeys(key string) (tea.Model, tea.Cmd) {
	tasks := m.FilteredTasks()

	switch key {
	case "q":
		m.Quitting = true
		return m, tea.Quit

	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
		m.clampScroll()

	case "down", "j":
		if m.Cursor < len(tasks)-1 {
			m.Cursor++
		}
		m.clampScroll()

	case "left", "h":
		m.FilterIndex = (m.FilterIndex - 1 + len(filterOrder)) % len(filterOrder)
		m.Cursor = 0
		m.Scroll = 0

	case "right", "l", "tab":
		m.FilterIndex = (m.FilterIndex + 1) % len(filterOrder)
		m.Cursor = 0
		m.Scroll = 0

	case "a":
		m.PrevScreen = ScreenTaskList
		m.Screen = ScreenAddTask
		m.Input.SetValue("")
		m.ErrMsg = ""
		m.StatusMsg = ""
		return m, m.Input.Focus()

	case "e":
		if t, ok := m.selectedTask(); ok {
			m.PrevScreen = ScreenTaskList
			m.Screen = ScreenEditTask
			m.EditingID = t.ID
			m.Input.SetValue(t.Description)
			m.Input.CursorEnd()
			m.ErrMsg = ""
			m.StatusMsg = ""
			return m, m.Input.Focus()
		}

	case "d":
		if t, ok := m.selectedTask(); ok {
			m.PrevScreen = ScreenTaskList
			m.Screen = ScreenDeleteConfirm
			m.DeleteID = t.ID
			m.DeleteDescription = t.Description
			m.ErrMsg = ""
			m.StatusMsg = ""
		}

	case " ", "n":
		// Cycle status forward: todo -> in_progress -> done -> todo
		if t, ok := m.selectedTask(); ok {
			next := nextStatus(t.Status)
			m.StatusMsg = ""
			m.ErrMsg = ""
			return m, markTaskCmd(m.TaskPath, t.ID, next)
		}

	case "1":
		if t, ok := m.selectedTask(); ok {
			return m, markTaskCmd(m.TaskPath, t.ID, task_tracker.Todo)
		}
	case "2":
		if t, ok := m.selectedTask(); ok {
			return m, markTaskCmd(m.TaskPath, t.ID, task_tracker.InProgress)
		}
	case "3":
		if t, ok := m.selectedTask(); ok {
			return m, markTaskCmd(m.TaskPath, t.ID, task_tracker.Done)
		}

	case "r":
		return m, loadTasksCmd(m.TaskPath)

	case "?":
		m.PrevScreen = ScreenTaskList
		m.Screen = ScreenHelp
	}

	return m, nil
}

func (m Model) handleAddTaskKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.Screen = m.PrevScreen
		m.Input.Blur()
		return m, nil
	case "enter":
		desc := trimmed(m.Input.Value())
		if desc == "" {
			m.ErrMsg = "Description can't be empty"
			return m, nil
		}
		m.Screen = m.PrevScreen
		m.Input.Blur()
		return m, addTaskCmd(m.TaskPath, desc)
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m Model) handleEditTaskKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.Screen = m.PrevScreen
		m.Input.Blur()
		return m, nil
	case "enter":
		desc := trimmed(m.Input.Value())
		if desc == "" {
			m.ErrMsg = "Description can't be empty"
			return m, nil
		}
		m.Screen = m.PrevScreen
		m.Input.Blur()
		return m, editTaskCmd(m.TaskPath, m.EditingID, desc)
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m Model) handleDeleteConfirmKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "y", "enter":
		id := m.DeleteID
		m.Screen = m.PrevScreen
		return m, deleteTaskCmd(m.TaskPath, id)
	case "n", "esc":
		m.Screen = m.PrevScreen
	}
	return m, nil
}

func (m Model) handleHelpKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc", "q", "enter", " ", "?":
		m.Screen = m.PrevScreen
	}
	return m, nil
}

// selectedTask returns the task under the cursor in the filtered list, if any.
func (m Model) selectedTask() (task_tracker.Task, bool) {
	tasks := m.FilteredTasks()
	if m.Cursor < 0 || m.Cursor >= len(tasks) {
		return task_tracker.Task{}, false
	}
	return tasks[m.Cursor], true
}

// clampScroll keeps the cursor within the visible window of the list.
func (m *Model) clampScroll() {
	visible := m.visibleRows()
	if visible <= 0 {
		return
	}
	if m.Cursor < m.Scroll {
		m.Scroll = m.Cursor
	}
	if m.Cursor >= m.Scroll+visible {
		m.Scroll = m.Cursor - visible + 1
	}
}

// visibleRows computes how many task rows fit in the current terminal height.
func (m Model) visibleRows() int {
	// Reserve space for title bar, tabs, borders, and footer.
	return max(m.Height-10, 3)
}

func nextStatus(s task_tracker.Status) task_tracker.Status {
	switch s {
	case task_tracker.Todo:
		return task_tracker.InProgress
	case task_tracker.InProgress:
		return task_tracker.Done
	default:
		return task_tracker.Todo
	}
}

func trimmed(s string) string {
	start, end := 0, len(s)
	for start < end && isSpace(s[start]) {
		start++
	}
	for end > start && isSpace(s[end-1]) {
		end--
	}
	return s[start:end]
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
