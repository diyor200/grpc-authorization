package auth

import (
	"context"
	"grpc-authorize/internal/domain/models"
	"log/slog"
	"time"
)

type UserStorage interface {
	Save(ctx context.Context, email string, passHash []byte) (int64, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
}

type AppProvider interface {
	Get(ctx context.Context, id int) (models.App, error)
}

type Auth struct {
	log      *slog.Logger
	userRepo UserStorage
	appRepo  AppProvider
	tokenTTL time.Duration
}

func New(l *slog.Logger, userRepo UserStorage, appRepo AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      l,
		userRepo: userRepo,
		appRepo:  appRepo,
		tokenTTL: tokenTTL,
	}
}
