package main

import (
    "fmt"
    "log"
    "net"
    "os"

    "github.com/cultivation-world/ai-scheduler/internal/service"
    "github.com/cultivation-world/shared/proto/pb"
    "google.golang.org/grpc"
)

func main() {
    aiSvc := service.NewAISchedulerService()

    grpcServer := grpc.NewServer()
    pb.RegisterAISchedulerServiceServer(grpcServer, aiSvc)

    port := 50054
    if p := os.Getenv("GRPC_PORT"); p != "" {
        fmt.Sscanf(p, "%d", &port)
    }

    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    log.Printf("AI Scheduler Service starting on :%d", port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
