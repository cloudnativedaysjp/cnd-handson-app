package telemetry

import (
	"context"
	"errors"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

// WaitAndStop は SIGTERM / SIGINT か serveErr（サーバーの異常終了）を待ち、stops を timeout 内で順に呼ぶ
func WaitAndStop(timeout time.Duration, serveErr <-chan error, stops ...func(context.Context) error) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	var err error
	select {
	case <-ctx.Done():
	case err = <-serveErr:
	}
	cancel() // 通知を解除し、停止処理中に 2 回目のシグナルが来たら既定どおり即終了させる

	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for _, stop := range stops {
		err = errors.Join(err, stop(ctx))
	}
	return err
}

// GRPCStop は GracefulStop し、ctx の期限までに終わらなければ Stop で打ち切る
func GRPCStop(s *grpc.Server) func(context.Context) error {
	return func(ctx context.Context) error {
		done := make(chan struct{})
		go func() { s.GracefulStop(); close(done) }()
		select {
		case <-done:
		case <-ctx.Done():
			s.Stop()
		}
		return nil
	}
}
