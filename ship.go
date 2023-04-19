package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/HaiHart/ShipdockServer/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Ship struct {
	Name    string
	Placed  int
	Iden    string
	Key     int32
	Detail  detail
	InTime  time.Time
	OutTime time.Time
	Length  uint32
}

type Doc struct {
	No           int
	Name         string
	Length       uint32
	BoarderRight int
	ShipList     []string
}

type ShipConn struct {
	pb.UnimplementedComServer
	context     context.Context
	cancel      context.CancelFunc
	ships       []Ship
	docks       []Doc
	clients     map[string]RelayConn
	toSend      map[string]chan *pb.Pack
	lock        sync.Mutex
	currCommand []Container
	port        int32
	log         []string
	detailLog   []string
}

func (s *ShipConn) FetchDocks(ctx context.Context, time *pb.Header) (*pb.Docks, error) {
	fmt.Println(time.Time)

	var list []*pb.Ship

	var list_2[]*pb.Doc

	var log []string = make([]string, 0)

	for _, v := range s.ships {
		list = append(list, &pb.Ship{
			Name: v.Name,
			Placed: int32(v.Placed),
			Iden: v.Iden,
			Key: v.Key,
			Detail: &pb.Detail{
				From: v.Detail.From,
				By: v.Detail.by,
				AtTime: v.Detail.atTime,
				Owner: v.Detail.owner,
			},
			InTime: timestamppb.New(v.InTime),
			OutTime: timestamppb.New(v.OutTime),
		})
	}
	for _, v := range s.log {
		log = append(log, v)
	}

	for _,v:=range s.docks{
		list_2=append(list_2, &pb.Doc{
			No: int32(v.No),
			Name: v.Name,
			Length: v.Length,
			BoarderRight: int32(v.BoarderRight),
			ShipList: v.ShipList,
		})
	}

	return &pb.Docks{
		L
	}
}