package main

import (
	"log"
	"net"

	pb "grpc-streaming/model"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedMessagesServer
}

func (s *Server) Ping(stream pb.Messages_PingServer) error {
	for {
		res, err := stream.Recv()
		if err != nil {
			return err
		}

		log.Printf("received: %s", res.Message)

		if err := stream.Send(&pb.PingResponse{Message: "pong"}); err != nil {
			return err
		}
	}
}

func NewServer() pb.MessagesServer {
	return &Server{}
}

func main() {
	tcpServer, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	messageServer := NewServer()
	pb.RegisterMessagesServer(grpcServer, messageServer)

	if err := grpcServer.Serve(tcpServer); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
