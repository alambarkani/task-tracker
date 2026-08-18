package task_tracker

import (
	"encoding/json"
	"fmt"
	"os"
)

func MarkTask(path string, id int, markStatus Status) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to open file %s", err)
	}
	defer file.Close()

	var tasks []Task
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		return fmt.Errorf("Error decoding file %s", err)
	}

	var foundData bool
	for i, task := range tasks {
		if task.ID == id {
			foundData = true
			tasks[i].Status = markStatus
			break
		}
	}

	if !foundData {
		return fmt.Errorf("Data with ID=%d was not found", id)
	}

	fileWrite, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error to write file: %s", err)
	}
	defer fileWrite.Close()
	encoder := json.NewEncoder(fileWrite)
	err = encoder.Encode(&tasks)
	if err != nil {
		return fmt.Errorf("Error to encode file: %s", err)
	}

	return nil
}