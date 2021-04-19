package main

import (
	"flag"
	"log"

	"github.com/qibobo/grpcgo/auth"
	"github.com/qibobo/grpcgo/auth/interceptor"
	"github.com/qibobo/grpcgo/client"
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
	authClient := auth.NewAuthClient(conn, "qibobo", "123456")
	clientAuthInterceptor, err := interceptor.NewClientInterceptor(authClient, authMethods(), auth.TokenExpiredDuration)
	if err != nil {
		log.Panicf("can not initialize auth interceptor %s\n", err)
	}
	conn, err = grpc.Dial(*address, grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(clientAuthInterceptor.Unary()),
		grpc.WithChainStreamInterceptor(clientAuthInterceptor.Stream()),
	)
	if err != nil {
		log.Panicf("can not dial to gprc server %s\n", err)
	}

	personServiceClient := client.NewPersonServiceClient(conn)
	personId, err := personServiceClient.SavePerson(&models.Person{
		Name:  "qibobo",
		Email: "lqiyangl@gmail.com",
	})
	if err != nil {
		log.Printf("failed to save person %s\n", err)
		return
	}
	log.Printf("save person successfully %s", personId)

	personServiceClient.GerPersonStream()
	personServiceClient.UploadImage("12345", "./images/qiye.jpg")

}
func authMethods() map[string]bool {
	rootPath := "/grpcgo.rpcservice.PersonService/"
	return map[string]bool{
		rootPath + "SavePerson":      true,
		rootPath + "GetPerson":       false,
		rootPath + "GetPersonStream": true,
		rootPath + "UploadImage":     true,
	}
}
