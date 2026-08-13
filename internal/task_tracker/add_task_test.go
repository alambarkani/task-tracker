package task_tracker

import "testing"

func TestAddTask(t *testing.T) {
	err := AddTask("Test Task 2")
	if err != nil {
		t.Error(err)
	}
}
