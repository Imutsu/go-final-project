package main

import (
	"log"
	"net/http"
	"os"

	"scheduler/pkg/db"
)

func main() {
	webDir := "web"
	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "7540"
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
