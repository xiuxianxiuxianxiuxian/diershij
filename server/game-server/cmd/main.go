package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/cultivation-world/game-server/internal/repository"
	"github.com/cultivation-world/game-server/internal/service"
	cultivation "github.com/cultivation-world/shared/proto/pb"
	"github.com/cultivation-world/shared/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

	entityRepo := repository.NewEntityRepository(db, redisClient,
		repository.NewSpiritStonesRepository(db),
		repository.NewKarmaRepository(db),
	)
	itemRepo := repository.NewPostgresItemRepository(db)
	inventoryRepo := repository.NewPostgresInventoryRepository(db)
	spellRepo := repository.NewPostgresSpellRepository(db)
	messageRepo := repository.NewPostgresMessageRepository(db)
	sectRepo := repository.NewSectRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)
	friendRepo := repository.NewFriendRepository(db)

	worldConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.WorldEngine.Host, cfg.WorldEngine.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to world engine: %v", err)
	}
	defer worldConn.Close()
	worldClient := service.NewWorldGrpcClient(worldConn)

	daoConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.DaoService.Host, cfg.DaoService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to heavenly dao: %v", err)
	}
	defer daoConn.Close()
	daoClient := service.NewHeavenlyDaoGrpcClient(daoConn)

	operationSvc := service.NewOperationService(entityRepo, itemRepo, inventoryRepo, spellRepo, messageRepo, worldClient, daoClient,
		service.NewSectRepoAdapter(sectRepo), service.NewRecipeRepoAdapter(recipeRepo), service.NewFriendRepoAdapter(friendRepo))
	gameSvc := service.NewGameService(entityRepo, operationSvc, spellRepo, itemRepo, inventoryRepo)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("PANIC recovered in %s: %v", info.FullMethod, r)
					err = status.Errorf(codes.Internal, "internal server error")
				}
			}()
			return handler(ctx, req)
		}),
	)
	cultivation.RegisterGameServiceServer(grpcServer, gameSvc)

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
