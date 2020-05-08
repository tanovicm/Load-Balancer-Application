package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tanovicm/tenderly/communication"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	loadBalancerAddress = "localhost:50051"
	defaultPort         = "60061"
)

func launchWorkerServer(server *server, port string) {

	s := grpc.NewServer()
	communication.RegisterWorkerServer(s, server)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}

func main() {

	// parsing args
	port := defaultPort
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	// Create and launch worker grpc server
	server := createWorkerServer()
	go launchWorkerServer(&server, port)

	// Set up a connection to the Load Balancer.
	conn, err := grpc.Dial(loadBalancerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect load balancer: %v", err)
	}
	defer conn.Close()
	
	c := communication.NewLoadBalancerClient(conn)

	// Deregister from the Load Balancer on SIGTERM
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		deregister(c, port)
		os.Exit(1)
	}()

	// Register to the Load Balancer when connection is ready (on LB (re)boot).
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		changed := conn.WaitForStateChange(ctx, conn.GetState())
		if !changed {
			continue
		}

		log.Printf("State changed %v", conn.GetState().String())

		if conn.GetState() == connectivity.Ready {
			register(c, port)
		}
	}
}
