package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/HaiHart/ShipdockServer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"
	// grpc "google.golang.org/grpc"
)

type Container struct {
	Name   string
	Placed int32
	Iden   string
	Key    int32
	inTime time.Time
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
	context     context.Context
	cancel      context.CancelFunc
	cache       []Container
	clients     map[string]RelayConn
	toSend      map[string]chan *pb.Pack
	currCommand []Container
	port        int32
}

func (s *SerConn) MoveContainer(msg pb.Com_MoveContainerServer) error {
	ctx := msg.Context()
	p, _ := peer.FromContext(ctx)
	peerID := p.Addr.String()

	if _, ok := s.toSend[peerID]; !ok {
		s.toSend[peerID] = make(chan *pb.Pack, 1000)
	}

	go func() {
		for {
			select {
			case toSend := <-s.toSend[peerID]:
				if err := msg.Send(toSend); err != nil {
					fmt.Println(err)
					continue
				}
			}
		}

	}()
	for {
		in, err := msg.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		var changes = Container{
			Name:   in.List.Name,
			Placed: (in.List.Place),
			Key:    0,
			Iden:   in.List.Id,
			inTime: in.List.Time.AsTime(),
		}
		s.Valid(&changes)
	}

	return nil
}

func (s *SerConn) Valid(changes *Container) {

	var new_move = &pb.Pack{
		List: &pb.Container{
			Name:  changes.Name,
			Id:    changes.Iden,
			Place: changes.Placed,
			Time:  timestamppb.New(changes.inTime),
		},
	}

	s.currCommand = append([]Container{*changes}, s.currCommand...)
	if s.CheckOnCache(changes) {
		for _, i := range s.toSend {
			i <- new_move
		}
	}

	return
}

func (s *SerConn) CheckOnCache(changes *Container) bool {
	return true
}

func (s *SerConn) FetchShip(ctx context.Context, msg *pb.ShipAccess) (*pb.ShipResponse, error) {
	return nil, nil
}
