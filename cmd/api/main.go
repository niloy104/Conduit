package main

import (
	"log"

	"github.com/ianschenck/envflag"
	"github.com/niloy104/Conduit/api/handler"
	"github.com/niloy104/Conduit/api/server"
	"github.com/niloy104/Conduit/api/storer"
	"github.com/niloy104/Conduit/db"
)

const minSecretKeySize = 32

func main() {
	var secretKey = envflag.String("SECRET_kEY", "01234567890123456789012345678901", "secret key for jwt signing")
	if len(*secretKey) < minSecretKeySize {
		log.Fatalf("SECRET_KEY must me at leas %d charachter", minSecretKeySize)
	}
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal("error opening database: %w", err)
	}
	defer db.Close()

	log.Println("Successfully connected to the database")

	// do something with the database
	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv, *secretKey)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")
}
