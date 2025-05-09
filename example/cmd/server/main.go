package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/juliendoutre/protoc-gen-go-mcp/example/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := flag.Int("port", 8000, "port to listen on")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterHelloWorldServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}

type server struct {
	pb.UnimplementedHelloWorldServer
}

func (s *server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	return &pb.GreetResponse{Greeting: fmt.Sprintf("Hello %s!", in.GetName())}, nil
}
