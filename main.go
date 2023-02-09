package server

import (
	"context"
	"time"

	pb "github.com/HaiHart/ShipdockServer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
	// grpc "google.golang.org/grpc"
)

type Container struct {
	Name   string
	Placed int
	Iden   string
	Key    int
}

type CacheField struct {
}

type RelayConn struct {
	conn      *grpc.ClientConn
	client    pb.ComClient
	reqStream pb.Com_MoveContainerClient
}

type ResponseMessage struct {
	msg      interface{}
	clientID string
}

type SerConn struct {
	pb.UnimplementedComServer
	context context.Context
	cancel  context.CancelFunc
	cache   []Container
	clients map[string]RelayConn

	port    int32
}

func (s *SerConn) MoveContainer(msg pb.Com_MoveContainerServer) error {
	return nil
}

func (s *SerConn) FetchShip(ctx context.Context, msg *pb.ShipAccess) (*pb.ShipResponse, error) {
	return nil, nil
}
