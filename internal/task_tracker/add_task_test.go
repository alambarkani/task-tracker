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