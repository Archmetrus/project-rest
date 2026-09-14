package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"

	"example.com/project-rest/internal/config"
	"example.com/project-rest/internal/store"
	"example.com/project-rest/migrations"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) != 3 {
		return fmt.Errorf("usage: migrate <user|hr> <up|status|down>")
	}
	service, action := os.Args[1], os.Args[2]
	c := config.Load()
	path := c.UserDB
	switch service {
	case "user":
	case "hr":
		path = c.HRDB
	default:
		return fmt.Errorf("unknown service: %s", service)
	}
	if action != "up" && action != "status" && action != "down" {
		return fmt.Errorf("unknown migration action: %s", action)
	}
	db, err := store.Open(path)
	if err != nil {
		return err
	}
	defer db.Close()
	files, err := fs.Sub(migrations.Files, service)
	if err != nil {
		return err
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, db, files)
	if err != nil {
		return err
	}
	ctx := context.Background()
	switch action {
	case "up":
		results, err := provider.Up(ctx)
		if err != nil {
			return err
		}
		for _, result := range results {
			fmt.Printf("applied %s\n", result.Source.Path)
		}
	case "down":
		result, err := provider.Down(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("rolled back %s\n", result.Source.Path)
	case "status":
		results, err := provider.Status(ctx)
		if err != nil {
			return err
		}
		for _, result := range results {
			fmt.Printf("%s %s\n", result.Source.Path, result.State)
		}
	}
	return nil
}
