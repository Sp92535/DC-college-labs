package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/Sp92535/pb"
	"google.golang.org/grpc"
)

type CriticalServer struct {
	pb.UnimplementedCriticalServer
	
	mu sync.Mutex
	shared int64 
}


func (c *CriticalServer) Enter(ctx context.Context, req *pb.Empty) (*pb.EnterResponse, error) {
	
	c.mu.Lock()
	c.shared++
	c.mu.Unlock()
	
	return &pb.EnterResponse{
		Shared: c.shared,
	}, nil
}

func main() {

	listner, err := net.Listen("tcp", "localhost:9595")

	if err != nil {
		log.Fatalf("listner error %v", err)
	}

	grpcServer := grpc.NewServer()

	criticalServer := &CriticalServer{
		shared: 0,
	}

	pb.RegisterCriticalServer(grpcServer, criticalServer)

	log.Printf("grpc critical server running on %s...", listner.Addr().String())
	if err = grpcServer.Serve(listner); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
