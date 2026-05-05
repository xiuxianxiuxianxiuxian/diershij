package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/cultivation-world/world-engine/internal/service"
	cultivation "github.com/cultivation-world/shared/proto/pb"
	"google.golang.org/grpc"
)

func main() {
	worldSvc := service.NewWorldEngineService()

	// 后台资源刷新协程：每5分钟刷新一次资源
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("World Engine: 开始刷新资源...")
			worldSvc.AdvanceEpoch()
		}
	}()

	grpcServer := grpc.NewServer()
	cultivation.RegisterWorldServiceServer(grpcServer, worldSvc)

	port := 50054
	if p := os.Getenv("GRPC_PORT"); p != "" {
		fmt.Sscanf(p, "%d", &port)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("World Engine Service starting on :%d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
