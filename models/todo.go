package models

import (
	"errors"
	"fmt"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

var (
	todos  []*Todo
	nextID = 1
)

func GetTodos() []*Todo {
	return todos
}

func AddTodo(t Todo) (Todo, error) {
	if t.ID != 0 {
		return Todo{}, errors.New("new todo must not include an ID, or it must be set to 0")
	}
	t.ID = nextID
	nextID++
	todos = append(todos, &t)
	return t, nil
}

func GetTodoByID(id int) (Todo, error) {
	for _, t := range todos {
		if t.ID == id {
			return *t, nil
		}
	}

	return Todo{}, fmt.Errorf("Todo with ID '%v' not found", id)
}

func UpdateTodo(t Todo) (Todo, error) {
	for i, item := range todos {
		if item.ID == t.ID {
			todos[i] = &t
			return t, nil
		}
	}

	return Todo{}, fmt.Errorf("Todo with ID '%v' not found", t.ID)
}

func RemoveTodoByID(id int) error {
	for i, item := range todos {
		if item.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Todo with ID '%v' not found", id)
}
