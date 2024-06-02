package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/kmulqueen/go-react-todo/models"
)

type todoController struct {
	todoIDPattern *regexp.Regexp
}

func (tc todoController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/todos" {
		switch r.Method {
		case http.MethodGet:
			tc.getTodos(w, r)
		case http.MethodPost:
			tc.createNewTodo(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		matches := tc.todoIDPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		switch r.Method {
		case http.MethodGet:
			tc.getTodoByID(id, w)
		case http.MethodPatch:
			tc.updateTodoByID(id, w, r)
		case http.MethodDelete:
			tc.deleteTodoByID(id, w)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

// List all Todos
func (tc *todoController) getTodos(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetTodos(), w)
}

// Get Todo by ID
func (tc *todoController) getTodoByID(id int, w http.ResponseWriter) {
	t, err := models.GetTodoByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(t, w)
}

// Create new Todo
func (tc *todoController) createNewTodo(w http.ResponseWriter, r *http.Request) {
	t, err := tc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse Todo object"))
		return
	}
	t, err = models.AddTodo(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(t, w)
}

// Update Todo by ID
func (tc *todoController) updateTodoByID(id int, w http.ResponseWriter, r *http.Request) {
	t, err := tc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse Todo object"))
		return
	}
	if id != t.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted Todo must match the ID in the URL"))
		return
	}
	t, err = models.UpdateTodo(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(t, w)
}

// Delete Todo by ID
func (tc *todoController) deleteTodoByID(id int, w http.ResponseWriter) {
	err := models.RemoveTodoByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (tc *todoController) parseRequest(r *http.Request) (models.Todo, error) {
	dec := json.NewDecoder(r.Body)
	var t models.Todo
	err := dec.Decode(&t)
	if err != nil {
		return models.Todo{}, err
	}
	return t, nil
}

func newTodoController() *todoController {
	return &todoController{
		todoIDPattern: regexp.MustCompile(`^/users/(\d+)/?`),
	}
}
