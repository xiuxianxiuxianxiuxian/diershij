package main

import (
    "log"
    "net"

    "github.com/cultivation-world/heavenly-dao/internal/service"
    "github.com/cultivation-world/shared/config"
    "google.golang.org/grpc"
)

func main() {
    cfg := config.LoadConfigFromEnv()

    heavenlyDaoSvc := service.NewHeavenlyDaoService()

    grpcServer := grpc.NewServer()
    game.RegisterHeavenlyDaoServiceServer(grpcServer, heavenlyDaoSvc)

    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    log.Printf("Heavenly Dao Service starting on :50052")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
