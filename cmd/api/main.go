package main

import (
	"log"

	"github.com/niloy104/Conduit/db"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal("error opening database: %w", err)
	}
	defer db.Close()

	log.Println("Successfully connected to the database")
}


// do something with the database
// st:= storer.NewStorer(db.GetDB())