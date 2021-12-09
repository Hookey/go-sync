package api

import (
	"context"
	"log"

	pb "github.com/Hookey/go-sync/api/pb"
	"github.com/Hookey/go-sync/core"
)

// server is used to implement APIServer
type Service struct {
	pb.UnimplementedAPIServer
	core.Storage
}

func (s *Service) Ls(ctx context.Context, in *pb.LsRequest) (*pb.LsReply, error) {
	path := in.GetPath()
	s.Storage.Ls(path)
	log.Printf("Received: %v", in.GetPath())
	return &pb.LsReply{Lsmessage: "Hello " + in.GetPath(), Lserror: ""}, nil
}
