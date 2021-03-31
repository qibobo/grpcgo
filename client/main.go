package main

import (
	"log"

	"google.golang.org/grpc"
)

const (
	address     = ":8080"
	defaultName = "world"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect to grpc server: %s", err)
	}
	defer conn.Close()

}
