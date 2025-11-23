package domain

type TaskListOut struct {
	Tasks      []Task
	TotalTasks int64
	TotalPages int64
}
