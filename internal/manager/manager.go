package manager

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"taskReceiver/internal/config"
	"taskReceiver/internal/tasks"
	"time"
	"uuid"
)

type Manager struct {
	http  *http.Server
	store tasks.TaskStore
	cfg   config.Config
}

func (m *Manager) AddHandlers() {

	// Define routes
	http.HandleFunc("GET /tasks/{id}", m.handleGetTask)
	http.HandleFunc("POST /tasks", m.handleCreateTask)
	http.HandleFunc("DELETE /tasks/{id}", m.handleDeleteTask)
	http.HandleFunc("GET /tasks", m.listTasksHandler)
	http.HandleFunc("GET /healthz", m.healthCheckHandler)
}

func (m *Manager) Start() error {
	slog.Info("Server started", "port", m.cfg.ListenPort)
	if err := m.http.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func NewManager(config *config.Config) *Manager {
	m := &Manager{
		http: &http.Server{
			Addr: fmt.Sprintf("%s:%d", config.ListenHost, config.ListenPort),
		},
		store: tasks.NewStore(config.StoreType),
		cfg:   *config,
	}
	m.AddHandlers()
	m.http.Handler = LoggingMiddleware(http.DefaultServeMux)
	return m
}

func (m *Manager) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (m *Manager) handleGetTask(w http.ResponseWriter, r *http.Request) {
	// Read ID from inline
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid TaskID. Expected: UUID"))
		return
	}
	if id.String() == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("TaskID is required"))
		return
	}

	// We have a real ID
	task, err := m.store.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	payload, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)

}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger := slog.Default().With("method", r.Method, "path", r.URL.Path)
		logger.Info("request started")
		next.ServeHTTP(w, r)
		logger.Info("request completed",
			"latency_ms", time.Since(start).Milliseconds(),
		)
	})
}

func (m *Manager) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	// Read the task from the request body
	var task tasks.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse Descrtiption
	if task.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Description is required"))
		return
	}

	// Generate a new UUID for the task
	task.ID = uuid.New()
	// Set initial Task status
	task.Status = tasks.TaskStatusSubmitted
	// Save the task
	err = m.store.Put(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Serialize task
	payload, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(payload)
	return
}

func (m *Manager) listTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch All tasks
	statusFilter := r.URL.Query().Get("status")
	var taskList []*tasks.Task
	var err error
	if statusFilter == "" {
		// Return Them all. Needs a Pager
		taskList, err = m.store.List()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// Return filtered by status
		taskList, err = m.store.ListByStatus(tasks.TaskStatus(statusFilter))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	payload, err := json.Marshal(taskList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
	return
	
}

func (m *Manager) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	// Read ID from inline
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid TaskID. Expected: UUID"))
		return
	}
	if id.String() == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("TaskID is required"))
		return
	}

	err = m.store.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
