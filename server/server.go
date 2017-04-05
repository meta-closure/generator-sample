package server

import (
	"fmt"
	"net/http"
)

type Hook struct {
	handler func(http.ResponseWriter, *http.Request) error
}

func (h Hook) Handler(w http.ResponseWriter, r *http.Request) {
	if err := h.handler(w, r); err != nil {
		fmt.Printf("FAILED: %s\n", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	return
}
