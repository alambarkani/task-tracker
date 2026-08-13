package task_tracker

import "testing"

func TestTaskList(t *testing.T) {
	tasks, err := TaskList()
	if err != nil {
		t.Error(err)
	}

	if len(tasks) == 0 {
		t.Error("No tasks found")
	}
}