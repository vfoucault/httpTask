package tasks

import (
	"uuid"
)

type TaskStore interface {
	Get(ID uuid.UUID) (*Task, error)
	Put(task *Task) error
	List() ([]*Task, error)
	ListByStatus(status TaskStatus) ([]*Task, error)
	Delete(ID uuid.UUID) error
}

func NewStore(storeType string) TaskStore {
	switch storeType {
	case "memory":
		return &MemoryTaskStore{
			tasks: make(map[uuid.UUID]*Task),
		}
	default:
		return nil
	}
}
