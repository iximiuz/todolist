package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/iximiuz/todolist/handlers"
	"github.com/iximiuz/todolist/models"
)

func TestCreateTodo(t *testing.T) {
	rdb, close := newRedisClient(t)
	defer close()

	sut := handlers.NewTodoHandler(rdb)

	req, err := http.NewRequest("POST", "/todos", bytes.NewBufferString(`{"task": "Test task"}`))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(sut.CreateTodo).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v: %v", status, http.StatusCreated, rr.Body.String())
	}
}

func TestListTodos(t *testing.T) {
	rdb, close := newRedisClient(t)
	defer close()

	sut := handlers.NewTodoHandler(rdb)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /todos", sut.CreateTodo)
	mux.HandleFunc("GET /todos", sut.ListTodos)

	req, err := http.NewRequest("POST", "/todos", bytes.NewBufferString(`{"task": "Test task"}`))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	req, err = http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v: %v", status, http.StatusOK, rr.Body.String())
	}

	var todos []models.Todo
	if err := json.NewDecoder(rr.Body).Decode(&todos); err != nil {
		t.Fatal(err)
	}

	if len(todos) != 1 {
		t.Errorf("handler returned wrong number of todos: got %v want %v", len(todos), 1)
	}
}

func TestDeleteTodo(t *testing.T) {
	rdb, close := newRedisClient(t)
	defer close()

	sut := handlers.NewTodoHandler(rdb)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /todos", sut.CreateTodo)
	mux.HandleFunc("GET /todos", sut.ListTodos)
	mux.HandleFunc("DELETE /todos/{id}", sut.DeleteTodo)

	req, err := http.NewRequest("POST", "/todos", bytes.NewBufferString(`{"task": "Test task"}`))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	var todo models.Todo
	if err := json.NewDecoder(rr.Body).Decode(&todo); err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("DELETE", "/todos/"+todo.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v: %v", status, http.StatusNoContent, rr.Body.String())
	}

	req, err = http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v: %v", status, http.StatusOK, rr.Body.String())
	}

	var todos []models.Todo
	if err := json.NewDecoder(rr.Body).Decode(&todos); err != nil {
		t.Fatal(err)
	}

	if len(todos) != 0 {
		t.Errorf("handler returned wrong number of todos: got %v want %v", len(todos), 0)
	}
}

func newRedisClient(t *testing.T) (*redis.Client, func()) {
	r := miniredis.RunT(t)

	return redis.NewClient(&redis.Options{
		Addr: r.Addr(),
	}), r.Close
}
