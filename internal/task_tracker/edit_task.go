package task_tracker

import (
	"encoding/json"
	"fmt"
	"os"
)

func EditTask(path string, id int, description string) error {

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Error opening file: %s", err)
	}
	defer file.Close()

	var tasks []Task
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		return fmt.Errorf("Error decoding JSON: %s", err)
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = description
			break
		}
	}

	fileData, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}
	defer fileData.Close()

	updatedTasks := json.NewEncoder(fileData)
	err = updatedTasks.Encode(tasks)
	if err != nil {
		return fmt.Errorf("Error encoding JSON: %s", err)
	}

	return nil
}
