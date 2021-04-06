package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"

	"github.com/qibobo/grpcgo/models"
	"google.golang.org/grpc"
)

func main() {

	address := flag.String("address", "", "gprc server address")
	flag.Parse()
	log.Printf("grpc server address is %s\n", *address)

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Panicf("can not dial to gprc server %s\n", err)
	}

	personClient := models.NewPersonServiceClient(conn)
	resp, err := personClient.SavePerson(context.Background(), &models.SavePersonRequest{
		Person: &models.Person{
			Name:  "qibobo",
			Email: "lqiyangl@gmail.com",
		},
	})
	if err != nil {
		log.Printf("failed to save person %s\n", err)
		return
	}
	log.Printf("save person successfully %s", resp.GetId())
	stream, err := personClient.GetPersonStream(context.Background(), &models.GetPersonStreamRequest{})
	if err != nil {
		log.Panicf("failed to save person %s\n", err)
		return
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("EOF %s\n", err)
				return
			}
			log.Printf("failed to get person from stream %s\n", err)

		}
		log.Printf("get person from stream %v\n", resp)
	}

}
