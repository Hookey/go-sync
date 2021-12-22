package client

import (
	"context"

	pb "github.com/Hookey/go-sync/api/pb"
	"google.golang.org/grpc"
)

// Client provides the client api.
type Client struct {
	api  pb.APIClient
	conn *grpc.ClientConn
}

// NewClient starts the client.
func NewClient(target string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		api:  pb.NewAPIClient(conn),
		conn: conn,
	}, nil
}

// Close closes the client's grpc connection and cancels any active requests.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Put uploads a single file from src to dst
func (c *Client) Put(src, dst string) (err error) {
	_, err = c.api.Put(context.Background(), &pb.PutRequest{Src: src, Dst: dst})
	return
}

// Get downloads a single file from src to dst
func (c *Client) Get(src, dst string) (err error) {
	_, err = c.api.Get(context.Background(), &pb.GetRequest{Src: src, Dst: dst})
	return
}
