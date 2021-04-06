package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/qibobo/grpcgo/models"
	"github.com/qibobo/grpcgo/service"
	"github.com/qibobo/grpcgo/service/store"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()
	log.Printf("the listening port is %d\n", *port)

	personServer := service.NewPersonServer(store.NewInMemoryStore())
	rpcServer := grpc.NewServer()
	models.RegisterPersonServiceServer(rpcServer, personServer)
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicf("can not listen server %s\n", err)
	}

	err = rpcServer.Serve(listener)
	if err != nil {
		log.Panicf("can not start grpc server %s\n", err)
	}

}
