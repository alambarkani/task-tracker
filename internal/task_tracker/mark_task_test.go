package task_tracker

import (
	"encoding/json"
	"os"
	"path"
	"testing"
	"time"
)

func TestMarkTask(t *testing.T) {
	dir := t.TempDir()
	path := path.Join(dir, "tasks.json")

	mockTask := []Task{
		{
			ID:          1,
			Description: "Task 1",
			Status:      Todo,
			CreatedAt:   time.Now().GoString(),
			UpdatedAt:   time.Now().GoString(),
		},
		{
			ID:          2,
			Description: "Task 2",
			Status:      Todo,
			CreatedAt:   time.Now().GoString(),
			UpdatedAt:   time.Now().GoString(),
		},
	}

	data, _ := json.Marshal(&mockTask)
	_ = os.WriteFile(path, data, 0644)

	err := MarkTask(path, 1, InProgress)
	if err != nil {
		t.Fatal("Mark fail")
	}

	var afterMark []Task
	fileAfterMark, _ := os.ReadFile(path)
	_ = json.Unmarshal(fileAfterMark, &afterMark)

	if afterMark[0].Status != InProgress {
		t.Fatal("Status not changed")
	}
}
