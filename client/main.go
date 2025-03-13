package main

import (
	"context"
	"log"
	"time"

	pb "grpc-streaming/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func LoopPing(ctx context.Context, client pb.MessagesClient, interval time.Duration) error {
	stream, err := client.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to call Ping: %v", err)
	}
	defer stream.CloseSend()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := stream.Send(&pb.PingRequest{Message: "ping"}); err != nil {
			return err
		}

		res, err := stream.Recv()
		if err != nil {
			return err
		}

		log.Printf("received: %s", res.Message)
	}

	return nil
}

func main() {
	ctx := context.Background()
	interval := 1 * time.Second

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewMessagesClient(conn)

	if err := LoopPing(ctx, client, interval); err != nil {
		log.Fatalf("failed to loop ping: %v", err)
	}
}
