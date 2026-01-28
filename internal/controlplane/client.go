package controlplane

import (
	"context"
	"log"
	"time"

	pb "github.com/<your-username>/hyperlane/internal/agent/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type AgentClient struct {
	client pb.VMServiceClient
	conn   *grpc.ClientConn
}

func NewAgentClient(addr, certFile string) (*AgentClient, error) {
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	return &AgentClient{
		client: pb.NewVMServiceClient(conn),
		conn:   conn,
	}, nil
}

func (a *AgentClient) CreateVM(name string, cpu, ram int, image string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := a.client.CreateVM(ctx, &pb.CreateVMRequest{
		Name:  name,
		CPU:   int32(cpu),
		RAM:   int32(ram),
		Image: image,
	})
	if err != nil {
		log.Println("CreateVM error:", err)
		return
	}
	log.Println("CreateVM response:", resp.Message)
}

func (a *AgentClient) Close() {
	a.conn.Close()
}
