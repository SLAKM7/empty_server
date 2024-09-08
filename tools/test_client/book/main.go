package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	bookpb "grpc-gateway-demo/pkg/proto/book"
	gatepb "grpc-gateway-demo/pkg/proto/gate"
)

const (
	address = "localhost:12345" // gRPC 服务器地址
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024))) // 设置最大接收消息大小为 10MB)
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := gatepb.NewGateClient(conn)

	// 设置上下文，添加超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	bookReq := &bookpb.GetBookRequest{
		Id: 2,
	}
	marshal, err := proto.Marshal(bookReq)
	if err != nil {
		return
	}
	req := &gatepb.RpcRequest{
		Service: "Book",
		Method:  "GetBook",
		Data:    marshal,
	}
	response, err := client.Rpc(ctx, req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// 输出回应
	res := &bookpb.GetBookResponse{}
	err = proto.Unmarshal(response.Data, res)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
		return
	}
	fmt.Printf("Greeting: %s\n", res.Data)
}
