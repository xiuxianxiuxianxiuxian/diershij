package service

import (
    "context"
    "fmt"
    "time"

    "github.com/cultivation-world/shared/proto/pb"
    "github.com/cultivation-world/shared/types"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type GameServiceClient struct {
    conn       *grpc.ClientConn
    gameClient pb.GameServiceClient
}

func NewGameServiceClient(host string, port int) (*GameServiceClient, error) {
    addr := fmt.Sprintf("%s:%d", host, port)
    conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, err
    }

    return &GameServiceClient{
        conn:       conn,
        gameClient: pb.NewGameServiceClient(conn),
    }, nil
}

func (c *GameServiceClient) Close() error {
    return c.conn.Close()
}

func (c *GameServiceClient) ExecuteOperation(op *types.Operation) (*types.OperationResult, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    params := make(map[string]string)
    for k, v := range op.Params {
        params[k] = fmt.Sprintf("%v", v)
    }

    resp, err := c.gameClient.ExecuteOperation(ctx, &pb.OperationRequest{
        OperationId: op.ID,
        ActorId:     string(op.ActorID),
        ActionType:  string(op.ActionType),
        Params:      params,
        Timestamp:   op.Timestamp,
    })

    if err != nil {
        return nil, err
    }

    effects := make(map[string]interface{})
    for k, v := range resp.Effects {
        effects[k] = v
    }

    return &types.OperationResult{
        Success:   resp.Success,
        Message:   resp.Message,
        Effects:   effects,
        Timestamp: resp.Timestamp,
    }, nil
}

func (c *GameServiceClient) CreateEntity(ctx context.Context, username, password string, entityType types.EntityType) (*types.Entity, error) {
    resp, err := c.gameClient.CreateEntity(ctx, &pb.CreateEntityRequest{
        Name:       username,
        EntityType: string(entityType),
    })

    if err != nil {
        return nil, err
    }

    return protoToEntity(resp.Entity), nil
}

func (c *GameServiceClient) AuthenticateEntity(ctx context.Context, username, password string) (*types.Entity, error) {
    resp, err := c.gameClient.AuthenticateEntity(ctx, &pb.AuthRequest{
        Username: username,
        Password: password,
    })

    if err != nil {
        return nil, err
    }

    return protoToEntity(resp.Entity), nil
}

func (c *GameServiceClient) GetEntity(ctx context.Context, entityID types.EntityID) (*types.Entity, error) {
    resp, err := c.gameClient.GetEntity(ctx, &pb.EntityRequest{
        EntityId: string(entityID),
    })

    if err != nil {
        return nil, err
    }

    return protoToEntity(resp.Entity), nil
}

func protoToEntity(e *pb.Entity) *types.Entity {
    if e == nil {
        return nil
    }

    return &types.Entity{
        ID:         types.EntityID(e.Id),
        EntityType: types.EntityType(e.EntityType),
        Name:       e.Name,
        Realm:      types.CultivationRealm(e.Realm),
        Position: types.WorldPosition{
            RegionID: e.Position.RegionId,
            X:        e.Position.X,
            Y:        e.Position.Y,
        },
        Attributes: types.Attributes{
            Qi:                  e.Attributes.Qi,
            MaxQi:               e.Attributes.MaxQi,
            SpiritualPower:      e.Attributes.SpiritualPower,
            MaxSpiritualPower:   e.Attributes.MaxSpiritualPower,
            DivineSense:         e.Attributes.DivineSense,
            Comprehension:       int(e.Attributes.Comprehension),
            Constitution:        int(e.Attributes.Constitution),
            Luck:                int(e.Attributes.Luck),
            CultivationProgress: e.Attributes.CultivationProgress,
            AttackPower:         e.Attributes.AttackPower,
            Defense:             e.Attributes.Defense,
            Speed:               e.Attributes.Speed,
            MentalStability:     int(e.Attributes.MentalStability),
            RemainingLifespan:   int(e.Attributes.RemainingLifespan),
            MaxLifespan:         int(e.Attributes.MaxLifespan),
        },
        Karma: types.Karma{
            KarmaValue:   int(e.Karma.KarmaValue),
            Merit:        int(e.Karma.Merit),
            HeavenlyMark: e.Karma.HeavenlyMark,
        },
        Status:    types.EntityStatus(e.Status),
        CreatedAt: time.Unix(e.CreatedAt, 0),
        UpdatedAt: time.Unix(e.UpdatedAt, 0),
    }
}
