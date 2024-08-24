package main

import (
	"context"
	"log"
	"time"

	pb "pro/pb"

	"google.golang.org/grpc"
)

const (
	address = "10.129.173.253:30030"
	path    = "/test2"
)

func main() {
	// 建立连接
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()
	c := pb.NewZkServiceClient(conn)

	// 创建一个用于取消操作的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 启动写入和读取的 goroutine
	go writeLoop(ctx, c)
	go readLoop(ctx, c)

	// 让主程序运行一段时间（比如60秒）
	time.Sleep(300 * time.Second)
}

func writeLoop(ctx context.Context, c pb.ZkServiceClient) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentTime := time.Now().Format(time.RFC3339)
			_, err := c.Set(ctx, &pb.PathAndData{Path: path, Data: currentTime})
			if err != nil {
				log.Printf("写入错误: %v", err)
			} else {
				log.Printf("写入成功: %s", currentTime)
			}
		}
	}
}

func readLoop(ctx context.Context, c pb.ZkServiceClient) {
	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r, err := c.Get(ctx, &pb.Path{Path: path})
			if err != nil {
				log.Printf("读取错误: %v", err)
			} else {
				log.Printf("读取结果: %s", r.Data)
			}
		}
	}
}
