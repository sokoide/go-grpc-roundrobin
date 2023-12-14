package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	hellopb "github.com/sokoide/go-grpc-roundrobin/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type options struct {
	host       string
	port       int
	calls      int
	roundRobin bool
}

var (
	client hellopb.GreetingServiceClient
	o      options
)

func parseArgs() {
	flag.StringVar(&o.host, "host", "localhost", "target host")
	flag.IntVar(&o.port, "port", 8080, "target port")
	flag.IntVar(&o.calls, "calls", 10, "number of calls")
	flag.BoolVar(&o.roundRobin, "roundrobin", false, "roundrobin")
	flag.Parse()
}

func getConn(endpoint string, roundRobin bool) (conn *grpc.ClientConn, err error) {
	resolver.SetDefaultScheme("dns")

	// https://github.com/grpc/grpc/blob/master/doc/load-balancing.md
	if roundRobin {
		conn, err = grpc.Dial(endpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		)
	} else {
		conn, err = grpc.Dial(endpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	}
	return conn, err
}

func main() {
	parseArgs()
	log.Println("start gRPC Client.")

	address := fmt.Sprintf("%s:%d", o.host, o.port)
	log.Printf("target=%s", address)

	conn, err := getConn(address, o.roundRobin)

	if err != nil {
		log.Fatal("Connection failed.")
		return
	}
	defer conn.Close()

	// create a client
	client = hellopb.NewGreetingServiceClient(conn)

	// make calls
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

loop:
	for {
		select {
		case <-ticker.C:
			Hello()
		case <-interrupt:
			break loop
		}
	}

}

func Hello() {
	req := &hellopb.HelloRequest{
		Name: "Hoge",
	}
	res, err := client.Hello(context.Background(), req)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(res.GetMessage())
	}
}
