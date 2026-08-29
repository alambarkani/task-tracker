package task_tracker

import (
	"encoding/json"
	"os"
	"path"
	"testing"
)

func TestAddTask(t *testing.T) {
	dir := t.TempDir()
	path := path.Join(dir, "tasks.json")
	err := AddTask(path, "Test Task")
	if err != nil {
		t.Error(err)
	}

	var tasks []Task
	file, err := os.ReadFile(path)
	if err != nil {
		t.Fatal("Cannot read file")
	}
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		t.Fatal("Cannot decode file")
	}

	if len(tasks) <= 0 {
		t.Error("Test Task not exist after added")
	}
}

// TestAddTaskAfterListEmptied guards against a regression where AddTask
// panicked with an index-out-of-range once the task file existed but held
// zero tasks (e.g. right after deleting the only task).
func TestAddTaskAfterListEmptied(t *testing.T) {
	dir := t.TempDir()
	taskPath := path.Join(dir, "tasks.json")

	if err := AddTask(taskPath, "first"); err != nil {
		t.Fatalf("add first task: %v", err)
	}
	if err := DeleteTask(taskPath, 1); err != nil {
		t.Fatalf("delete first task: %v", err)
	}
	if err := AddTask(taskPath, "second"); err != nil {
		t.Fatalf("add task after list emptied: %v", err)
	}

	tasks, err := TaskList(taskPath)
	if err != nil {
		t.Fatalf("list tasks: %v", err)
	}
	// IDs restart from 1 once the list is empty again, matching the
	// behavior of adding to a task file that was never created.
	if len(tasks) != 1 || tasks[0].Description != "second" || tasks[0].ID != 1 {
		t.Fatalf("expected a single task {ID:1, Description:second}, got %+v", tasks)
	}
}