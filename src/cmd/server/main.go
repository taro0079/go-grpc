package main

import (
	"context"
	"fmt"
	hellogrpc "go-grpc/pkg/grpc"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type myServer struct {
	hellogrpc.UnimplementedGreetingServiceServer
}

func (s *myServer) Hello(ctx context.Context, req *hellogrpc.HelloRequest) (*hellogrpc.HelloResponse, error) {
	return &hellogrpc.HelloResponse{
		Message: fmt.Sprintf("Hello, %s", req.GetName()),
	}, nil
}

func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	hellogrpc.RegisterGreetingServiceServer(s, NewMyServer())

	reflection.Register(s)

	go func() {
		log.Printf("start GRPC server port: %v", port)
		s.Serve(listener)
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping GRPC server ...")
	s.GracefulStop()

}
