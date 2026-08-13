package task_tracker

import (
	"encoding/json"
	"fmt"
	"os"

	"alambarkani.com/task_tracker/internal/dir_helper"
)

func TaskList() ([]Task, error) {
	taskFile, err := os.Open(dir_helper.TaskDir())
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %s", err)
	}
	defer taskFile.Close()

	var tasks []Task
	decoder := json.NewDecoder(taskFile)
	err = decoder.Decode(&tasks)
	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON: %s", err)
	}

	return tasks, nil
}
