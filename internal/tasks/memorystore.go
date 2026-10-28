package tasks

import (
	"fmt"
	"sync"
	"uuid"
)

type MemoryTaskStore struct {
	tasks map[uuid.UUID]*Task
	mutex sync.RWMutex
}

func (m *MemoryTaskStore) ListByStatus(status TaskStatus) ([]*Task, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	tasks := []*Task{}
	for _, task := range m.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (m *MemoryTaskStore) Get(ID uuid.UUID) (*Task, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if val, ok := m.tasks[ID]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("task with ID %s not found", ID.String())
}

func (m *MemoryTaskStore) Put(task *Task) error {
	// Overwrite allowed
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.tasks[task.ID] = task
	return nil
}

func (m *MemoryTaskStore) Delete(ID uuid.UUID) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.tasks[ID]; ok {
		delete(m.tasks, ID)
		return nil
	}
	return fmt.Errorf("task with ID %s not found", ID.String())
}

func (m *MemoryTaskStore) List() ([]*Task, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	tasks := make([]*Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}
