package geo

import (
	"context"
	"time"

	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/generated/clients/geosrv/geopb"
	"delivery/internal/pkg/errs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ ports.GeoClient = &Client{}

type Client struct {
	conn    *grpc.ClientConn
	client  geopb.GeoClient
	timeout time.Duration
}

func NewClient(host string) (*Client, error) {
	if host == "" {
		return nil, errs.NewValueIsRequiredError("host")
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		client:  geopb.NewGeoClient(conn),
		timeout: 5 * time.Second,
	}, nil
}

func (c *Client) GetGeolocation(ctx context.Context, street string) (kernel.Location, error) {
	req := &geopb.GetGeolocationRequest{Street: street}
	res, err := c.client.GetGeolocation(ctx, req)
	if err != nil {
		return kernel.Location{}, err
	}

	return kernel.NewLocation(int(res.Location.X), int(res.Location.Y))
}

func (c *Client) Close() error {
	return c.conn.Close()
}
