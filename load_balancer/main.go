package main

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tanovicm/tenderly/communication"
	grpc "google.golang.org/grpc"
)

func launchRouter(server *server) {

	r := mux.NewRouter()
	registerAuthRoutes(r)
	registerExpenseRoutes(r, server)
	registerBankRoutes(r, server)

	if err := http.ListenAndServe(":8090", r); err != nil {
		log.Fatalf("failed to serve http: %v", err)
	}
}

func launchLoadBalancerServer(server *server) {

	s := grpc.NewServer()
	communication.RegisterLoadBalancerServer(s, server)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}

func main() {

	server := server{}

	go launchRouter(&server)
	go launchLoadBalancerServer(&server)

	for {
	}
}
