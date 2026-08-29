package tui

import (
	"alambarkani.com/task_tracker/internal/task_tracker"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Screen identifies which view the Model is currently rendering.
type Screen int

const (
	ScreenSplash Screen = iota
	ScreenTaskList
	ScreenAddTask
	ScreenEditTask
	ScreenDeleteConfirm
	ScreenHelp
)

// filter values used to narrow the task list by status.
const (
	FilterAll        = "all"
	FilterTodo       = "todo"
	FilterInProgress = "in_progress"
	FilterDone       = "done"
)

var filterOrder = []string{FilterAll, FilterTodo, FilterInProgress, FilterDone}

// Model holds all TUI application state.
type Model struct {
	Screen     Screen
	PrevScreen Screen

	Width  int
	Height int

	TaskPath string
	Tasks    []task_tracker.Task

	Cursor      int // selected row within the filtered task list
	Scroll      int // first visible row within the filtered task list
	FilterIndex int // index into filterOrder

	Input     textinput.Model
	EditingID int

	DeleteID          int
	DeleteDescription string

	StatusMsg string
	ErrMsg    string

	Quitting bool
}

// New builds the initial Model for the task tracker TUI.
func New(taskPath string) Model {
	ti := textinput.New()
	ti.Placeholder = "What needs to be done?"
	ti.CharLimit = 200
	ti.Width = 48

	return Model{
		Screen:      ScreenSplash,
		PrevScreen:  ScreenSplash,
		TaskPath:    taskPath,
		Input:       ti,
		FilterIndex: 0,
	}
}

// Init kicks off the initial commands: load tasks and start the splash timer.
func (m Model) Init() tea.Cmd {
	return tea.Batch(loadTasksCmd(m.TaskPath), splashTickCmd())
}

// CurrentFilter returns the active status filter.
func (m Model) CurrentFilter() string {
	return filterOrder[m.FilterIndex]
}

// FilteredTasks returns Tasks narrowed by the active filter, in stable order.
func (m Model) FilteredTasks() []task_tracker.Task {
	f := m.CurrentFilter()
	if f == FilterAll {
		return m.Tasks
	}
	out := make([]task_tracker.Task, 0, len(m.Tasks))
	for _, t := range m.Tasks {
		if string(t.Status) == f {
			out = append(out, t)
		}
	}
	return out
}

// GetScreenTitle returns the title shown in the title bar for the current screen.
func (m Model) GetScreenTitle() string {
	switch m.Screen {
	case ScreenSplash:
		return "Task Tracker"
	case ScreenTaskList:
		return "Task Tracker"
	case ScreenAddTask:
		return "New Task"
	case ScreenEditTask:
		return "Edit Task"
	case ScreenDeleteConfirm:
		return "Delete Task"
	case ScreenHelp:
		return "Help"
	default:
		return "Task Tracker"
	}
}
