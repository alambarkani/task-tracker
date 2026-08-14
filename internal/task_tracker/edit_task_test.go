package task_tracker

import "testing"

func TestEditTask(t *testing.T) {
	err := EditTask(1, "Test Edit Task 2")
	if err != nil {
		t.Error(err)
	}
}
