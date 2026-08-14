package task_tracker

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"alambarkani.com/task_tracker/internal/dir_helper"
)

func AddTask(description string) error {
	var newTask Task
	lastTask, err := TaskList()
	if err != nil {
		newTask.ID = 1
	} else {
		newTask.ID = lastTask[len(lastTask)-1].ID + 1
	}

	newTask.Description = description
	newTask.Status = Todo
	newTask.CreatedAt = time.Now().String()
	newTask.UpdatedAt = time.Now().String()

	var tasks []Task

	path := dir_helper.TaskDir()

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
