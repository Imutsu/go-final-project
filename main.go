package main

import (
	"net/http"
	"os"
)

func main() {
	webDir := "web"
	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.ListenAndServe(":"+port, nil)
}
