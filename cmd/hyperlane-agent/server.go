package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/<your-username>/hyperlane/internal/agent/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type vmServer struct {
	pb.UnimplementedVMServiceServer
}

func (s *vmServer) CreateVM(ctx context.Context, req *pb.CreateVMRequest) (*pb.VMResponse, error) {
	fmt.Printf("Agent: CreateVM called: %+v\n", req)
	// TODO: integrate with KVM lifecycle
	return &pb.VMResponse{Success: true, Message: "VM creation stubbed"}, nil
}

func (s *vmServer) StopVM(ctx context.Context, req *pb.StopVMRequest) (*pb.VMResponse, error) {
	fmt.Printf("Agent: StopVM called: %+v\n", req)
	return &pb.VMResponse{Success: true, Message: "VM stop stubbed"}, nil
}

func (s *vmServer) DestroyVM(ctx context.Context, req *pb.DestroyVMRequest) (*pb.VMResponse, error) {
	fmt.Printf("Agent: DestroyVM called: %+v\n", req)
	return &pb.VMResponse{Success: true, Message: "VM destroy stubbed"}, nil
}

func (s *vmServer) ReportState(ctx context.Context, req *pb.StateReportRequest) (*pb.VMResponse, error) {
	fmt.Printf("Agent: ReportState called: %+v\n", req.Vms)
	return &pb.VMResponse{Success: true, Message: "State report stubbed"}, nil
}

func main() {
	certFile := "certs/server.crt"
	keyFile := "certs/server.key"

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterVMServiceServer(grpcServer, &vmServer{})

	fmt.Println("Hyperlane agent gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
