package main

import (
	"example.com/project-rest/internal/config"
	"example.com/project-rest/internal/transport"
	"log"
)

func main() {
	c := config.Load()
	handler, err := transport.Gateway(c.UserAddress, c.HRAddress)
	if err != nil {
		log.Fatal(err)
	}
	if err := transport.Serve(c.GatewayPort, handler); err != nil {
		log.Fatal(err)
	}
}
