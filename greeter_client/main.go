package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/bithavoc/hellogrpcstdin/common"
	pb "github.com/bithavoc/hellogrpcstdin/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	defaultName = "world"
)

var logger *log.Logger

func init() {
	debugFile, err := os.Create("client.log")
	if err != nil {
		panic(err)
	}
	logger = log.New(debugFile, "", log.LstdFlags)
}

func main() {
	logger.Println("client starting")
	pipe := common.NewStdStreamJoint(os.Stdin, os.Stdout)
	// Set up a connection to the server.
	conn, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		logger.Println("addresss", addr)
		return pipe, nil
	}))
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	name := defaultName
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		logger.Fatalf("could not greet: %v", err)
	}
	logger.Printf("Greeting: %s", r.Message)
}
