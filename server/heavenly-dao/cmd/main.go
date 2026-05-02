package main

import (
    "fmt"
    "log"
    "net"
    "os"

    "github.com/cultivation-world/heavenly-dao/internal/service"
    "github.com/cultivation-world/shared/proto/pb"
    "google.golang.org/grpc"
)

func main() {
    heavenlyDaoSvc := service.NewHeavenlyDaoService()

    grpcServer := grpc.NewServer()
    pb.RegisterHeavenlyDaoServiceServer(grpcServer, heavenlyDaoSvc)

    port := 50053
    if p := os.Getenv("GRPC_PORT"); p != "" {
        fmt.Sscanf(p, "%d", &port)
    }

    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    log.Printf("Heavenly Dao Service starting on :%d", port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
