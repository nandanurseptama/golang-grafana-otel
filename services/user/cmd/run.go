package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/nandanurseptama/golang-grafana-otel/services/user"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/bootstrap"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/internal"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/pkg/otel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := otel.SetupSDK(ctx)

	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	env, err := bootstrap.GetEnv()
	if err != nil {
		return
	}

	db, err := bootstrap.OpenDB(env.DBPath)
	if err != nil {
		return
	}
	address := fmt.Sprintf(":%s", env.Port)

	lis, err := net.Listen("tcp", address)

	if err != nil {
		return
	}
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	serverImpl, err := internal.NewServer(db)
	if err != nil {
		return
	}

	user.RegisterUserServiceServer(s, serverImpl)
	log.Printf("server listening at %v", lis.Addr())

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- s.Serve(lis)
	}()

	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}
	return
}
