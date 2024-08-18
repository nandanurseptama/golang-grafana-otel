package otel

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/metric"
)

type OtelMeters struct {
	FailedLoginCounter     metric.Int64Counter
	SuccessLoginCounter    metric.Int64Counter
	FailedRegisterCounter  metric.Int64Counter
	SuccessRegisterCounter metric.Int64Counter
	FailedMeCounter        metric.Int64Counter
	SuccessMeCounter       metric.Int64Counter
}

func SetupMeters(ctx context.Context, meter metric.Meter) (*OtelMeters, error) {
	failedLoginCounter, err := meter.Int64Counter("failed_login")

	if err != nil {
		return nil, errors.Join(errors.New("failed to setup failed login counter"), err)
	}

	successLoginCounter, err := meter.Int64Counter("success_login")

	if err != nil {
		return nil, errors.Join(errors.New("failed to setup success login counter"), err)
	}

	failedMeCounter, err := meter.Int64Counter("failed_me")

	if err != nil {
		return nil, errors.Join(errors.New("failed to setup failed me counter"), err)
	}

	successMeCounter, err := meter.Int64Counter("success_me")

	if err != nil {
		return nil, errors.Join(errors.New("failed to setup success me counter"), err)
	}

	failedRegisterCounter, err := meter.Int64Counter("failed_register")

	if err != nil {
		return nil, errors.Join(errors.New("failed to setup failed register counter"), err)
	}

	successRegisterCounter, err := meter.Int64Counter("success_register")

	if err != nil {
		return nil, errors.Join(errors.New("failed to setup success register counter"), err)
	}

	return &OtelMeters{
		FailedLoginCounter:     failedLoginCounter,
		SuccessLoginCounter:    successLoginCounter,
		FailedMeCounter:        failedMeCounter,
		SuccessMeCounter:       successMeCounter,
		FailedRegisterCounter:  failedRegisterCounter,
		SuccessRegisterCounter: successRegisterCounter,
	}, nil
}
