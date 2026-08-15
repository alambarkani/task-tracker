package task_tracker

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

func DeleteTask(path string, id int) error {
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

	var found bool
	for _, data := range tasks {
		if data.ID == id {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Task doesn't exist")
	}

	tasks = slices.DeleteFunc(tasks, func(t Task) bool {
		return t.ID == id
	})

	newFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}
	defer newFile.Close()

	encoder := json.NewEncoder(newFile)
	err = encoder.Encode(&tasks)
	if err != nil {
		return fmt.Errorf("Error encoding file: %s", err)
	}

	return nil
}
