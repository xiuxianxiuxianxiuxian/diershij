package main

import (
    "fmt"
    "log"
    "net"
    "os"

    "github.com/cultivation-world/game-server/internal/repository"
    "github.com/cultivation-world/game-server/internal/service"
    "github.com/cultivation-world/shared/proto/pb"
    "github.com/cultivation-world/shared/config"
    "google.golang.org/grpc"
)

func main() {
    cfg := config.LoadConfigFromEnv()

    db, err := repository.NewDatabase(&cfg.Database)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    redisClient := repository.NewRedisClient(&cfg.Redis)
    defer redisClient.Close()

    entityRepo := repository.NewEntityRepository(db, redisClient)
    operationSvc := service.NewOperationService(entityRepo)
    gameSvc := service.NewGameService(entityRepo, operationSvc)

    grpcServer := grpc.NewServer()
    pb.RegisterGameServiceServer(grpcServer, gameSvc)

    port := 50051
    if p := os.Getenv("GRPC_PORT"); p != "" {
        fmt.Sscanf(p, "%d", &port)
    }

    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    log.Printf("Game Server starting on :%d", port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
