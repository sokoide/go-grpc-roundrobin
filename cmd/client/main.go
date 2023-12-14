package main

import (
	"bufio"
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	hellopb "github.com/sokoide/go-grpc-roundrobin/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	scanner *bufio.Scanner
	client  hellopb.GreetingServiceClient
)

func main() {
	log.Println("start gRPC Client.")

	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,

		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("Connection failed.")
		return
	}
	defer conn.Close()

	// create a client
	client = hellopb.NewGreetingServiceClient(conn)
	// make a call
	Hello()
}

func Hello() {
	log.Println("Please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}
	res, err := client.Hello(context.Background(), req)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(res.GetMessage())
	}
}
