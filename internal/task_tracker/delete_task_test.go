package task_tracker

import (
	"encoding/json"
	"os"
	"path"
	"testing"
	"time"
)

func TestDeleteTask(t *testing.T) {
	dir := t.TempDir()
	path := path.Join(dir, "tasks.json")

	seed := []Task{
		{
			ID:          1,
			Description: "a",
			Status:      Todo,
			CreatedAt:   time.Now().String(),
			UpdatedAt:   time.Now().String(),
		},
		{
			ID:          2,
			Description: "b",
			Status:      Todo,
			CreatedAt:   time.Now().String(),
			UpdatedAt:   time.Now().String(),
		},
	}
	data, err := json.Marshal(seed)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = DeleteTask(path, 2)
	if err != nil {
		t.Error(err)
	}

	var result []Task
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	json.Unmarshal(b, &result)

	for _, task := range result {
		if task.ID == 2 {
			t.Error("Task 2 still exist")
		}
	}
}
