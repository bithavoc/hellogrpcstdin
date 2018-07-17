//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"sync"

	"github.com/bithavoc/hellogrpcstdin/common"
	pb "github.com/bithavoc/hellogrpcstdin/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println("Server replying, check greeter_server/client.log")
	return &pb.HelloReply{Message: "Hello " + in.Name + " from server"}, nil
}

func main() {
	fmt.Println("Server starting")
	cmd := exec.Command("../greeter_client/greeter_client")
	inPipeReader, inPipeWriter := io.Pipe()
	outPipeReader, outPipeWriter := io.Pipe()

	cmd.Stdin = inPipeReader
	cmd.Stdout = outPipeWriter

	pipe := common.NewStdStreamJoint(outPipeReader, inPipeWriter)
	lis := &stdinListener{}
	go startServer(lis)

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("agent process running, will ready listener")
	lis.Ready(pipe)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Pluggin stopped")
}

func startServer(lis net.Listener) {
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type stdinListener struct {
	closed   bool
	wait     sync.WaitGroup
	onlyConn net.Conn
}

func (lis *stdinListener) Ready(conn net.Conn) {
	lis.onlyConn = conn
	lis.wait.Done()
}

func (lis *stdinListener) Accept() (net.Conn, error) {
	lis.wait.Add(1)
	fmt.Println("accepting, waiting for only conn")
	lis.wait.Wait()
	fmt.Println("accepting, only conn ready")
	return lis.onlyConn, nil
}

func (lis *stdinListener) Close() error {
	lis.closed = true
	return nil
}

func (lis *stdinListener) Addr() net.Addr {
	return common.NewStdinAddr("listener")
}
