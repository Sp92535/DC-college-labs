package main

import (
	"context"
	"log"
	"net"
	"os"
	"slices"
	"time"

	pb "github.com/Sp92535/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NodeServer struct {
	pb.UnimplementedNodeServer
	co       string
	nodes    map[int64]string
	pid      int64
	add      string
	hasToken bool
}

func (n *NodeServer) Election(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (n *NodeServer) GiveToken(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	n.hasToken = true
	return &pb.Empty{}, nil
}

func (n *NodeServer) Ping(ctx context.Context, req *pb.Empty) (*pb.EnterResponse, error) {

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	criticalConn, err := grpc.NewClient(":9595", opts)
	if err != nil {
		log.Fatalf("error connecting critical :%v", err)
	}

	criticalClient := pb.NewCriticalClient(criticalConn)
	res, err := criticalClient.Enter(context.Background(), &pb.Empty{})
	criticalConn.Close()
	return res, err
}

func (n *NodeServer) Cordinator(ctx context.Context, req *pb.CordinatorRequest) (*pb.Empty, error) {
	n.co = req.Address
	log.Printf("New Cordinator %s Pid %d", req.Address, req.Pid)
	return &pb.Empty{}, nil
}

func main() {

	listner, err := net.Listen("tcp", "localhost:0")

	if err != nil {
		log.Fatalf("listner error %v", err)
	}

	grpcServer := grpc.NewServer()

	nodeServer := &NodeServer{
		pid:      int64(os.Getpid()),
		add:      listner.Addr().String(),
		hasToken: os.Args[1] == "token",
	}

	pb.RegisterNodeServer(grpcServer, nodeServer)

	log.Printf("grpc node server running on %s...", nodeServer.add)
	go func() {
		if err = grpcServer.Serve(listner); err != nil {
			log.Fatalf("Failed to Serve: %v", err)
		}
	}()

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	trackerConn, err := grpc.NewClient(":6969", opts)
	if err != nil {
		log.Fatalf("error connecting tracker :%v", err)
	}

	trackerClient := pb.NewTrackerClient(trackerConn)
	trackerClient.Register(context.Background(), &pb.RegisterRequest{
		Address: listner.Addr().String(),
		Pid:     nodeServer.pid,
	})

	// periodically get nodes
	go func() {
		ticker := time.NewTicker(10 * time.Second)

		res, _ := trackerClient.GetNodes(context.Background(), &pb.Empty{})
		nodeServer.nodes = res.Address

		for {
			select {
			case <-ticker.C:
				res, _ := trackerClient.GetNodes(context.Background(), &pb.Empty{})
				nodeServer.nodes = res.Address
			default:
			}
		}

	}()

	ticker := time.NewTicker(10 * time.Second)
	for os.Args[1] == "cent" {
		select {
		case <-ticker.C:
			nodeConn, _ := grpc.NewClient(nodeServer.co, opts)
			nodeClient := pb.NewNodeClient(nodeConn)

			res, err := nodeClient.Ping(context.Background(), &pb.Empty{})
			log.Print(res)
			nodeConn.Close()

			if err != nil {
				log.Println(err)
				// var arr []int64
				// nodeServer.ring(arr)
				nodeServer.bully()
			}

		default:
		}
	}

	for {
		if nodeServer.hasToken {
			select {
			case <-ticker.C:
				nodeServer.ring()
			default:
			}
		}
	}

}

func (n *NodeServer) bully() {
	log.Printf("Starting Bully Election....")
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	var cor bool = true
	for pid, add := range n.nodes {
		if pid > n.pid {
			nodeConn, _ := grpc.NewClient(add, opts)
			nodeClient := pb.NewNodeClient(nodeConn)

			_, err := nodeClient.Election(context.Background(), &pb.Empty{})
			if err == nil {
				cor = false
			}
			nodeConn.Close()
		}
	}

	if cor {
		for _, add := range n.nodes {
			nodeConn, _ := grpc.NewClient(add, opts)
			nodeClient := pb.NewNodeClient(nodeConn)

			nodeClient.Cordinator(context.Background(), &pb.CordinatorRequest{Address: n.add, Pid: n.pid})

			nodeConn.Close()
		}
	}
}

func (n *NodeServer) ring() {

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	criticalConn, err := grpc.NewClient(":9595", opts)
	if err != nil {
		log.Fatalf("error connecting critical :%v", err)
	}

	criticalClient := pb.NewCriticalClient(criticalConn)

	res, _ := criticalClient.Enter(context.Background(), &pb.Empty{})
	criticalConn.Close()
	log.Println(res)

	var sorted []int64

	for pid := range n.nodes {
		sorted = append(sorted, pid)
	}
	slices.Sort(sorted)

	for id, val := range sorted {
		if val > n.pid {
			sorted = slices.Concat(sorted[id:], sorted[:id])
			break
		}
	}
	for _, v := range sorted {
		nodeConn, _ := grpc.NewClient(n.nodes[v], opts)
		nodeClient := pb.NewNodeClient(nodeConn)

		_, err := nodeClient.GiveToken(context.Background(), &pb.Empty{})
		if err == nil {
			n.hasToken = false
			break
		}
		nodeConn.Close()
	}

}
