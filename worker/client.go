package main

import (
	"context"
	"log"
	"time"

	"github.com/tanovicm/tenderly/communication"
)

func register(c communication.LoadBalancerClient, port string) {
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.Register(ctx, &communication.RegisterRequest{Addr: "localhost:" + port})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("Registered worker")
}

func deregister(c communication.LoadBalancerClient, port string) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.DeRegister(ctx, &communication.DeRegisterRequest{Addr: "localhost:" + port})
	if err != nil {
		log.Fatalf("could not de register: %v", err)
	}

	log.Printf("De Registered worker")
}
