package models_test

import (
	"encoding/json"
	"testing"

	"github.com/iximiuz/todolist/models"
)

func TestTodo_MarshalJSON(t *testing.T) {
	todo := models.Todo{
		ID:   "123",
		Task: "Test task",
	}

	b, err := json.Marshal(todo)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"id":"123","task":"Test task"}`
	if string(b) != expected {
		t.Errorf("expected %v, got %v", expected, string(b))
	}
}

func TestTodo_UnmarshalJSON(t *testing.T) {
	b := []byte(`{"id":"123","task":"Test task"}`)

	var todo models.Todo
	if err := json.Unmarshal(b, &todo); err != nil {
		t.Fatal(err)
	}

	if todo.ID != "123" {
		t.Errorf("expected %v, got %v", "123", todo.ID)
	}

	if todo.Task != "Test task" {
		t.Errorf("expected %v, got %v", "Test task", todo.Task)
	}
}
