package main

import (
	"log"
	"net/http"
)

func handleIndex(log *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
		},
	)
}
