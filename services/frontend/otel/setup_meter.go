package otel

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/metric"
)

type OtelMeters struct {
	FailedLoginCounter  metric.Int64Counter
	SuccessLoginCounter metric.Int64Counter
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

	return &OtelMeters{
		FailedLoginCounter:  failedLoginCounter,
		SuccessLoginCounter: successLoginCounter,
	}, nil
}
