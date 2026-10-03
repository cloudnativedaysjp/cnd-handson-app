package telemetry

import (
	"context"
	"errors"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

// WaitAndStop は SIGTERM / SIGINT か serveErr（サーバーの異常終了）を待ち、stops を順に呼ぶ。
// 前の停止が期限を使い切っても後ろ（telemetry の flush など）が動けるよう、期限は stop ごとに与える
func WaitAndStop(timeout time.Duration, serveErr <-chan error, stops ...func(context.Context) error) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	var err error
	select {
	case <-ctx.Done():
	case err = <-serveErr:
	}
	cancel() // 通知を解除し、停止処理中に 2 回目のシグナルが来たら既定どおり即終了させる

	for _, stop := range stops {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		err = errors.Join(err, stop(ctx))
		cancel()
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
