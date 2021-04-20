package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/qibobo/grpcgo/auth"
	"github.com/qibobo/grpcgo/auth/interceptor"
	"github.com/qibobo/grpcgo/models"
	"github.com/qibobo/grpcgo/service"
	"github.com/qibobo/grpcgo/service/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	serverCert = "cert/server-cert.pem"
	serverKey  = "cert/server-key.pem"
	clientCA   = "cert/ca-cert.pem"
)

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()
	log.Printf("the listening port is %d\n", *port)

	personServer := service.NewPersonServer(store.NewInMemoryStore(), &store.DiskImageStore{})
	userStore := store.NewInMemoryUserStore()
	u, err := models.NewUser("qibobo", "123456", "admin")
	if err != nil {
		log.Panicf("can not new user %s\n", err)
	}
	err = userStore.Save(u)
	if err != nil {
		log.Panicf("can not save user %s\n", err)
	}
	jwtManager := auth.NewJWTManager("secretkey", time.Hour)

	userServer := service.NewAuthServer(userStore, jwtManager)
	autoInterceptor := interceptor.NewServerInterceptor(accessiableRoles(), jwtManager)

	//load server cert
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Panicf("can not load cert %s\n", err)
	}
	//load client ca

	clientCA, err := ioutil.ReadFile(clientCA)
	if err != nil {
		log.Panicf("can not load client ca cert %s\n", err)
	}
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(clientCA)
	if !ok {
		log.Panicf("can not append client ca cert %s\n", err)
	}
	rpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(autoInterceptor.Unary()),
		grpc.ChainStreamInterceptor(autoInterceptor.Stream()),
		grpc.Creds(credentials.NewTLS(
			&tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientAuth:   tls.RequireAndVerifyClientCert,
				ClientCAs:    certPool,
			}),
		),
	)
	models.RegisterPersonServiceServer(rpcServer, personServer)
	models.RegisterLoginServiceServer(rpcServer, userServer)
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

func accessiableRoles() map[string][]string {
	rootPath := "/grpcgo.rpcservice.PersonService/"
	return map[string][]string{
		rootPath + "SavePerson":      {"admin"},
		rootPath + "GetPerson":       {"admin", "user"},
		rootPath + "GetPersonStream": {"admin", "user"},
		rootPath + "UploadImage":     {"admin"},
	}
}
