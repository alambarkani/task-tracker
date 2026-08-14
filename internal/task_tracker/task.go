package task_tracker

type Task struct {
	ID          int
	Description string
	Status      Status
	CreatedAt   string
	UpdatedAt   string
}

type Status string

const (
	Todo       Status = "todo"
	InProgress Status = "in_progress"
	Done       Status = "done"
)
