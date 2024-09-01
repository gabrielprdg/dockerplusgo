package main

import "log"

func main() {
	storage, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	if err := storage.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(":8080", storage)
	server.Run()
}
