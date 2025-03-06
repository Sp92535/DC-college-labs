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
	players map[string]*pb.StatsResponse
}

func NewCricketServer() *CricketServer {
	return &CricketServer{
		players: map[string]*pb.StatsResponse{
			"Virat Kohli":      {Name: "Virat Kohli", Average: 57.3, Centuries: 76},
			"Sachin Tendulkar": {Name: "Sachin Tendulkar", Average: 53.8, Centuries: 100},
			"Rohit Sharma":         {Name: "Rohit Sharma", Average: 50.6, Centuries: 49},
		},
	}
}

func (s *CricketServer) GetTopScorers(ctx context.Context, req *pb.Empty) (*pb.TopScoreResponse, error) {
	topScorer := ""
	highestAvg := 0.0

	for name, stats := range s.players {
		if stats.Average > highestAvg {
			topScorer = name
			highestAvg = stats.Average
		}
	}

	return &pb.TopScoreResponse{Name: topScorer, Average: highestAvg}, nil
}

func (s *CricketServer) GetCenturions(ctx context.Context, req *pb.Empty) (*pb.CenturionsResponse, error) {
	topCenturion := ""
	var highestCenturies uint32

	for name, stats := range s.players {
		if stats.Centuries > highestCenturies {
			topCenturion = name
			highestCenturies = stats.Centuries
		}
	}

	return &pb.CenturionsResponse{Name: topCenturion, Centuries: highestCenturies}, nil
}

func (s *CricketServer) GetPlayerStats(ctx context.Context, req *pb.PlayerRequest) (*pb.StatsResponse, error) {
	if stats, exists := s.players[req.Name]; exists {
		return stats, nil
	}
	return nil, fmt.Errorf("player not found")
}

func (s *CricketServer) UpdatePlayerScore(ctx context.Context, req *pb.UpdateScoreRequest) (*pb.Empty, error) {
	if stats, exists := s.players[req.Name]; exists {
		stats.Average = ((stats.Average * float64(stats.Centuries)) + float64(req.Runs)) / float64(stats.Centuries+1)
		stats.Centuries++
		log.Printf("Updated %s: Avg %.2f, Centuries %d", stats.Name, stats.Average, stats.Centuries)
	} else {
		s.players[req.Name] = &pb.StatsResponse{Name: req.Name, Average: float64(req.Runs), Centuries: 1}
		log.Printf("Added new player %s with avg %.2f and 1 century", req.Name, float64(req.Runs))
	}
	return &pb.Empty{}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":6969")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	cricketServer := NewCricketServer()

	pb.RegisterCricketServer(grpcServer, cricketServer)

	fmt.Println("gRPC server running on port 6969...")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}