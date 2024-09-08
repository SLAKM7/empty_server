package gate

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	bookpb "grpc-gateway-demo/pkg/proto/book"
	pb "grpc-gateway-demo/pkg/proto/gate"
)

type Service struct {
	pb.UnimplementedGateServer
	clients []bookpb.BookServiceClient
}

func NewGateService(clients ...bookpb.BookServiceClient) *Service {
	return &Service{
		UnimplementedGateServer: pb.UnimplementedGateServer{},
		clients:                 clients,
	}
}

func (svc *Service) Rpc(ctx context.Context, req *pb.RpcRequest) (*pb.RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rpc not implemented")
}
