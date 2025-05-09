package main

import (
	"flag"
	"log"

	"github.com/juliendoutre/protoc-gen-go-mcp/example/internal/pb"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/grpc"
)

func main() {
	target := flag.String("target", "localhost:8000", "gRPC service target")
	flag.Parse()

	conn, err := grpc.NewClient(*target)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.ServeStdio(pb.NewHelloWorldMCPServer(pb.NewHelloWorldClient(conn))); err != nil {
		log.Fatal(err)
	}
}
