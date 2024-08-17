package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth"
	"github.com/nandanurseptama/golang-grafana-otel/services/frontend/bootstrap"
	"github.com/nandanurseptama/golang-grafana-otel/services/frontend/otel"
	otelSdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	otelShutdown, err := otel.SetupSDK(ctx)

	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	env, err := bootstrap.GetEnv(ctx)
	if err != nil {
		return
	}

	meter := otelSdk.Meter("github.com/nandanurseptama/golang-grafana-otel/services/frontend")
	meters, err := otel.SetupMeters(ctx, meter)

	if err != nil {
		return
	}
	authSvcClient, err := bootstrap.AuthServiceClient(ctx, env.AuthServiceAddress)

	if err != nil {
		return
	}
	go doLogin(authSvcClient, meters)

	<-ctx.Done()
	stop()
	return
}

func doLogin(authClient auth.AuthServiceClient, meters *otel.OtelMeters) {
	tracer := otelSdk.Tracer("github.com/nandanurseptama/golang-grafana-otel/services/frontend")
	var loginSuccessCase = false
	for {
		email, password := func() (string, string) {
			if loginSuccessCase {
				return "doe@gmail.com", "12345"
			}
			return faker.Email(), faker.Password()
		}()

		ctx := context.Background()
		ctx, span := tracer.Start(ctx, "doLogin")

		_, err := authClient.Login(ctx, &auth.LoginRequest{
			Email:    email,
			Password: password,
		})

		if err != nil {
			slog.ErrorContext(ctx, "failed doLogin", slog.Any("reason", err.Error()))

			span.RecordError(err)
			span.SetStatus(codes.Error, "failed doLogin")
			meters.FailedLoginCounter.Add(ctx, 1)
			loginSuccessCase = !loginSuccessCase
			span.End()
			time.Sleep(10 * time.Second)
			continue
		}
		slog.InfoContext(ctx, "success doLogin")
		meters.SuccessLoginCounter.Add(ctx, 1)
		span.SetStatus(codes.Ok, "success doLogin")
		loginSuccessCase = !loginSuccessCase
		span.End()
		time.Sleep(10 * time.Second)
	}
}
