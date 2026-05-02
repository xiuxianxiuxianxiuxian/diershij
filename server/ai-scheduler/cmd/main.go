package main

import (
    "log"
    "net"

    "github.com/cultivation-world/ai-scheduler/internal/service"
    "github.com/cultivation-world/shared/config"
    "google.golang.org/grpc"
)

func main() {
    cfg := config.LoadConfigFromEnv()

    aiSchedulerSvc := service.NewAISchedulerService(cfg)

    grpcServer := grpc.NewServer()
    game.RegisterAISchedulerServiceServer(grpcServer, aiSchedulerSvc)

    lis, err := net.Listen("tcp", ":50053")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    log.Printf("AI Scheduler Service starting on :50053")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
