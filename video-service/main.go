package main

import (
	"fmt"
	"log"
	"net"

	proto "github.com/AJC232/InfinityStream-backend/common/protoc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterVideoServiceServer(srv, &Video{})
	reflection.Register(srv)

	fmt.Println("Video Service running on :8082")
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
