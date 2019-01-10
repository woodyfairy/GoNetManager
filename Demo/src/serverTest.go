package main

import (
	"../../NetManager"
	"log"
)

func main() {
	port := ":3000"

	log.Println("Start Server", port)
	NetManager.Listen(port)
}
