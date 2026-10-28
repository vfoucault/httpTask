package tasks

type TaskStatus string

var (
	TaskStatusPending    = TaskStatus("pending")
	TaskStatusInProgress = TaskStatus("in_progress")
	TaskStatusCompleted  = TaskStatus("completed")
	TaskStatusSubmitted  = TaskStatus("submitted")
)
