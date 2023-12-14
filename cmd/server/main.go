package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	hellopb "github.com/sokoide/go-grpc-roundrobin/pkg/grpc"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	msg := fmt.Sprintf("Hello, %s!", req.GetName())
	log.Printf("Hello: returning %s\n", msg)
	return &hellopb.HelloResponse{
		Message: msg,
	}, nil
}

func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	log.Infoln("server")
	port := 8080

	// listen
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// create a server
	s := grpc.NewServer()
	hellopb.RegisterGreetingServiceServer(s, NewMyServer())

	// start a server
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// ctrl-c
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
