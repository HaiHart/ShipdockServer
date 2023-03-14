package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	pb "github.com/HaiHart/ShipdockServer/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	lock        sync.Mutex
	currCommand []Container
	port        int32
	log         []string
	detailLog   []string
}

func (s *SerConn) FetchList(ctx context.Context, time *pb.Header) (*pb.ShipList, error) {
	fmt.Println(time.Time)

	var list []*pb.ContainerSet

	var log []string = make([]string, 0)

	for _, v := range s.cache {
		list = append(list, &pb.ContainerSet{
			Name:  v.Name,
			Id:    v.Iden,
			Key:   v.Key,
			Place: v.Placed,
		})
	}
	for _, v := range s.log {
		log = append(log, v)
	}

	return &pb.ShipList{
		List: list,
		Log:  log,
	}, nil
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
		fmt.Printf("got command at %v\n", time.Now().UTC())
		var changes = Container{
			Name:   in.List[0].Name,
			Placed: (in.List[0].Place),
			Key:    0,
			Iden:   in.List[0].Id,
			inTime: in.List[0].Time.AsTime(),
		}
		var new_place = in.List[0].NewPlace
		var swap = in.Swap
		if swap {
			var changes_2 = Container{
				Name:   in.List[1].Name,
				Placed: (in.List[1].Place),
				Key:    0,
				Iden:   in.List[1].Id,
				inTime: in.List[1].Time.AsTime(),
			}
			s.ValidSwap(&changes, &changes_2, peerID)
		}

		s.ValidMove(&changes, new_place, peerID)
	}

	return nil
}

func (s *SerConn) ValidMove(changes *Container, new_place int32, peerID string) {

	if changes == nil {
		return
	}

	var new_move = &pb.Pack{
		List: []*pb.Container{
			{
				Name:     changes.Name,
				Id:       changes.Iden,
				Place:    changes.Placed,
				Time:     timestamppb.New(changes.inTime),
				NewPlace: new_place,
			},
		},
		Swap: false,
		Err:  "None",
	}

	// fmt.Println(new_place)
	if s.CheckOnCacheMove(changes, new_place) {
		s.currCommand = append([]Container{*changes}, s.currCommand...)
		var detail = fmt.Sprintf("%v:%v:%v is moved to %v ay %v", peerID,len(s.log), changes.Iden, new_place, time.Now().Format(time.ANSIC))
		s.detailLog = append(s.detailLog, detail)
		s.Swap(changes.Iden, int(new_place))
		for _, i := range s.toSend {
			i <- new_move
		}

	} else {
		new_move.Err = "Miss match starting point"
		s.toSend[peerID] <- new_move
	}

	return
}

func (s *SerConn) ValidSwap(changes *Container, changes_2 *Container, peerOD string) {
	if changes == nil || changes_2 == nil {
		return
	}

	var new_move = &pb.Pack{
		List: []*pb.Container{
			{
				Name:     changes.Name,
				Id:       changes.Iden,
				Place:    changes.Placed,
				Time:     timestamppb.New(changes.inTime),
				NewPlace: changes_2.Placed,
			},
			{
				Name:     changes_2.Name,
				Id:       changes_2.Iden,
				Place:    changes_2.Placed,
				Time:     timestamppb.New(changes_2.inTime),
				NewPlace: changes.Placed,
			},
		},
		Swap: true,
		Err:  "None",
	}
	if s.CheckOnCacheSwap(changes, changes_2) {
		s.currCommand = append([]Container{*changes, *changes_2}, s.currCommand...)
		for _, i := range s.toSend {
			i <- new_move
		}
		var detail = fmt.Sprintf("%v:%v:%v is switched with %v at %v", peerOD,len(s.log), changes.Iden, changes_2.Iden, time.Now().Format(time.ANSIC))
		s.detailLog = append(s.detailLog, detail)
		s.Swap(changes.Iden, int(changes_2.Placed))
	} else {
		new_move.Err = "Miss match starting point"
		s.toSend[peerOD] <- new_move
	}

	return
}

func (s *SerConn) CheckOnCacheMove(changes *Container, new_place int32) bool {

	if changes == nil {
		return false
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if changes.Placed == new_place {
		return false
	}
	for _, v := range s.cache {
		if v.Placed == new_place && new_place != -1 {
			return false
		}
	}
	for _, v := range s.cache {
		if v.Iden == changes.Iden {
			if v.Placed != changes.Placed {
				return false
			}
		}
	}
	return true
}

func (s *SerConn) CheckOnCacheSwap(changes *Container, changes_2 *Container) bool {
	if changes == nil || changes_2 == nil {
		return false
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if changes.Placed == changes_2.Placed {
		return false
	}
	for _, v := range s.cache {
		if v.Iden == changes.Iden && v.Placed != changes.Placed {
			return false
		}
		if v.Iden == changes_2.Iden && v.Placed != changes_2.Placed {
			return false
		}
	}
	fmt.Println("Here")
	for _, v := range s.cache {
		if v.Iden == changes.Iden {
			v.Placed = changes_2.Placed
		}
		if v.Iden == changes_2.Iden {
			v.Placed = changes.Placed
		}
	}
	return true
}

func (s *SerConn) FetchShip(ctx context.Context, msg *pb.ShipAccess) (*pb.ShipResponse, error) {
	return nil, nil
}

func (s *SerConn) Swap(x string, place int) {
	// s.lock.Lock()
	// defer s.lock.Unlock()
	var rv string = ""
	for k, v := range s.cache {
		if v.Iden == x {
			if place == -1 {
				(s.cache)[k] = Container{
					Iden:   v.Iden,
					Name:   v.Name,
					Placed: int32(place),
				}
				rv = string(fmt.Sprintf("%v is moved to %v at %v", x, place, time.Now().Format(time.ANSIC)))
			} else {
				for i, j := range s.cache {
					if j.Placed == int32(place) {
						if j.Iden != x {
							(s.cache)[i] = Container{
								Iden:   j.Iden,
								Name:   j.Name,
								Placed: v.Placed,
							}
							rv = string(fmt.Sprintf("%v is switched with %v at %v", x, j.Iden, time.Now().Format(time.ANSIC)))
						}
					}
				}
			}
			(s.cache)[k] = Container{
				Iden:   v.Iden,
				Name:   v.Name,
				Placed: int32(place),
			}
			if len(rv) < 1 {
				rv = string(fmt.Sprintf("%v is moved to %d at %v", x, place, time.Now().Format(time.ANSIC)))
			}

		}
	}
	fmt.Println(rv, place)
	s.log = append(s.log, rv)

}

func (s *SerConn) Start() error {
	ln, _ := net.Listen("tcp", fmt.Sprintf("localhost:%v", s.port))
	var opts []grpc.ServerOption
	var grpcServer = grpc.NewServer(opts...)
	pb.RegisterComServer(grpcServer, s)
	err := grpcServer.Serve(ln)
	if err != nil {
		return err
	}
	return nil
}

func (s *SerConn) RunServer() error {
	defer s.cancel()
	var err error = s.Start()

	return err
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := SerConn{
		port:    8080,
		context: ctx,
		cancel:  cancel,
		cache: []Container{
			{
				Name:   "1",
				Placed: -1,
				Iden:   "1",
				Key:    0,
				inTime: time.Now(),
			},
			{
				Name:   "2",
				Placed: -1,
				Iden:   "2",
				Key:    1,
				inTime: time.Now(),
			},
			{
				Name:   "3",
				Placed: -1,
				Iden:   "3",
				Key:    2,
				inTime: time.Now(),
			},
			{
				Name:   "4",
				Placed: -1,
				Iden:   "4",
				Key:    3,
				inTime: time.Now(),
			},
			{
				Name:   "5",
				Placed: -1,
				Iden:   "5",
				Key:    4,
				inTime: time.Now(),
			},
			{
				Name:   "6",
				Placed: -1,
				Iden:   "6",
				Key:    5,
				inTime: time.Now(),
			},
		},
		toSend:    make(map[string]chan *pb.Pack),
		clients:   make(map[string]RelayConn),
		log:       make([]string, 0),
		detailLog: make([]string, 0),
	}
	var group errgroup.Group
	fmt.Println("here")
	group.Go(server.RunServer)
	fmt.Println("Started")

	err := group.Wait()
	if err != nil {
		fmt.Println(err)
		return
	}
}
