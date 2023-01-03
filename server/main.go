package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/geekr-dev/go-grpc-demo/proto"
	"google.golang.org/grpc"
)

type GreeterServer struct {
	*pb.UnimplementedGreeterServer
}

func (gs *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello world"}, nil
}

func (gs *GreeterServer) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
	for n := 0; n <= 6; n++ {
		err := stream.Send(&pb.HelloReply{
			Message: "hello list",
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (gs *GreeterServer) SayRecord(stream pb.Greeter_SayRecordServer) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReply{Message: "say.record"})
		}
		if err != nil {
			return err
		}
		log.Printf("resp: %v", resp)
	}
}

func (gs *GreeterServer) SayRoute(stream pb.Greeter_SayRouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.HelloReply{Message: fmt.Sprintf("say.route%d", n)})
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n++
		log.Printf("resp: %v", resp)
	}
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis, _ := net.Listen("tcp", ":9000")
	server.Serve(lis)
}
