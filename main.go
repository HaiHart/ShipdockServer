package server

import (
	"context"

	pb "github.com/HaiHart/ShipdockServer/proto"
	"google.golang.org/grpc"
)

type ClientConn struct{
	
}

type ResponseMessage struct{
	msg interface{}
	clientID string
}

type SerConn struct {
	pb.UnimplementedComServer
	context               context.Context
	cancel                context.CancelFunc
	port int32
	
}