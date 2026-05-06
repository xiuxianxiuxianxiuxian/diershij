package service

import (
	"context"
	"fmt"
	"time"

	cultivation "github.com/cultivation-world/shared/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WorldEventClient struct {
	conn        *grpc.ClientConn
	worldClient cultivation.WorldServiceClient
}

func NewWorldEventClient(host string, port int) (*WorldEventClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &WorldEventClient{
		conn:        conn,
		worldClient: cultivation.NewWorldServiceClient(conn),
	}, nil
}

func (c *WorldEventClient) Close() error {
	return c.conn.Close()
}

func (c *WorldEventClient) GetWorldState(ctx context.Context) (*cultivation.WorldState, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.worldClient.GetWorldState(ctx, &cultivation.WorldStateRequest{})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	return resp.State, nil
}
