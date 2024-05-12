//go:build e2e

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/iximiuz/todolist/models"
)

func TestCRUD(t *testing.T) {
	todo := createTodo(t, "E2E task")

	todos := listTodos(t)
	found := false
	for _, x := range todos {
		if x.ID == todo.ID {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("todo not found")
	}

	deleteTodo(t, todo.ID)

	todos = listTodos(t)
	found = false
	for _, x := range todos {
		if x.ID == todo.ID {
			found = true
			break
		}
	}

	if found {
		t.Fatal("todo not deleted")
	}
}

func createTodo(t *testing.T, task string) models.Todo {
	resp, err := http.Post("http://localhost:8080/todos", "application/json", bytes.NewBufferString(`{"task": "`+task+`"}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %v", resp.StatusCode)
	}

	var created models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	return created
}

func listTodos(t *testing.T) []models.Todo {
	resp, err := http.Get("http://localhost:8080/todos")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %v", resp.StatusCode)
	}

	var todos []models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		t.Fatal(err)
	}

	return todos
}

func deleteTodo(t *testing.T, id string) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/todos/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204 No Content, got %v", resp.StatusCode)
	}
}
