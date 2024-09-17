package main

import (
	"context"
	"fmt"
	service "grpc-gateway-demo/api/gate"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	"grpc-gateway-demo/api/gate/middleware"
	bookconfig "grpc-gateway-demo/bin/book"
	config "grpc-gateway-demo/bin/gate"
	handler "grpc-gateway-demo/internal/gate"
	bookpb "grpc-gateway-demo/pkg/proto/book"
	pb "grpc-gateway-demo/pkg/proto/gate"
)

func run() error {
	log := grpclog.NewLoggerV2(os.Stdout, io.Discard, io.Discard)
	grpclog.SetLoggerV2(log)

	ctx := context.Background()

	op := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(middleware.AuthInterceptor),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// 创建grpc连接 127.0.0.1:8001
	conn, err := grpc.NewClient(config.GetRpcAddr(), op...)
	if err != nil {
		log.Error("dial failed: %v", err)
		return err
	}

	sp := []runtime.ServeMuxOption{
		runtime.WithForwardResponseOption(middleware.Forward),
		runtime.WithRoutingErrorHandler(middleware.RoutingErrorHandler),
		runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
			if s == "X-User-Id" {
				return s, true
			}
			return runtime.DefaultHeaderMatcher(s)
		}),
	}
	mux := runtime.NewServeMux(sp...)

	if err = mux.HandlePath(http.MethodPost, "/book/objects", handler.Upload); err != nil {
		return err
	}

	if err = mux.HandlePath(http.MethodGet, "/book/objects/{name}", handler.Download); err != nil {
		return err
	}

	err = newGateway(ctx, conn, mux)
	if err != nil {
		log.Error("register handler failed: %v", err)
		return err
	}
	go newGrpcGate(log)
	err = newGrpcHttpServer(mux, log)
	if err != nil {
		return err
	}
	return nil
}

func newGrpcHttpServer(mux *runtime.ServeMux, log grpclog.LoggerV2) error {
	server := http.Server{
		Addr:    config.GetHttpAddr(), // 127.0.0.1:8002
		Handler: mux,
	}

	log.Infof("Serving Http on %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

// newGrpcGate 启动rpc服务
func newGrpcGate(log grpclog.LoggerV2) error {
	conn, err := grpc.NewClient(bookconfig.GetRpcAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024))) // 设置最大接收消息大小为 10MB)
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}

	client := bookpb.NewBookServiceClient(conn)

	grpcAddr := config.GetRpcAddr()
	// 127.0.0.1:12345
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		grpclog.Errorf("tcp listen failed: %v", err)
		return err
	}
	defer func() {
		if err = l.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	s := grpc.NewServer()
	// 	注册服务
	gateService := service.NewGateService(client)
	pb.RegisterGateServer(s, gateService)

	log.Infof("Serving gRPC on %s", l.Addr())

	return s.Serve(l)
}

func newGateway(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) error {
	err := bookpb.RegisterBookServiceHandler(ctx, mux, conn)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// grpc 启动
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
