package api

import (
	"context"

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
	err := s.Storage.Ls(path)
	//log.Printf("Received: %v", in.GetPath())
	return &pb.LsReply{Result: "Hello " + in.GetPath()}, err
}

func (s *Service) Put(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
	src := in.GetSrc()
	dst := in.GetDst()
	//workers := in.GetWorkers()
	//chunksize := in.GetChunksize()
	err := s.Storage.Put(src, dst)
	//log.Printf("Received: %v", in.GetPath())
	return &pb.PutReply{}, err
}
