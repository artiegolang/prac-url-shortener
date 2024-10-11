package main

import (
	"net/http"
	"practicum-middle/pkg/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.RootHandler)

	if err := http.ListenAndServe(":8085", mux); err != nil {
		panic(err)
	}
}
