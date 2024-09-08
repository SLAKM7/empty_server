package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var notFoundMsg = "路径不存在"

func AuthInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	return invoker(ctx, method, req, reply, cc, opts...)
}
func Forward(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
	writer.Header().Set("test-aa", "bbbb")
	writer.WriteHeader(200)
	return nil
}

func RoutingErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, i int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNotFound)
	msg := fmt.Sprintf("%s %s", request.RequestURI, notFoundMsg)
	writer.Write([]byte(msg))
	return
}
