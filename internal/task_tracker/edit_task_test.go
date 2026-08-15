package task_tracker

import (
	"encoding/json"
	"os"
	"path"
	"testing"
	"time"
)

func TestEditTask(t *testing.T) {
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
	}

	data, err := json.Marshal(&seed)
	if err = os.WriteFile(path, data, 0644); err != nil {
		t.Fatal("Cannot encode data to json")
	}

	err = EditTask(path, 1, "Test Edit")
	if err != nil {
		t.Error(err)
	}

	var tasks []Task
	fileData, err := os.ReadFile(path)
	if err != nil {
		t.Fatal("Cannot read file")
	}
	if err = json.Unmarshal(fileData, &tasks); err != nil {
		t.Fatal("Cannot decode file")
	}

	for _, data := range tasks {
		if data.ID == 1 {
			if data.Description != "Test Edit" {
				t.Error("Data 1 Not Edited")
			}
		}
	}
}
