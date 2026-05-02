package main

import (
    "log"
    "net"

    "github.com/cultivation-world/shared/config"
    "github.com/cultivation-world/world-engine/internal/service"
    "google.golang.org/grpc"
)

func main() {
    cfg := config.LoadConfigFromEnv()

    worldEngineSvc := service.NewWorldEngineService()

    grpcServer := grpc.NewServer()
    game.RegisterWorldServiceServer(grpcServer, worldEngineSvc)

    lis, err := net.Listen("tcp", ":50054")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    log.Printf("World Engine Service starting on :50054")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
