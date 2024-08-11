package main

import (
	"fmt"
	"log"
	"net"

	"github.com/nandanurseptama/golang-grafana-otel/services/user"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/bootstrap"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/internal"
	"google.golang.org/grpc"
)

func main() {
	env, err := bootstrap.GetEnv()
	if err != nil {
		panic(err)
	}

	db, err := bootstrap.OpenDB(env.DBPath)
	if err != nil {
		panic(err)
	}
	address := fmt.Sprintf(":%s", env.Port)

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	serverImpl, err := internal.NewServer(db)
	if err != nil {
		panic(err)
	}

	user.RegisterUserServiceServer(s, serverImpl)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
