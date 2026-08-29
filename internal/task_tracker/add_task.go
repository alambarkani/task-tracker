package task_tracker

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func AddTask(path string, description string) error {
	var tasks []Task

	fileData, err := os.Open(path)
	switch {
	case os.IsNotExist(err):
		break
	case err != nil:
		return fmt.Errorf("Error opening file: %s", err)
	default:
		defer fileData.Close()
		decoder := json.NewDecoder(fileData)
		err = decoder.Decode(&tasks)
		if err != nil {
			return fmt.Errorf("Error decoding JSON: %s", err)
		}
	}

	var newTask Task
	if len(tasks) == 0 {
		newTask.ID = 1
	} else {
		newTask.ID = tasks[len(tasks)-1].ID + 1
	}

	newTask.Description = description
	newTask.Status = Todo
	newTask.CreatedAt = time.Now().String()
	newTask.UpdatedAt = time.Now().String()

	tasks = append(tasks, newTask)

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}
	defer file.Close()

	updatedTasks := json.NewEncoder(file)
	err = updatedTasks.Encode(tasks)
	if err != nil {
		return fmt.Errorf("Error encoding JSON: %s", err)
	}

	return nil
}
