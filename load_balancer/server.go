package main

import (
	"errors"
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tanovicm/tenderly/communication"
	grpc "google.golang.org/grpc"
)

type server struct {
	communication.UnimplementedLoadBalancerServer

	conns []*grpc.ClientConn
	next  int
}

func (s *server) Register(ctx context.Context, in *communication.RegisterRequest) (*empty.Empty, error) {

	log.Printf("Register: %v", in.GetAddr())

	conn, err := grpc.Dial(in.GetAddr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Dial failed")
		return nil, err
	}

	log.Printf("Connection established")
	s.conns = append(s.conns, conn)

	return &empty.Empty{}, nil
}

func (s *server) DeRegister(ctx context.Context, in *communication.DeRegisterRequest) (*empty.Empty, error) {

	log.Printf("De Register: %v", in.GetAddr())

	for i, conn := range s.conns {
		if conn.Target() == in.GetAddr() {
			s.conns = append(s.conns[:i], s.conns[i+1:]...)
			break
		}
	}

	return &empty.Empty{}, nil
}

func (s *server) GetWorker() (communication.WorkerClient, error) {

	if len(s.conns) == 0 {
		return nil, errors.New("No workers active")
	}

	worker := communication.NewWorkerClient(s.conns[s.next])
	s.next = (s.next + 1) % len(s.conns)
	return worker, nil
}
