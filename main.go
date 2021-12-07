package main

import (
	"context"
	"log"
	"net"

	pb "example.com/sync/api/pb"
	"example.com/sync/dropboxsdk"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement APIServer
type server struct {
	pb.UnimplementedAPIServer
	// TODO: put cloud config here
}

func (s *server) Ls(ctx context.Context, in *pb.LsRequest) (*pb.LsReply, error) {
	path := in.GetPath()
	dropboxsdk.Ls(path)
	log.Printf("Received: %v", in.GetPath())
	return &pb.LsReply{Lsmessage: "Hello " + in.GetPath(), Lserror: ""}, nil
}

func main() {
	//TODO: cobra cli
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAPIServer(s, &server{})
	// TODO: sdk interface
	dropboxsdk.Execute()
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
