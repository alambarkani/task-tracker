package task_tracker

import (
	"encoding/json"
	"os"
	"path"
	"testing"
	"time"
)

func TestTaskList(t *testing.T) {
	dir := t.TempDir()
	path := path.Join(dir, "tasks.json")

	seed := []Task {
		{
			ID: 1,
			Description: "a",
			Status: Todo,
			CreatedAt: time.Now().String(),
			UpdatedAt: time.Now().String(),
		},
	}
	data, err := json.Marshal(seed)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	
	tasks, err := TaskList(path)
	if err != nil {
		t.Error(err)
	}

	if len(tasks) == 0 {
		t.Error("No tasks found")
	}
}
