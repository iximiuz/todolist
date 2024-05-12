# Todolist - Your TODO list as a service

This toy web application is created specifically for the [iximiuz Labs' Dagger course](https://labs.iximiuz.com/courses/dagger).

To start the application using Docker Compose, run the following command:

```sh
make up
```

Then create a few TODO items:

```sh
curl -X POST -d '{"task": "Finish the course"}' http://localhost:8080/todos
curl -X POST -d '{"task": "Nail the sales!!!"}' http://localhost:8080/todos
```

Retrieve the list of all the tasks:

```sh
curl -X GET http://localhost:8080/todos
```

Delete a task:

```sh
curl -X DELETE http://localhost:8080/todos/{id}
```

Run unit tests:

```sh
make test
```

Run end-to-end tests (requires `make up` to be running):

```sh
make test-e2e
```

Check the [Makefile](Makefile) for all available commands.
