package gate

import (
	"context"
	"google.golang.org/protobuf/proto"
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
	var resData []byte
	switch req.Method {
	case bookpb.BookService_CreateBook_FullMethodName:
		gateReq := &bookpb.CreateBookRequest{}
		err := proto.Unmarshal(req.Data, gateReq)
		if err != nil {
			return nil, err
		}
		res, err := svc.clients[0].CreateBook(ctx, gateReq)
		if err != nil {
			return nil, err
		}
		resData, err = proto.Marshal(res)
		if err != nil {
			return nil, err
		}
	case bookpb.BookService_GetBook_FullMethodName:
		gateReq := &bookpb.GetBookRequest{}
		err := proto.Unmarshal(req.Data, gateReq)
		if err != nil {
			return nil, err
		}
		res, err := svc.clients[0].GetBook(ctx, gateReq)
		if err != nil {
			return nil, err
		}
		resData, err = proto.Marshal(res)
		if err != nil {
			return nil, err
		}
	}

	return &pb.RpcResponse{
		Data: resData,
	}, nil
}
