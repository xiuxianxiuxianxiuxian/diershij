package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/cultivation-world/shared/proto/pb"
	"google.golang.org/grpc"
)

type OperationResult struct {
	Success   bool
	Message   string
	Effects   map[string]interface{}
}

type GameServiceClient interface {
	ExecuteOperation(ctx context.Context, actorID string, actionType string, params map[string]string) (*OperationResult, error)
	GetEntity(ctx context.Context, entityID string) (*EntityInfo, error)
}

type EntityInfo struct {
	ID         string
	Name       string
	Realm      string
	RegionID   string
	EntityType string
}

type gameGrpcClient struct {
	client pb.GameServiceClient
}

func NewGameGrpcClient(cc grpc.ClientConnInterface) GameServiceClient {
	return &gameGrpcClient{client: pb.NewGameServiceClient(cc)}
}

func (g *gameGrpcClient) ExecuteOperation(ctx context.Context, actorID string, actionType string, params map[string]string) (*OperationResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if params == nil {
		params = make(map[string]string)
	}

	resp, err := g.client.ExecuteOperation(ctx, &pb.OperationRequest{
		OperationId: fmt.Sprintf("npc-%s-%d", actorID, time.Now().UnixNano()),
		ActorId:     actorID,
		ActionType:  actionType,
		Params:      params,
		Timestamp:   time.Now().UnixNano(),
	})
	if err != nil {
		return nil, fmt.Errorf("game-server ExecuteOperation failed: %w", err)
	}

	effects := make(map[string]interface{})
	for k, v := range resp.Effects {
		var decoded interface{}
		if err := json.Unmarshal([]byte(v), &decoded); err == nil {
			effects[k] = decoded
		} else {
			effects[k] = v
		}
	}

	return &OperationResult{
		Success: resp.Success,
		Message: resp.Message,
		Effects: effects,
	}, nil
}

func (g *gameGrpcClient) GetEntity(ctx context.Context, entityID string) (*EntityInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := g.client.GetEntity(ctx, &pb.EntityRequest{
		EntityId: entityID,
	})
	if err != nil {
		return nil, fmt.Errorf("game-server GetEntity failed: %w", err)
	}
	if !resp.Found || resp.Entity == nil {
		return nil, nil
	}

	regionID := ""
	if resp.Entity.Position != nil {
		regionID = resp.Entity.Position.RegionId
	}

	return &EntityInfo{
		ID:         resp.Entity.Id,
		Name:       resp.Entity.Name,
		Realm:      resp.Entity.Realm,
		RegionID:   regionID,
		EntityType: resp.Entity.EntityType,
	}, nil
}
