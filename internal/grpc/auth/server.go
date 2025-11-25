package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ssov1 "github.com/diyor200/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(ctx context.Context, email, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error)
}

func Register(grpcServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(grpcServer, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Login(ctx context.Context, in *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "appId is required")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.EInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *ServerAPI) Register(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	uid, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &ssov1.RegisterResponse{UserId: uid}, nil
}
