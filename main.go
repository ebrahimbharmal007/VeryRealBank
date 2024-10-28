package main

import "log"

func main() {

	store, err := NewInMemoryRepository()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	service, err := NewService(store)
	if err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer("3000", *service)
	server.Run()
}
