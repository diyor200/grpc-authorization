package auth

import (
	"context"

	ssov1 "github.com/diyor200/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(ctx context.Context, email, password, appID string) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error)
}

func Register(grpcServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(grpcServer, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Login()
