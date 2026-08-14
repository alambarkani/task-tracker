package task_tracker

import "testing"

func TestDeleteTask(t *testing.T) {
	err := DeleteTask(2)
	if err != nil {
		t.Error(err)
	}
}
