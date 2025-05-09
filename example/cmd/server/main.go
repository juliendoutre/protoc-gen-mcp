package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

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

	randomServer := &pb.HelloWorldRandomServer{
		Rand: rand.New(rand.NewSource(time.Now().Unix())),
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterHelloWorldServer(grpcServer, randomServer)
	grpcServer.Serve(lis)
}
