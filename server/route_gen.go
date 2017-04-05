package server

import (
	"net/http"
)

func Run() error {
	http.HandleFunc("/user", Hook{handler: PostUserHandler}.Handler)
	return http.ListenAndServe(":8080", nil)
}
