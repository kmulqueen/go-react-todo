package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()

	todos := []Todo{}

	// List all Todos
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(todos)
	})

	// Create new Todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Todo body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(http.StatusCreated).JSON(todo)
	})

	// Get Todo by ID
	app.Get("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for _, todo := range todos {
			if strconv.Itoa(todo.ID) == id {
				return c.Status(http.StatusOK).JSON(todo)
			}
		}

		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	})

	// Update a Todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if strconv.Itoa(todo.ID) == id {
				todos[i].Completed = !todos[i].Completed
				return c.Status(http.StatusOK).JSON(todos[i])
			}
		}

		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	})

	// Delete a Todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if strconv.Itoa(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(http.StatusOK).JSON(fiber.Map{"msg": fmt.Sprint(todo.Body) + " was deleted"})
			}
		}

		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	})

	log.Fatal(app.Listen(":4000"))
}
