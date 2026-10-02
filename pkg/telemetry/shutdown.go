package telemetry

import (
	"context"
	"errors"
	"os/signal"
	"syscall"
	"time"
)

// WaitAndStop は SIGTERM / SIGINT を受けるまで待ち、stops を timeout 内で順に呼ぶ
func WaitAndStop(timeout time.Duration, stops ...func(context.Context) error) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	<-ctx.Done()
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var errs []error
	for _, stop := range stops {
		errs = append(errs, stop(ctx))
	}
	return errors.Join(errs...)
}
