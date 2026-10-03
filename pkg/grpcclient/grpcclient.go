// Package grpcclient はサービス間の gRPC 接続を作る
package grpcclient

import (
	"fmt"
	"os"

	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/telemetry"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/userid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Dial は addrEnv の接続先に、trace context と x-user-id を引き継ぐ client で接続する。
// 平文なのは、クラスタ内の暗号化を Istio の mTLS に任せるため
func Dial(addrEnv string) (*grpc.ClientConn, error) {
	addr := os.Getenv(addrEnv)
	if addr == "" {
		return nil, fmt.Errorf("%s is required", addrEnv)
	}
	return grpc.NewClient(addr, append(telemetry.ClientOptions(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(userid.Forward()))...)
}
