package tasks

import "uuid"

type Task struct {
	ID          uuid.UUID  `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
}
