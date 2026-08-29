package tui

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func newTestModel(t *testing.T) Model {
	t.Helper()
	dir := t.TempDir()
	m := New(filepath.Join(dir, "tasks.json"))
	m.Width, m.Height = 80, 30
	m.Screen = ScreenTaskList // skip splash for these tests
	return m
}

// runCmd executes a tea.Cmd synchronously and feeds the resulting message(s)
// back through Update, returning the final model. Nil-safe and unwraps
// tea.BatchMsg so command chains (e.g. save -> reload) settle fully.
func runCmd(t *testing.T, m Model, cmd tea.Cmd) Model {
	t.Helper()
	for cmd != nil {
		msg := cmd()
		if msg == nil {
			return m
		}
		var updated tea.Model
		updated, cmd = m.Update(msg)
		m = updated.(Model)
	}
	return m
}

func typeRunes(t *testing.T, m Model, s string) Model {
	t.Helper()
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)})
	m = updated.(Model)
	return runCmd(t, m, cmd)
}

func pressKey(t *testing.T, m Model, k tea.KeyMsg) Model {
	t.Helper()
	updated, cmd := m.Update(k)
	m = updated.(Model)
	return runCmd(t, m, cmd)
}

func TestSplashTransitionsToTaskList(t *testing.T) {
	m := New(filepath.Join(t.TempDir(), "tasks.json"))
	if m.Screen != ScreenSplash {
		t.Fatalf("expected initial screen to be ScreenSplash, got %v", m.Screen)
	}

	updated, _ := m.Update(splashDoneMsg{})
	m = updated.(Model)
	if m.Screen != ScreenTaskList {
		t.Fatalf("expected ScreenTaskList after splashDoneMsg, got %v", m.Screen)
	}
}

func TestAddTaskFlow(t *testing.T) {
	m := newTestModel(t)

	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	if m.Screen != ScreenAddTask {
		t.Fatalf("expected ScreenAddTask, got %v", m.Screen)
	}

	m = typeRunes(t, m, "Buy milk")
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	if m.Screen != ScreenTaskList {
		t.Fatalf("expected to return to ScreenTaskList, got %v", m.Screen)
	}
	if m.ErrMsg != "" {
		t.Fatalf("unexpected error: %s", m.ErrMsg)
	}
	tasks := m.FilteredTasks()
	if len(tasks) != 1 || tasks[0].Description != "Buy milk" {
		t.Fatalf("expected one task 'Buy milk', got %+v", tasks)
	}
}

func TestAddTaskRejectsEmptyDescription(t *testing.T) {
	m := newTestModel(t)
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	if m.Screen != ScreenAddTask {
		t.Fatalf("expected to stay on ScreenAddTask on empty input, got %v", m.Screen)
	}
	if m.ErrMsg == "" {
		t.Fatalf("expected an error message for empty description")
	}
}

func TestEditAndDeleteFlow(t *testing.T) {
	m := newTestModel(t)

	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	m = typeRunes(t, m, "Original")
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})
	if m.Screen != ScreenEditTask {
		t.Fatalf("expected ScreenEditTask, got %v", m.Screen)
	}
	if m.Input.Value() != "Original" {
		t.Fatalf("expected edit input pre-filled with 'Original', got %q", m.Input.Value())
	}

	// Clear and retype.
	m.Input.SetValue("")
	m = typeRunes(t, m, "Updated")
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	tasks := m.FilteredTasks()
	if len(tasks) != 1 || tasks[0].Description != "Updated" {
		t.Fatalf("expected one task 'Updated', got %+v", tasks)
	}

	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	if m.Screen != ScreenDeleteConfirm {
		t.Fatalf("expected ScreenDeleteConfirm, got %v", m.Screen)
	}

	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	if m.Screen != ScreenTaskList {
		t.Fatalf("expected ScreenTaskList after confirming delete, got %v", m.Screen)
	}
	if len(m.FilteredTasks()) != 0 {
		t.Fatalf("expected no tasks after delete, got %+v", m.FilteredTasks())
	}
}

func TestDeleteCancel(t *testing.T) {
	m := newTestModel(t)
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	m = typeRunes(t, m, "Keep me")
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})

	if m.Screen != ScreenTaskList {
		t.Fatalf("expected ScreenTaskList after cancel, got %v", m.Screen)
	}
	if len(m.FilteredTasks()) != 1 {
		t.Fatalf("expected task to survive cancel, got %+v", m.FilteredTasks())
	}
}

func TestStatusCycling(t *testing.T) {
	m := newTestModel(t)
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	m = typeRunes(t, m, "Cycle me")
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	if got := string(m.FilteredTasks()[0].Status); got != "todo" {
		t.Fatalf("expected initial status todo, got %s", got)
	}

	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeySpace})
	if got := string(m.FilteredTasks()[0].Status); got != "in_progress" {
		t.Fatalf("expected in_progress after one cycle, got %s", got)
	}

	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeySpace})
	if got := string(m.FilteredTasks()[0].Status); got != "done" {
		t.Fatalf("expected done after two cycles, got %s", got)
	}

	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeySpace})
	if got := string(m.FilteredTasks()[0].Status); got != "todo" {
		t.Fatalf("expected wrap to todo after three cycles, got %s", got)
	}
}

func TestFilterCycleAndNavigationClamping(t *testing.T) {
	m := newTestModel(t)

	if m.CurrentFilter() != FilterAll {
		t.Fatalf("expected default filter to be FilterAll, got %s", m.CurrentFilter())
	}

	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	if m.CurrentFilter() != FilterTodo {
		t.Fatalf("expected filter to advance to FilterTodo, got %s", m.CurrentFilter())
	}

	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	if m.CurrentFilter() != FilterAll {
		t.Fatalf("expected filter to go back to FilterAll, got %s", m.CurrentFilter())
	}

	// Cursor must not move past an empty list.
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if m.Cursor != 0 {
		t.Fatalf("expected cursor to stay at 0 on empty list, got %d", m.Cursor)
	}
}

func TestHelpScreenTogglesBack(t *testing.T) {
	m := newTestModel(t)
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	if m.Screen != ScreenHelp {
		t.Fatalf("expected ScreenHelp, got %v", m.Screen)
	}
	m = pressKey(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	if m.Screen != ScreenTaskList {
		t.Fatalf("expected ScreenTaskList after esc, got %v", m.Screen)
	}
}
