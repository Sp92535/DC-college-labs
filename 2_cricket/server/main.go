package main

import (
	"context"
	pb "cricket/pb"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
)

type CricketServer struct {
	pb.UnimplementedCricketServer
}

func (s *CricketServer) GetTopScorers(ctx context.Context, req *pb.Empty) (*pb.TopScoreResponse, error) {
	return &pb.TopScoreResponse{Name: "Virat Kohli", Average: 60.11}, nil
}

func (s *CricketServer) GetCenturions(ctx context.Context, req *pb.Empty) (*pb.CenturionsResponse, error) {
	return &pb.CenturionsResponse{Name: "Sachin Tendulkar", Centuries: 100}, nil
}

func (s *CricketServer) GetPlayerStats(ctx context.Context, req *pb.PlayerRequest) (*pb.StatsResponse, error) {
	return &pb.StatsResponse{Name: req.Name, Average: 50.23, Centuries: 30}, nil
}

func (s *CricketServer) UpdatePlayerScore(ctx context.Context, req *pb.UpdateScoreRequest) (*pb.Empty, error) {
	log.Printf("Updated score for %s by %d runs", req.Name, req.Runs)
	return &pb.Empty{}, nil
}

// Starting the server
func main() {
	listener, err := net.Listen("tcp", ":6969")

	if err != nil {
		log.Printf("ERROR")
		return
	}

	grpcServer := grpc.NewServer()

	pb.RegisterCricketServer(grpcServer, &CricketServer{})

	fmt.Println("gRPC server running on port 6969...")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
