package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/iximiuz/todolist/models"
)

type TodoHandler struct {
	rdb *redis.Client
}

func NewTodoHandler(rdb *redis.Client) *TodoHandler {
	return &TodoHandler{rdb: rdb}
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.ID = uuid.New().String()

	if err := h.rdb.HSet(r.Context(), "todos", todo.ID, todo.Task).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	todosMap, err := h.rdb.HGetAll(r.Context(), "todos").Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todos := []models.Todo{}
	for id, task := range todosMap {
		todos = append(todos, models.Todo{ID: id, Task: task})
	}

	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")
	if taskID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if err := h.rdb.HDel(r.Context(), "todos", taskID).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
