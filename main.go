package main

import (
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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"msg": "Hello world"})
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

	// Update a Todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			todoID := strconv.Itoa(todo.ID)
			if todoID == id {
				todos[i].Completed = !todos[i].Completed
				return c.Status(http.StatusOK).JSON(todos[i])
			}
		}

		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	})

	log.Fatal(app.Listen(":4000"))
}
