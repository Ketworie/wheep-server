package main

import (
	"log"
	"wheep-server/server"
)

func main() {
	log.Fatal(server.StartServer())
}
