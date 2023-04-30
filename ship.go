package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/HaiHart/ShipdockServer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
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

type ShipRelayConn struct {
	conn      *grpc.ClientConn
	client    pb.ComClient
	reqStream pb.Com_MoveShipClient
}

type ShipConn struct {
	pb.UnimplementedComServer
	context     context.Context
	cancel      context.CancelFunc
	ships       []Ship
	docks       []Doc
	clients     map[string]ShipRelayConn
	toSend      map[string]chan *pb.PlaceShip
	lock        sync.Mutex
	currCommand []Container
	port        int32
	log         []string
	detailLog   []string
}

func (s *ShipConn) FetchDocks(ctx context.Context, time *pb.Header) (*pb.Docks, error) {
	fmt.Println(" From Docks Side ", time.Time)

	var list []*pb.Ship

	var list_2 []*pb.Doc

	var log []string = make([]string, 0)

	for _, v := range s.ships {
		list = append(list, &pb.Ship{
			Name:   v.Name,
			Placed: int32(v.Placed),
			Iden:   v.Iden,
			Key:    v.Key,
			Detail: &pb.Detail{
				From:   v.Detail.From,
				By:     v.Detail.by,
				AtTime: v.Detail.atTime,
				Owner:  v.Detail.owner,
			},
			InTime:  timestamppb.New(v.InTime),
			OutTime: timestamppb.New(v.OutTime),
			Length:  v.Length,
		})
	}
	for _, v := range s.log {
		log = append(log, v)
	}

	for _, v := range s.docks {
		list_2 = append(list_2, &pb.Doc{
			No:           int32(v.No),
			Name:         v.Name,
			Length:       v.Length,
			BoarderRight: int32(v.BoarderRight),
			ShipList:     append([]string{}, v.ShipList...),
		})
	}

	return &pb.Docks{
		Log:   log,
		Ships: list,
		Docks: list_2,
	}, nil
}

func (s *ShipConn) checkShipList(sl [][]string) bool {

	if len(sl) != len(s.docks) {
		return false
	}
	for i := 0; i < len(sl); i++ {
		if len(sl[i]) != len(s.docks[i].ShipList) {
			return false
		}
		for j := 0; j < len(s.docks[i].ShipList); j++ {
			if sl[i][j] != s.docks[i].ShipList[j] {
				return false
			}
		}
	}

	return true
}

func (s *ShipConn) getShip(Name string) int {
	for k, i := range s.ships {
		if i.Name == Name {
			return k
		}
	}
	return -1
}

func (s *ShipConn) getDock(DocPlace int) int {
	for k, i := range s.docks {
		if i.No == (DocPlace) {
			return k
		}
	}
	return -1
}

func (s *ShipConn) checkFit(doc int, ship int) []int {
	idx := doc
	count := make([]int, 1)
	count[0] = s.docks[idx].No
	Length := s.ships[ship].Length
	for Length > s.docks[idx].Length {
		Length -= s.docks[idx].Length
		if s.docks[idx].BoarderRight != -1 {
			idx = s.getDock(int(s.docks[idx].BoarderRight))
		} else {
			return make([]int, 0)
		}
		count = append(count, int(s.docks[idx].No))

	}
	fmt.Println(count)
	return count
}

func (s *ShipConn) checkTime(doc int, ship int) int {
	out := s.ships[ship].OutTime
	in := s.ships[ship].InTime
	dock := s.docks[doc]
	for i := len(dock.ShipList) - 1; i >= 0; i-- {
		temp := s.ships[s.getShip(dock.ShipList[i])]
		if temp.OutTime.Before(in) {
			return i + 1
		} else {
			if temp.InTime.Before(out) {
				fmt.Println(doc, " ", dock.ShipList)
				fmt.Println("Error here")
				return -1
			}
		}
	}
	return 0
}

func (s *ShipConn) setShip(doc int, name string, idx int) {
	if len(s.docks[doc].ShipList) <= idx {
		temp := make([]string, idx+1)
		for i, k := range s.docks[doc].ShipList {
			temp[i] = k
		}
		temp[idx] = name
		s.docks[doc].ShipList = temp
		return
	}
	// if s.docks[doc].ShipList[idx] == "" {
	// 	s.docks[doc].ShipList[idx] = name
	// 	return
	// }
	s.docks[doc].ShipList = append(s.docks[doc].ShipList[:idx+1], s.docks[doc].ShipList[idx:]...)
	s.docks[doc].ShipList[idx] = name
	// i := s.docks[doc].BoarderRight
	// for i != -1 {
	// 	for rv, j := range s.docks[i].ShipList {
	// 		checkpoint := false
	// 		for _, k := range s.docks[doc].ShipList[idx+1:] {
	// 			if j == k {
	// 				s.setShip(int(i), j, rv+1)
	// 				checkpoint = true
	// 				break
	// 			}
	// 		}
	// 		if checkpoint {

	// 			break
	// 		}
	// 	}
	// 	i = s.docks[i].BoarderRight
	// }

}

func (s *ShipConn) placeShip(DocPlace int, Name string) bool {
	fmt.Println(DocPlace, "  ", Name)
	if DocPlace == -1 {
		return s.removeShip(Name)
	}
	var rv string = ""
	ship := s.getShip(Name)
	doc := s.getDock(DocPlace)
	if s.ships[ship].Placed == DocPlace {
		return false
	}
	temp := s.ships[ship].Placed
	if temp != -1 {
		s.removeShip(Name)
	}
	listDoc := s.checkFit(doc, ship)
	idxes := make([]int, 0)
	if len(listDoc) == 0 {
		s.placeShip(int(temp), Name)
		fmt.Println("End no place")
		return false
	}
	for _, i := range listDoc {
		j := s.checkTime(i, ship)
		idxes = append(idxes, j)
	}

	for _, i := range idxes {
		if i == -1 {
			s.placeShip(int(temp), Name)
			fmt.Println("End no place")
			return false
		}
	}

	s.ships[ship].Placed = doc
	for k, i := range listDoc {
		s.setShip(i, Name, idxes[k])
	}
	if len(listDoc) == 0 || len(idxes) == 0 {
		fmt.Println("empty")
		return false
	}
	rv = fmt.Sprintf("Ship %s has been set to dock(s): %v ; at position %v ; at time %v", Name, listDoc, idxes, time.Now())
	s.log = append(s.log, rv)
	fmt.Println(rv)
	return true
}

func (s *ShipConn) removeElement(DocPlace int, Name string) {
	for i, k := range s.docks[DocPlace].ShipList {
		if k == Name {
			s.docks[DocPlace].ShipList = append(s.docks[DocPlace].ShipList[:i], s.docks[DocPlace].ShipList[i+1:]...)
			return
		}
	}
}

func (s *ShipConn) removeShip(Name string) bool {
	ship := s.getShip(Name)
	if s.ships[ship].Placed == -1 {
		return false
	}
	s.ships[ship].Placed = -1
	for i, _ := range s.docks {
		s.removeElement(i, Name)
	}
	rv := fmt.Sprintf("Removed ship %s at time %v", Name, time.Now())
	s.log = append(s.log, rv)
	return true
}

func (s *ShipConn) parserTime(raw string) time.Time {
	raw = strings.Replace(raw, "Z", "", -1)
	major := strings.Split(raw, "T")
	var date []string = strings.Split(major[0], "-")
	var t []string = strings.Split(major[1], ":")
	datei := make([]int, 0)
	for _, i := range date {
		temp, _ := strconv.Atoi(i)
		datei = append(datei, temp)
	}
	for _, i := range t {
		temp, _ := strconv.Atoi(i)
		datei = append(datei, temp)
	}
	var rv time.Time = time.Date(datei[0], time.Month(datei[1]), datei[2], datei[3], datei[4], 0, 0, time.UTC)
	return rv
}

func (s *ShipConn) setTime(Name string, in time.Time, out time.Time) string {

	if in.Equal(time.Time{}) || out.Equal(time.Time{}) {
		return "fail"
	}
	if !in.Before(out) {
		return "Wrong Time order"
	}
	ship := s.getShip(Name)
	s.ships[ship].InTime = in
	s.ships[ship].OutTime = out
	s.log = append(s.log, fmt.Sprintf("Changed time to ship %s at time %v", Name, time.Now()))
	return "success"
}

func (s *ShipConn) propergate(newMove pb.PlaceShip) {
	for _, i := range s.toSend {
		i <- &newMove
	}
}

func (s *ShipConn) MoveShip(msg pb.Com_MoveShipServer) error {
	ctx := msg.Context()
	p, _ := peer.FromContext(ctx)
	peerID := p.Addr.String()

	if _, ok := s.toSend[peerID]; !ok {
		s.toSend[peerID] = make(chan *pb.PlaceShip, 1000)
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

		s.lock.Lock()

		if in.ChangeTime {
			var name = in.Ship.Name
			var inTime = in.Ship.InTime.AsTime()
			var outTime time.Time = in.Ship.OutTime.AsTime()
			rv := s.setTime(name, inTime, outTime)
			if rv == "success" {
				var res = &pb.PlaceShip{
					Err: "success",
					Ship: &pb.Ship{
						Name:    name,
						InTime:  timestamppb.New(inTime),
						OutTime: timestamppb.New(outTime),
					},
					ChangeTime: true,
				}
				s.propergate(*res)
				s.lock.Unlock()
				continue
			}
			var res = &pb.PlaceShip{
				Err: rv,
			}
			s.toSend[peerID] <- res
			s.lock.Unlock()
			continue
		}
		var changes = Ship{
			Name:   in.Ship.Name,
			Placed: (int(in.Ship.Placed)),
			Key:    in.Ship.Key,
			Iden:   in.Ship.Iden,
			InTime: in.Ship.InTime.AsTime(),
			Detail: detail{
				From:   in.Ship.Detail.From,
				atTime: in.Ship.Detail.AtTime,
				by:     in.Ship.Detail.By,
				owner:  in.Ship.Detail.Owner,
			},
			OutTime: in.Ship.OutTime.AsTime(),
			Length:  in.Ship.Length,
		}
		var shipList [][]string
		for _, i := range in.ShipList {
			shipList = append(shipList, i.List)
		}
		if s.checkShipList(shipList) && s.ships[s.getShip(in.Ship.Name)].Placed == int(in.Ship.Placed) {
			if s.placeShip(int(in.Place), changes.Name) {
				var res = &pb.PlaceShip{
					Ship: &pb.Ship{
						Name:    in.Ship.Name,
						Iden:    in.Ship.Iden,
						Key:     in.Ship.Key,
						Placed:  in.Ship.Placed,
						Length:  in.Ship.Length,
						InTime:  in.Ship.GetInTime(),
						OutTime: in.Ship.GetOutTime(),
						Detail:  in.Ship.GetDetail(),
					},
					Place:      in.GetPlace(),
					ChangeTime: false,
				}
				s.propergate(*res)
			}
		} else {
			var res = pb.PlaceShip{
				Err: "desync",
			}
			s.toSend[peerID] <- &res
		}
		s.lock.Unlock()
	}

	return nil
}

func (s *ShipConn) Start() error {
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

func (s *ShipConn) RunServer() error {
	defer s.cancel()
	var err error = s.Start()

	return err
}
