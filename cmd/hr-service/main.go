package main

import (
	"example.com/project-rest/internal/config"
	"example.com/project-rest/internal/store"
	"example.com/project-rest/internal/transport"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	c := config.Load()
	db, err := store.Open(c.HRDB)
	if err != nil {
		return err
	}
	defer db.Close()
	return transport.Serve(c.HRPort, transport.LeaveHandler(store.LeaveStore{DB: db}))
}
