package main

import (
	"GoCacher/internal"
	"log"
)

func main() {
	server := internal.NewServer(internal.Config{})
	log.Fatal(server.Start())
}
