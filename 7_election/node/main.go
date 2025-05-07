package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"slices"

	pb "github.com/Sp92535/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NodeServer struct {
	pb.UnimplementedNodeServer
	co    string
	nodes map[int64]string
	pid   int64
	add   string
}

func (n *NodeServer) Election(ctx context.Context, req *pb.ElectionRequest) (*pb.Empty, error) {
	if req.Ring {
		n.ring(req.Pid)
	}
	return &pb.Empty{}, nil
}

func (n *NodeServer) Ping(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
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
		pid: int64(os.Getpid()),
		add: listner.Addr().String(),
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
	for {
		select {
		case <-ticker.C:
			nodeConn, _ := grpc.NewClient(nodeServer.co, opts)
			nodeClient := pb.NewNodeClient(nodeConn)

			_, err := nodeClient.Ping(context.Background(), &pb.Empty{})
			nodeConn.Close()

			if err != nil {
				var arr []int64
				nodeServer.ring(arr)
			}

		default:
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

			_, err := nodeClient.Election(context.Background(), &pb.ElectionRequest{Ring: false})
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

func (n *NodeServer) ring(arr []int64) {

	log.Printf("Starting Ring Election....")
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	if len(arr) > 0 && arr[0] == n.pid {

		max := slices.Max(arr)
		for _, add := range n.nodes {

			nodeConn, _ := grpc.NewClient(add, opts)
			nodeClient := pb.NewNodeClient(nodeConn)

			nodeClient.Cordinator(context.Background(), &pb.CordinatorRequest{Address: n.nodes[max], Pid: max})

			nodeConn.Close()
		}

		return
	}

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

	arr = append(arr, n.pid)
	for _, v := range sorted {
		nodeConn, _ := grpc.NewClient(n.nodes[v], opts)
		nodeClient := pb.NewNodeClient(nodeConn)

		_, err := nodeClient.Election(context.Background(), &pb.ElectionRequest{
			Ring: true,
			Pid:  arr,
		})
		if err == nil {
			break
		}
		nodeConn.Close()
	}

}
