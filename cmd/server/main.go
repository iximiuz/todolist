package main

import (
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"

	"github.com/iximiuz/todolist/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is not set")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	h := handlers.NewTodoHandler(rdb)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /todos", h.CreateTodo)
	mux.HandleFunc("GET /todos", h.ListTodos)
	mux.HandleFunc("DELETE /todos/{id}", h.DeleteTodo)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
