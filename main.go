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

type detail struct {
	From   string
	atTime string
	by     string
	owner  string
}
type Container struct {
	Name   string
	Placed int32
	Iden   string
	Key    int32
	Detail detail
	inTime time.Time
}

type RelayConn struct {
	conn      *grpc.ClientConn
	client    pb.ComClient
	reqStream pb.Com_MoveContainerClient
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
			Detail: &pb.Detail{
				From:   v.Detail.From,
				By:     v.Detail.by,
				AtTime: v.Detail.atTime,
				Owner:  v.Detail.owner,
			},
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
			Detail: detail{
				From:   in.List[0].Detail.From,
				atTime: in.List[0].Detail.AtTime,
				by:     in.List[0].Detail.By,
				owner:  in.List[0].Detail.Owner,
			},
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
				Detail: detail{
					From:   in.List[1].Detail.From,
					atTime: in.List[1].Detail.AtTime,
					by:     in.List[1].Detail.By,
					owner:  in.List[1].Detail.Owner,
				},
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
				Detail: &pb.Detail{
					From:   changes.Detail.From,
					AtTime: changes.Detail.atTime,
					By:     changes.Detail.by,
					Owner:  changes.Detail.owner,
				},
			},
		},
		Swap: false,
		Err:  "None",
	}

	// fmt.Println(new_place)
	if s.CheckOnCacheMove(changes, new_place) {
		s.currCommand = append([]Container{*changes}, s.currCommand...)
		var detail = fmt.Sprintf("%v:%v:%v is moved to %v ay %v", peerID, len(s.log), changes.Iden, new_place, time.Now().Format(time.ANSIC))
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
				Detail: &pb.Detail{
					From:   changes.Detail.From,
					AtTime: changes.Detail.atTime,
					By:     changes.Detail.by,
					Owner:  changes.Detail.owner,
				},
			},
			{
				Name:     changes_2.Name,
				Id:       changes_2.Iden,
				Place:    changes_2.Placed,
				Time:     timestamppb.New(changes_2.inTime),
				NewPlace: changes.Placed,
				Detail: &pb.Detail{
					From:   changes_2.Detail.From,
					AtTime: changes_2.Detail.atTime,
					By:     changes_2.Detail.by,
					Owner:  changes_2.Detail.owner,
				},
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
		var detail = fmt.Sprintf("%v:%v:%v is switched with %v at %v", peerOD, len(s.log), changes.Iden, changes_2.Iden, time.Now().Format(time.ANSIC))
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
	if changes.Name == "x" && new_place == -1 {
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

	if changes.Name == changes_2.Name {
		return false
	}

	if changes.Name == "x" || changes_2.Name == "x" {
		return false
	}

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
	// fmt.Println("Here")
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

	var rv string = ""
	for k, v := range s.cache {
		if v.Iden == x {
			if place == -1 {
				(s.cache)[k].Placed = int32(place)
				rv = string(fmt.Sprintf("%v is moved to %v at %v", v.Name, place, time.Now().Format(time.ANSIC)))
			} else {
				for i, j := range s.cache {
					if j.Placed == int32(place) {
						if j.Iden != x {
							// (s.cache)[i] = Container{
							// 	Iden:   j.Iden,
							// 	Name:   j.Name,
							// 	Placed: v.Placed,
							// 	Detail: j.Detail,
							// }
							(s.cache)[i].Placed = v.Placed
							rv = string(fmt.Sprintf("%v is switched with %v at %v", v.Name, j.Name, time.Now().Format(time.ANSIC)))
						}
					}
				}
			}
			// (s.cache)[k] = Container{
			// 	Iden:   v.Iden,
			// 	Name:   v.Name,
			// 	// Key:    v.Key,
			// 	Placed: int32(place),
			// 	Detail: v.Detail,
			// }
			(s.cache)[k].Placed = int32(place)
			if len(rv) < 1 {
				rv = string(fmt.Sprintf("%v is moved to %d at %v", v.Name, place, time.Now().Format(time.ANSIC)))
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
				Detail: detail{
					by:     "Planner A",
					From:   "USA",
					atTime: "12/5/2022",
					owner:  "Chevorale",
				},
			},
			{
				Name:   "2",
				Placed: -1,
				Iden:   "2",
				Key:    1,
				inTime: time.Now(),
				Detail: detail{
					by:     "Planner B",
					From:   "BR",
					atTime: "12/5/2022",
					owner:  "Rio Cargo",
				},
			},
			{
				Name:   "3",
				Placed: -1,
				Iden:   "3",
				Key:    2,
				inTime: time.Now(),
				Detail: detail{
					by:     "Planner C",
					From:   "CEZ",
					atTime: "12/5/2022",
					owner:  "Winston Housing",
				},
			},
			{
				Name:   "4",
				Placed: -1,
				Iden:   "4",
				Key:    3,
				inTime: time.Now(),
				Detail: detail{
					by:     "Planner D",
					From:   "JP",
					atTime: "12/5/2022",
					owner:  "Mitsubishi Elc",
				},
			},
			{
				Name:   "5",
				Placed: -1,
				Iden:   "5",
				Key:    4,
				inTime: time.Now(),
				Detail: detail{
					by:     "Planner E",
					From:   "SA",
					atTime: "12/5/2022",
					owner:  "SA trade assoc",
				},
			},
			{
				Name:   "6",
				Placed: -1,
				Iden:   "6",
				Key:    5,
				inTime: time.Now(),
				Detail: detail{
					by:     "Planner F",
					From:   "CN",
					atTime: "12/5/2022",
					owner:  "North Start inc",
				},
			},
			{
				Name:   "x",
				Placed: 6,
				Iden:   "7",
				Key:    6,
				inTime: time.Now(),
				Detail: detail{
					by:     "Ship",
					From:   "Ship",
					atTime: "12/5/2022",
					owner:  "Ship owner",
				},
			},
			{
				Name:   "x",
				Placed: 7,
				Iden:   "8",
				Key:    7,
				inTime: time.Now(),
				Detail: detail{
					by:     "Ship",
					From:   "Ship",
					atTime: "12/5/2022",
					owner:  "Ship owner",
				},
			},
			{
				Name:   "x",
				Placed: 8,
				Iden:   "9",
				Key:    8,
				inTime: time.Now(),
				Detail: detail{
					by:     "Ship",
					From:   "Ship",
					atTime: "12/5/2022",
					owner:  "Ship owner",
				},
			},
			{
				Name:   "x",
				Placed: 9,
				Iden:   "10",
				Key:    9,
				inTime: time.Now(),
				Detail: detail{
					by:     "Ship",
					From:   "Ship",
					atTime: "12/5/2022",
					owner:  "Ship owner",
				},
			},
			{
				Name:   "x",
				Placed: 10,
				Iden:   "11",
				Key:    10,
				inTime: time.Now(),
				Detail: detail{
					by:     "Ship",
					From:   "Ship",
					atTime: "12/5/2022",
					owner:  "Ship owner",
				},
			},
			{
				Name:   "x",
				Placed: 11,
				Iden:   "12",
				Key:    11,
				inTime: time.Now(),
				Detail: detail{
					by:     "Ship",
					From:   "Ship",
					atTime: "12/5/2022",
					owner:  "Ship owner",
				},
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
