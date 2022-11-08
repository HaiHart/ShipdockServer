package server

import (
	"context"

	pb "github.com/server/proto"
	"google.golang.org/grpc"
)

type ClientConn struct{
	
}

type SerConn struct {
	pb.UnimplementedComServer
	context               context.Context
	cancel                context.CancelFunc
	port int32
	
}