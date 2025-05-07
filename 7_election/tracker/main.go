package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Sp92535/pb"
	"google.golang.org/grpc"
)

type TrackerServer struct {
	pb.UnimplementedTrackerServer
	nodes map[int64]string
}

func (t *TrackerServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.Empty, error) {
	t.nodes[req.Pid] = req.Address
	return &pb.Empty{}, nil
}

func (t *TrackerServer) GetNodes(ctx context.Context, req *pb.Empty) (*pb.GetNodesResponse, error) {
	return &pb.GetNodesResponse{
		Address: t.nodes,
	}, nil
}

func main() {

	listner, err := net.Listen("tcp", "localhost:6969")

	if err != nil {
		log.Fatalf("listner error %v", err)
	}

	grpcServer := grpc.NewServer()

	trackerServer := &TrackerServer{
		nodes: make(map[int64]string),
	}

	pb.RegisterTrackerServer(grpcServer, trackerServer)

	log.Printf("grpc tracker server running on %s...", listner.Addr().String())
	if err = grpcServer.Serve(listner); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}

}
