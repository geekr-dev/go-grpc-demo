package main

import (
	"context"
	"io"
	"log"

	pb "github.com/geekr-dev/go-grpc-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, _ := grpc.Dial(
		":9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	err := SayHello(client)
	if err != nil {
		log.Fatalf("SayHello err: %v", err)
	}
	err = SayList(client, &pb.HelloRequest{Name: "geekr1"})
	if err != nil {
		log.Fatalf("SayList err: %v", err)
	}
	err = SayRecord(client, &pb.HelloRequest{Name: "geekr2"})
	if err != nil {
		log.Fatalf("SayRecord err: %v", err)
	}
	err = SayRoute(client, &pb.HelloRequest{Name: "geekr3"})
	if err != nil {
		log.Fatalf("SayRoute err: %v", err)
	}
}

func SayHello(client pb.GreeterClient) error {
	resp, err := client.SayHello(
		context.Background(),
		&pb.HelloRequest{Name: "geekr"},
	)
	if err != nil {
		return err
	}
	log.Printf("client.SayHello resp: %s", resp.Message)
	return nil
}

func SayList(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayList(context.Background(), r)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: %v", resp)
	}
	return nil
}

func SayRecord(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayRecord(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n <= 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("resp: %v", resp)
	return nil
}

func SayRoute(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayRoute(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n <= 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return nil
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: %v", resp)
	}

	err = stream.CloseSend()
	if err != nil {
		return err
	}
	return nil
}
