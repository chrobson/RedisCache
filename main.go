package main

import (
	"log"

	_ "github.com/lib/pq"
)

func main() {
	database, err := NewPostgresDatabase()
	if err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":8080", database)
	server.Run()
}
