package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth"
	"github.com/nandanurseptama/golang-grafana-otel/services/frontend/bootstrap"
	"github.com/nandanurseptama/golang-grafana-otel/services/frontend/pkg/otel"
	"go.opentelemetry.io/otel/codes"
)

func run() (err error) {
	runtime.GOMAXPROCS(2)
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

	meters, err := otel.SetupMeters(ctx, otel.Meter)

	if err != nil {
		return
	}
	authSvcClient, err := bootstrap.AuthServiceClient(ctx, env.AuthServiceAddress)

	if err != nil {
		return
	}

	// sleep for 20 seconds
	// for waiting all container services booted up
	time.Sleep(time.Second * 20)

	go doLogin(authSvcClient, meters)
	go authMe(authSvcClient, meters)
	go register(authSvcClient, meters)

	<-ctx.Done()
	stop()
	fmt.Println("stop")
	return
}

var authToken = ""

func getDelaySeconds() time.Duration {
	// Random delay seconds
	maxDelaySeconds := 300
	minDelaySeconds := 60
	delaySecond := rand.IntN(maxDelaySeconds-minDelaySeconds) + minDelaySeconds

	return time.Duration(delaySecond)
}

func getNumOfRequests() int {
	maxRequest := 100
	minRequest := 1
	numRequest := rand.IntN(maxRequest-minRequest) + minRequest

	return numRequest
}
func authMe(authClient auth.AuthServiceClient, meters *otel.OtelMeters) {

	doJob := func() {
		ctx := context.Background()
		ctx, span := otel.Tracer.Start(ctx, "authMe")
		defer span.End()

		_, err := authClient.Me(ctx, &auth.MeRequest{
			Token: authToken,
		})

		if err != nil {
			slog.ErrorContext(ctx, "failed authMe", slog.Any("reason", err.Error()))
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed authMe")
			meters.FailedMeCounter.Add(ctx, 1)
			return
		}

		slog.InfoContext(ctx, "success authMe")
		meters.SuccessMeCounter.Add(ctx, 1)
		span.SetStatus(codes.Ok, "success authMe")
	}

	for {
		delaySecond := getDelaySeconds()
		numRequest := getNumOfRequests()
		fmt.Printf("Do %d auth job\n", numRequest)
		for i := 0; i < numRequest; i++ {
			go doJob()
		}
		fmt.Printf("Delay for %d seconds\n", delaySecond)
		time.Sleep(time.Second * time.Duration(delaySecond))
		fmt.Println("Delay end")
	}

}
func doLogin(authClient auth.AuthServiceClient, meters *otel.OtelMeters) {

	doJob := func() {
		loginSuccessCase := rand.IntN(2) == 1
		email, password := func() (string, string) {
			if loginSuccessCase {
				return "doe@gmail.com", "12345"
			}
			return faker.Email(), faker.Password()
		}()

		ctx := context.Background()
		ctx, span := otel.Tracer.Start(ctx, "doLogin")

		authResult, err := authClient.Login(ctx, &auth.LoginRequest{
			Email:    email,
			Password: password,
		})

		if err != nil {
			authToken = ""
			slog.ErrorContext(ctx, "failed doLogin", slog.Any("reason", err.Error()))

			span.RecordError(err)
			span.SetStatus(codes.Error, "failed doLogin")

			meters.FailedLoginCounter.Add(ctx, 1)

			loginSuccessCase = !loginSuccessCase
			span.End()
			return
		}
		authToken = authResult.Token
		slog.InfoContext(ctx, "success doLogin")

		meters.SuccessLoginCounter.Add(ctx, 1)

		span.SetStatus(codes.Ok, "success doLogin")
		loginSuccessCase = !loginSuccessCase
		span.End()
	}
	var mtx sync.Mutex
	var wg sync.WaitGroup
	for {
		delaySecond := getDelaySeconds()
		numRequest := getNumOfRequests()
		fmt.Printf("Do %d login job\n", numRequest)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < numRequest; i++ {
				mtx.Lock()
				doJob()
				mtx.Unlock()
			}
		}()
		wg.Wait()
		fmt.Printf("Delay login job for %d seconds\n", delaySecond)
		time.Sleep(time.Second * time.Duration(delaySecond))
		fmt.Println("Delay login job end")
	}
}

func register(authClient auth.AuthServiceClient, meters *otel.OtelMeters) {

	doJob := func() {
		registerSuccessCase := rand.IntN(2) == 1
		email, password := func() (string, string) {
			if !registerSuccessCase {
				return "doe@gmail.com", "12345"
			}
			return faker.Email(), faker.Password()
		}()

		ctx := context.Background()
		ctx, span := otel.Tracer.Start(ctx, "doRegister")

		authResult, err := authClient.Register(ctx, &auth.LoginRequest{
			Email:    email,
			Password: password,
		})

		if err != nil {
			authToken = ""
			slog.ErrorContext(ctx, "failed doRegister", slog.Any("reason", err.Error()))

			span.RecordError(err)
			span.SetStatus(codes.Error, "failed doRegister")

			meters.FailedRegisterCounter.Add(ctx, 1)

			registerSuccessCase = !registerSuccessCase
			span.End()
			return
		}
		authToken = authResult.Token
		slog.InfoContext(ctx, "success doRegister")

		meters.SuccessRegisterCounter.Add(ctx, 1)

		span.SetStatus(codes.Ok, "success doRegister")
		registerSuccessCase = !registerSuccessCase
		span.End()
	}
	var mtx sync.Mutex
	var wg sync.WaitGroup
	for {
		delaySecond := getDelaySeconds()
		numRequest := getNumOfRequests()
		fmt.Printf("Do %d register job\n", numRequest)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < numRequest; i++ {
				mtx.Lock()
				doJob()
				mtx.Unlock()
			}
		}()
		wg.Wait()
		fmt.Printf("Delay register job for %d seconds\n", delaySecond)
		time.Sleep(time.Second * time.Duration(delaySecond))
		fmt.Println("Delay register job end")
	}
}
