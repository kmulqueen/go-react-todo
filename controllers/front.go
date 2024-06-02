package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

func RegisterControllers() {
	tc := newTodoController()

	http.Handle("/todos", *tc)
	http.Handle("/todos/", *tc)

}

func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
