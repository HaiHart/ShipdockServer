package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
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

type Cordinates struct {
	bay  int
	row  int
	tier int
}
type Container struct {
	Name   string
	Cor    Cordinates
	Type   int
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

type ShipContainer struct {
	containers []Container
	invalids   []Cordinates
	bays       int
	rows       int
	tiers      int
}

type SerConn struct {
	pb.UnimplementedComServer
	context     context.Context
	cancel      context.CancelFunc
	shipsList   map[string]ShipContainer
	clients     map[string]RelayConn
	toSend      map[string]chan *pb.Pack
	nameToShip  map[string]string
	lock        sync.Mutex
	currCommand []Container
	port        int32
	log         map[string][]string
	detailLog   []string
}

func (s *SerConn) FetchList(ctx context.Context, time *pb.Header) (*pb.ShipList, error) {
	fmt.Println(time.Time)

	var list []*pb.ContainerSet

	var log []string = make([]string, 0)

	var inval []*pb.Cordinate

	if _, ok := s.nameToShip[time.Name]; !ok {
		s.nameToShip[time.Name] = time.ShipId
	}

	if s.nameToShip[time.Name] != time.Name {
		s.nameToShip[time.Name] = time.Name
	}

	// peerID := time.Name

	// if _, ok := s.accessList[peerID]; !ok {
	// 	s.accessList[peerID] = time.ShipId
	// }

	for _, v := range s.shipsList[time.ShipId].containers {
		list = append(list, &pb.ContainerSet{
			Name: v.Name,
			Id:   v.Iden,
			Key:  v.Key,
			Place: &pb.Cordinate{
				Bay:  int32(v.Cor.bay),
				Row:  int32(v.Cor.row),
				Tier: int32(v.Cor.tier),
			},
			Detail: &pb.Detail{
				From:   v.Detail.From,
				By:     v.Detail.by,
				AtTime: v.Detail.atTime,
				Owner:  v.Detail.owner,
			},
			Type: int32(v.Type),
		})
	}
	for _, v := range s.log[time.ShipId] {
		log = append(log, v)
	}

	for _, v := range s.shipsList[time.ShipId].invalids {
		inval = append(inval, &pb.Cordinate{
			Bay:  int32(v.bay),
			Row:  int32(v.row),
			Tier: int32(v.tier),
		})
	}

	return &pb.ShipList{
		List:  list,
		Log:   log,
		Inval: inval,
		Sizes: &pb.Cordinate{
			Bay:  int32(s.shipsList[time.ShipId].bays),
			Row:  int32(s.shipsList[time.ShipId].rows),
			Tier: int32(s.shipsList[time.ShipId].tiers),
		},
	}, nil
}

func (s *SerConn) MoveContainer(msg pb.Com_MoveContainerServer) error {
	ctx := msg.Context()
	p, _ := peer.FromContext(ctx)
	peerID := p.Addr.String()

	if _, ok := s.toSend[peerID]; !ok {
		s.toSend[peerID] = make(chan *pb.Pack, 1000)
	}

	end:= false

	go func() {
		for {
			select {
			case toSend := <-s.toSend[peerID]:

				// if toSend.ShipName != s.nameToShip[] {
				// 	continue
				// }


				if err := msg.Send(toSend); err != nil {
					fmt.Println(err)
				}
				// fmt.Println("sent to ",peerID)
			case <-ctx.Done():
				delete(s.nameToShip, peerID)
				delete(s.toSend, peerID)
				end =true
				return 
			}
		}

	}()
	for {
		if(end){
			break
		}
		in, err := msg.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("got command at %v\n", time.Now().UTC())
		

		var changes = Container{
			Name: in.List[0].Name,
			Cor: Cordinates{
				bay:  int(in.List[0].Place.Bay),
				row:  int(in.List[0].Place.Row),
				tier: int(in.List[0].Place.Tier),
			},
			Key:    0,
			Iden:   in.List[0].Id,
			inTime: in.List[0].Time.AsTime(),
			Detail: detail{
				From:   in.List[0].Detail.From,
				atTime: in.List[0].Detail.AtTime,
				by:     in.List[0].Detail.By,
				owner:  in.List[0].Detail.Owner,
			},
			Type: int(in.List[0].Type),
		}
		var new_place = Cordinates{
			bay:  int(in.List[0].NewPlace.Bay),
			row:  int(in.List[0].NewPlace.Row),
			tier: int(in.List[0].NewPlace.Tier),
		}
		var swap = in.Swap
		if swap {
			var changes_2 = Container{
				Name: in.List[1].Name,
				Cor: Cordinates{
					bay:  int(in.List[1].Place.Bay),
					row:  int(in.List[1].Place.Row),
					tier: int(in.List[1].Place.Tier),
				},
				Key:    0,
				Iden:   in.List[1].Id,
				inTime: in.List[1].Time.AsTime(),
				Detail: detail{
					From:   in.List[1].Detail.From,
					atTime: in.List[1].Detail.AtTime,
					by:     in.List[1].Detail.By,
					owner:  in.List[1].Detail.Owner,
				},
				Type: int(in.List[1].Type),
			}
			s.ValidSwap(&changes, &changes_2, peerID, in.ShipName)
		}

		s.ValidMove(&changes, new_place, peerID, in.ShipName)
	}

	return nil
}

func (s *SerConn) ValidMove(changes *Container, new_place Cordinates, peerID string, shipName string) {

	if changes == nil {
		return
	}

	var new_move = &pb.Pack{
		List: []*pb.Container{
			{
				Name: changes.Name,
				Id:   changes.Iden,
				Place: &pb.Cordinate{
					Bay:  int32(changes.Cor.bay),
					Row:  int32(changes.Cor.row),
					Tier: int32(changes.Cor.tier),
				},
				Time: timestamppb.New(changes.inTime),
				NewPlace: &pb.Cordinate{
					Bay:  int32(new_place.bay),
					Row:  int32(new_place.row),
					Tier: int32(new_place.tier),
				},
				Detail: &pb.Detail{
					From:   changes.Detail.From,
					AtTime: changes.Detail.atTime,
					By:     changes.Detail.by,
					Owner:  changes.Detail.owner,
				},
			},
		},
		Swap:     false,
		Err:      "None",
		ShipName: shipName,
	}

	if s.CheckOnCacheMove(changes, new_place, shipName) {
		s.currCommand = append([]Container{*changes}, s.currCommand...)
		var detail = fmt.Sprintf("%v:%v:%v is moved to %v ay %v", peerID, len(s.log), changes.Iden, new_place, time.Now().Format(time.ANSIC))
		s.detailLog = append(s.detailLog, detail)
		s.Swap(changes.Iden, new_place, shipName)
		for _, i := range s.toSend {
			i <- new_move
			// fmt.Println(name)
		}
		// fmt.Println("not failed")

	} else {
		new_move.Err = "Miss match occur"
		s.toSend[peerID] <- new_move
		fmt.Println("failed")
	}

	return
}

func (s *SerConn) ValidSwap(changes *Container, changes_2 *Container, peerOD string, shipName string) {
	if changes == nil || changes_2 == nil {
		return
	}

	var new_move = &pb.Pack{
		List: []*pb.Container{
			{
				Name: changes.Name,
				Id:   changes.Iden,
				Place: &pb.Cordinate{
					Bay:  int32(changes.Cor.bay),
					Row:  int32(changes.Cor.row),
					Tier: int32(changes.Cor.tier),
				},
				Time: timestamppb.New(changes.inTime),
				NewPlace: &pb.Cordinate{
					Bay:  int32(changes_2.Cor.bay),
					Row:  int32(changes_2.Cor.row),
					Tier: int32(changes_2.Cor.tier),
				},
				Detail: &pb.Detail{
					From:   changes.Detail.From,
					AtTime: changes.Detail.atTime,
					By:     changes.Detail.by,
					Owner:  changes.Detail.owner,
				},
			},
			{
				Name: changes_2.Name,
				Id:   changes_2.Iden,
				Place: &pb.Cordinate{
					Bay:  int32(changes_2.Cor.bay),
					Row:  int32(changes_2.Cor.row),
					Tier: int32(changes_2.Cor.tier),
				},
				Time: timestamppb.New(changes_2.inTime),
				NewPlace: &pb.Cordinate{
					Bay:  int32(changes.Cor.bay),
					Row:  int32(changes.Cor.row),
					Tier: int32(changes.Cor.tier),
				},
				Detail: &pb.Detail{
					From:   changes_2.Detail.From,
					AtTime: changes_2.Detail.atTime,
					By:     changes_2.Detail.by,
					Owner:  changes_2.Detail.owner,
				},
			},
		},
		Swap:     true,
		Err:      "None",
		ShipName: shipName,
	}
	if s.CheckOnCacheSwap(changes, changes_2, shipName) {
		s.currCommand = append([]Container{*changes, *changes_2}, s.currCommand...)
		for _, i := range s.toSend {
			i <- new_move
		}
		var detail = fmt.Sprintf("%v:%v:%v is switched with %v at %v", peerOD, len(s.log), changes.Iden, changes_2.Iden, time.Now().Format(time.ANSIC))
		s.detailLog = append(s.detailLog, detail)
		s.Swap(changes.Iden, changes_2.Cor, shipName)
	} else {
		new_move.Err = "Miss match starting point"
		s.toSend[peerOD] <- new_move
	}

	return
}

func (s *SerConn) CheckOnCacheMove(changes *Container, new_place Cordinates, shipName string) bool {

	if changes == nil {
		return false
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	if new_place.bay > s.shipsList[shipName].bays || new_place.row > s.shipsList[shipName].rows || new_place.tier > s.shipsList[shipName].tiers {
		fmt.Println("failed _1")
		return false
	}

	if changes.Cor.bay == new_place.bay && changes.Cor.row == new_place.row && changes.Cor.tier == new_place.tier {
		fmt.Println("failed_2")
		return false
	}

	if changes.Type == 0 && new_place.bay%2 == 1 && new_place.bay != -1 {
		fmt.Println("failed_3")
		return false
	}

	if changes.Type == 1 && new_place.bay%2 == 0 && new_place.bay != -1 {
		fmt.Println("failed_4")
		return false
	}

	if changes.Name == "x" && new_place.bay == -1 {
		fmt.Println("failed_5")
		return false
	}
	for _, v := range s.shipsList[shipName].containers {
		if v.Cor.bay == new_place.bay && v.Cor.row == new_place.row && v.Cor.tier == new_place.tier && new_place.bay != -1 && new_place.row != -1 && new_place.tier != -1 {
			fmt.Println("failed_6 ", v.Cor, " ", new_place)
			return false
		}
	}
	for _, v := range s.shipsList[shipName].containers {
		if v.Iden == changes.Iden {
			if v.Cor.bay != changes.Cor.bay || v.Cor.row != changes.Cor.row || v.Cor.tier != changes.Cor.tier {
				fmt.Println("failed_7 ", v.Cor, " ", changes.Cor)
				return false
			}
		}
	}
	for _, v := range s.shipsList[shipName].invalids {
		if v.bay == changes.Cor.bay && v.row == changes.Cor.row && v.tier == changes.Cor.tier {
			fmt.Println("failed_8")
			return false
		}
	}
	return true
}

func (s *SerConn) CheckOnCacheSwap(changes *Container, changes_2 *Container, shipName string) bool {
	if changes == nil || changes_2 == nil {
		return false
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	if changes_2.Cor.bay == -1 {
		return false
	}

	if changes.Name == changes_2.Name {
		return false
	}

	if changes.Name == "x" || changes_2.Name == "x" {
		return false
	}

	if changes.Cor.row == changes_2.Cor.row && changes.Cor.tier == changes_2.Cor.tier {
		return false
	}
	for _, v := range s.shipsList[shipName].containers {
		if v.Iden == changes.Iden && (v.Cor.bay != changes.Cor.bay || v.Cor.row != changes.Cor.row || v.Cor.tier != changes.Cor.tier) {
			return false
		}
		if v.Iden == changes_2.Iden && (v.Cor.bay != changes_2.Cor.bay || v.Cor.row != changes_2.Cor.row || v.Cor.tier != changes_2.Cor.tier) {
			return false
		}
	}
	// for _, v := range s.shipsList[shipName].invalids {

	// }
	// for _, v := range s.shipsList[shipName].containers {
	// 	if v.Iden == changes.Iden {
	// 		v.Placed = changes_2.Placed
	// 	}
	// 	if v.Iden == changes_2.Iden {
	// 		v.Placed = changes.Placed
	// 	}
	// }
	return true
}

func (s *SerConn) FetchShip(ctx context.Context, msg *pb.ShipAccess) (*pb.ShipResponse, error) {
	return nil, nil
}

func (s *SerConn) generatePackage(shipName string) {
	tmp := s.shipsList[shipName]
	for i := 0; i < 3; i++ {
		for j := 0; j < tmp.tiers/3; j++ {
			for k := 0; k < tmp.rows; k++ {
				tmp.containers = append(tmp.containers, Container{
					Name: "x",
					Cor: Cordinates{
						bay:  i,
						row:  k,
						tier: j,
					},
					Type:   i % 2,
					Iden:   strconv.Itoa(7 + i*tmp.rows*tmp.tiers/3 + j*tmp.rows + k),
					Key:    int32(6 + i*tmp.rows*tmp.tiers/3 + k*tmp.rows + k),
					inTime: time.Now(),
					Detail: detail{
						by:     "Ship",
						From:   "Ship",
						atTime: "12/5/2022",
						owner:  "Ship owner",
					},
				})
			}
		}
	}
	s.shipsList[shipName] = tmp
}

func (s *SerConn) Swap(x string, place Cordinates, shipName string) {

	var rv string = ""
	for k, v := range s.shipsList[shipName].containers {
		if v.Iden == x {
			if place.bay == -1 {
				(s.shipsList[shipName].containers)[k].Cor = place
				rv = string(fmt.Sprintf("%v is moved to %v at %v", v.Name, place, time.Now().Format(time.ANSIC)))
			} else {
				for i, j := range s.shipsList[shipName].containers {
					if j.Cor.bay == place.bay && j.Cor.row == place.row && j.Cor.tier == place.tier {
						if j.Iden != x {

							(s.shipsList[shipName].containers)[i].Cor = v.Cor
							rv = string(fmt.Sprintf("%v is switched with %v at %v", v.Name, j.Name, time.Now().Format(time.ANSIC)))
						}
					}
				}
			}

			(s.shipsList[shipName].containers)[k].Cor = place
			if len(rv) < 1 {
				rv = string(fmt.Sprintf("%v is moved to %d at %v", v.Name, place, time.Now().Format(time.ANSIC)))
			}

		}
	}
	fmt.Println(rv, place)
	s.log[shipName] = append(s.log[shipName], rv)

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
	ships := ShipContainer{
		bays:  10,
		rows:  10,
		tiers: 6,
		invalids: []Cordinates{
			{
				bay:  0,
				row:  0,
				tier: 5,
			},
			{
				bay:  0,
				row:  9,
				tier: 5,
			},
			{
				bay:  1,
				row:  5,
				tier: 5,
			},
			{
				bay:  1,
				row:  4,
				tier: 5,
			},
			{
				bay:  2,
				row:  3,
				tier: 5,
			},
			{
				bay:  2,
				row:  6,
				tier: 5,
			},
		},
		containers: []Container{
			{
				Name: "1",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   0,
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
				Name: "2",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   1,
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
				Name: "3",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   0,
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
				Name: "4",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   1,
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
				Name: "5",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   0,
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
				Name: "6",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   1,
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
		},
	}
	ships_2 := ShipContainer{
		bays:  10,
		rows:  10,
		tiers: 6,
		invalids: []Cordinates{
			{
				bay:  0,
				row:  0,
				tier: 5,
			},
			{
				bay:  0,
				row:  9,
				tier: 5,
			},
			{
				bay:  1,
				row:  5,
				tier: 5,
			},
			{
				bay:  1,
				row:  4,
				tier: 5,
			},
			{
				bay:  2,
				row:  3,
				tier: 5,
			},
			{
				bay:  2,
				row:  6,
				tier: 5,
			},
		},
		containers: []Container{
			{
				Name: "1",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   0,
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
				Name: "2",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   1,
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
				Name: "3",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   0,
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
				Name: "4",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   1,
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
				Name: "5",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   0,
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
				Name: "6",
				Cor: Cordinates{
					bay:  -1,
					row:  -1,
					tier: -1,
				},
				Type:   1,
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
				Name: "x",
				Cor: Cordinates{
					bay:  0,
					row:  0,
					tier: 0,
				},
				Type:   0,
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
				Name: "x",
				Cor: Cordinates{
					bay:  0,
					row:  1,
					tier: 0,
				},
				Type:   0,
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
				Name: "x",
				Cor: Cordinates{
					bay:  0,
					row:  2,
					tier: 0,
				},
				Type:   0,
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
				Name: "x",
				Cor: Cordinates{
					bay:  1,
					row:  0,
					tier: 0,
				},
				Type:   1,
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
				Name: "x",
				Cor: Cordinates{
					bay:  1,
					row:  0,
					tier: 1,
				},
				Type:   1,
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
				Name: "x",
				Cor: Cordinates{
					bay:  1,
					row:  1,
					tier: 0,
				},
				Type:   1,
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
	}
	server := SerConn{
		port:       8080,
		context:    ctx,
		cancel:     cancel,
		shipsList:  make(map[string]ShipContainer),
		toSend:     make(map[string]chan *pb.Pack,100),
		clients:    make(map[string]RelayConn),
		log:        make(map[string][]string),
		detailLog:  make([]string, 0),
		nameToShip: make(map[string]string),
	}
	server.shipsList["Ship_1"] = ships
	server.shipsList["Ship_2"] = ships_2
	server.log["Ship_1"] = make([]string, 0)
	server.log["Ship_2"] = make([]string, 0)
	server.generatePackage("Ship_1")
	shipServer := ShipConn{
		port:      8050,
		context:   ctx,
		cancel:    cancel,
		toSend:    make(map[string]chan *pb.PlaceShip),
		clients:   make(map[string]ShipRelayConn),
		log:       make([]string, 0),
		detailLog: make([]string, 0),
		ships: []Ship{{
			Name:   "Ship 0",
			Placed: -1,
			Iden:   "0",
			Key:    0,
			Detail: detail{
				From:   "AF",
				atTime: "15/8/2022",
				by:     "Planner A",
				owner:  "Ris inc",
			},
			Length:  128,
			InTime:  time.Date(2023, time.February, 15, 18, 30, 0, 0, time.UTC),
			OutTime: time.Date(2023, time.February, 22, 18, 30, 0, 0, time.UTC),
		},
			{
				Name:   "Ship 1",
				Placed: -1,
				Iden:   "1",
				Key:    1,
				Detail: detail{
					From:   "AF_1",
					atTime: "15/8/2022",
					by:     "Planner A",
					owner:  "Ris inc",
				},
				Length:  250,
				InTime:  time.Date(2023, time.February, 20, 18, 30, 0, 0, time.UTC),
				OutTime: time.Date(2023, time.February, 27, 18, 30, 0, 0, time.UTC),
			},
			{
				Name:   "Ship 2",
				Placed: -1,
				Iden:   "2",
				Key:    2,
				Detail: detail{
					From:   "AF_2",
					atTime: "15/8/2022",
					by:     "Planner A",
					owner:  "Ris inc",
				},
				Length:  128,
				InTime:  time.Date(2023, time.February, 23, 18, 30, 0, 0, time.UTC),
				OutTime: time.Date(2023, time.March, 1, 18, 30, 0, 0, time.UTC),
			},
			{
				Name:   "Ship 3",
				Placed: -1,
				Iden:   "3",
				Key:    3,
				Detail: detail{
					From:   "AF_3",
					atTime: "15/8/2022",
					by:     "Planner A",
					owner:  "Ris inc",
				},
				Length:  128,
				InTime:  time.Date(2023, time.March, 5, 18, 30, 0, 0, time.UTC),
				OutTime: time.Date(2023, time.March, 12, 18, 30, 0, 0, time.UTC),
			},
			{
				Name:   "Ship 4",
				Placed: -1,
				Iden:   "4",
				Key:    4,
				Detail: detail{
					From:   "AF_2",
					atTime: "15/8/2022",
					by:     "Planner A",
					owner:  "Ris inc",
				},
				Length:  250,
				InTime:  time.Date(2023, time.March, 14, 18, 30, 0, 0, time.UTC),
				OutTime: time.Date(2023, time.March, 22, 18, 30, 0, 0, time.UTC),
			},
			{
				Name:   "Ship 5",
				Placed: -1,
				Iden:   "5",
				Key:    5,
				Detail: detail{
					From:   "AF_2",
					atTime: "15/8/2022",
					by:     "Planner A",
					owner:  "Ris inc",
				},
				Length:  128,
				InTime:  time.Date(2023, time.March, 23, 18, 30, 0, 0, time.UTC),
				OutTime: time.Date(2023, time.March, 30, 18, 30, 0, 0, time.UTC),
			},
		},
		docks: []Doc{
			{
				No:           0,
				Name:         "doc_0",
				Length:       200,
				BoarderRight: 1,
				ShipList:     make([]string, 0),
			},
			{
				No:           1,
				Name:         "doc_1",
				Length:       200,
				BoarderRight: 2,
				ShipList:     make([]string, 0),
			},
			{
				No:           2,
				Name:         "doc_2",
				Length:       200,
				BoarderRight: -1,
				ShipList:     make([]string, 0),
			},
			{
				No:           3,
				Name:         "doc_3",
				Length:       150,
				BoarderRight: 4,
				ShipList:     make([]string, 0),
			},
			{
				No:           4,
				Name:         "doc_4",
				Length:       200,
				BoarderRight: 5,
				ShipList:     make([]string, 0),
			},
			{
				No:           5,
				Name:         "doc_5",
				Length:       200,
				BoarderRight: -1,
				ShipList:     make([]string, 0),
			},
			{
				No:           6,
				Name:         "doc_6",
				Length:       300,
				BoarderRight: -1,
				ShipList:     make([]string, 0),
			},
			{
				No:           7,
				Name:         "doc_7",
				Length:       250,
				BoarderRight: -1,
				ShipList:     make([]string, 0),
			},
		}}
	var group errgroup.Group
	fmt.Println("here")
	group.Go(server.RunServer)
	group.Go(shipServer.RunServer)
	fmt.Println("Started")

	err := group.Wait()
	if err != nil {
		fmt.Println(err)
		return
	}
}
