package main

import (
	"log"
	"net"

	"github.com/ianschenck/envflag"
	"github.com/niloy104/Conduit/db"
	"github.com/niloy104/Conduit/grpc/pb"
	"github.com/niloy104/Conduit/grpc/server"
	"github.com/niloy104/Conduit/grpc/storer"
	"google.golang.org/grpc"
)

func main() {

	var (
		svcAddr = envflag.String("SVC_ADDR", "0.0.0.0:9091", "address where the grpc service is listening on")
		dbAddr  = envflag.String("DB_ADDR", "127.0.0.1:3306", "address where the database is running on")
	)
	envflag.Parse()

	//instntiate db
	db, err := db.NewDatabase(*dbAddr)
	if err != nil {
		log.Fatal("error opening database: %w", err)
	}
	defer db.Close()
	log.Println("Successfully connected to the database")

	// do something with the database
	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)

	//register our server with gRPC server

	grpcSrv := grpc.NewServer()
	pb.RegisterEcommServer(grpcSrv, srv)

	listener, err := net.Listen("tcp", *svcAddr)
	if err != nil {
		log.Fatalf("listener failed: %v", err)
	}

	log.Printf("server is listening on %s", *svcAddr)
	err = grpcSrv.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
