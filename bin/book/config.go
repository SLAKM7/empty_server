package book

import "fmt"

const (
	Host     = "127.0.0.1"
	GrpcPort = "8001"
	HttpPort = "8002"
)

func GetRpcAddr() string {
	return fmt.Sprintf("%s:%s", Host, GrpcPort)
}

func GetHttpAddr() string {
	return fmt.Sprintf("%s:%s", Host, HttpPort)
}
