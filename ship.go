package main

import (
	"context"
	pb "github.com/HaiHart/ShipdockServer/proto"
	"sync"
	"time"
)

type Ship struct {
	Name    string
	Placed  int32
	Iden    string
	Key     int32
	Detail  detail
	inTime  time.Time
	outTime time.Time
	length  uint32
}

type Doc struct {
	No           int32
	Name         string
	length       uint32
	boarderLeft  int32
	boarderRight int32
	occ          uint32
	ShipList     map[string][]string
}

type ShipConn struct {
	pb.UnimplementedComServer
	context     context.Context
	cancel      context.CancelFunc
	cache       []Container
	clients     map[string]RelayConn
	toSend      map[string]chan *pb.Pack
	lock        sync.Mutex
	currCommand []Container
	port        int32
	log         []string
	detailLog   []string
}
