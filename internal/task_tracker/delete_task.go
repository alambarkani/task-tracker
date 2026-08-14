package task_tracker

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"alambarkani.com/task_tracker/internal/dir_helper"
)

func DeleteTask(id int) error {
	path := dir_helper.TaskDir()

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Error opening file: %s", err)
	}
	defer file.Close()

	var tasks []Task
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		return fmt.Errorf("Error decoding file: %s", err)
	}

	for i, data := range tasks {
		if data.ID == id {
			tasks = slices.DeleteFunc(tasks, func(t Task) bool {
				return t.ID == id
			})
			break
		}

		if i >= len(tasks)-1 {
			return fmt.Errorf("Task doesn't exist")
		}
	}

	newFile, err := os.Create(path)
	encoder := json.NewEncoder(newFile)
	err = encoder.Encode(&tasks)
	if err != nil {
		return fmt.Errorf("Error encoding file: %s", err)
	}
	defer newFile.Close()

	return nil
}
