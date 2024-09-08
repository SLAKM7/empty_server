package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	protoDir := "./proto"
	outputDir := "./pkg/proto"

	// 检查并创建输出目录
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// 遍历 proto 目录
	err := filepath.Walk(protoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .proto 文件
		if filepath.Ext(path) == ".proto" {
			// 构建 protoc 命令
			cmd := exec.Command("protoc",
				"-I", protoDir,
				"--grpc-gateway_out="+outputDir,
				"--grpc-gateway_opt", "logtostderr=true",
				"--grpc-gateway_opt", "paths=source_relative",
				"--go_out="+outputDir,
				"--go_opt", "paths=source_relative",
				"--go-grpc_out="+outputDir,
				"--go-grpc_opt", "paths=source_relative",
				path,
			)

			// 设置命令输出到控制台
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// 执行命令
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to generate %s: %v", path, err)
			}
			fmt.Printf("Generated code for: %s\n", path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
	}
}
