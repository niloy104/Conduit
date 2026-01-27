package main

import (
	"log"

	"github.com/ianschenck/envflag"
	"github.com/niloy104/Conduit/api/handler"
	"github.com/niloy104/Conduit/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const minSecretKeySize = 32

func main() {
	var (
		secretKey = envflag.String("SECRET_kEY", "01234567890123456789012345678901", "secret key for jwt signing")
		svcAddr   = envflag.String("GRPC_SVC_ADDR", "0.0.0.0:9091", "address where the grpc service is listening on")
	)

	if len(*secretKey) < minSecretKeySize {
		log.Fatalf("SECRET_KEY must me at leas %d charachter", minSecretKeySize)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(*svcAddr, opts...)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewEcommClient(conn)

	hdl := handler.NewHandler(client, *secretKey)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")
}
