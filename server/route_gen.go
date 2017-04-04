package server

import "net/http"

func Run() error {
	http.HandleFunc("/user", Hook{handler: postUserHandler}.Handler)
	return http.ListenAndServe(":8080", nil)
}
