package client

import (
	"context"
	"log/slog"
	"os"
	"time"
	"url-shortener/internal/clients/sso/grpc/ping"
	pinggrpc "url-shortener/internal/clients/sso/grpc/ping"
	"url-shortener/internal/lib/logger/sl"
)

type GRPC struct {
	ping *ping.Client
}

func New(ctx context.Context,
	log *slog.Logger,
	address string,
	timeout time.Duration,
	retriesCount int,
	appId int64,
	appName string,
	appSecret string,
) *GRPC {
	const op = "client.New"

	ping, err := pinggrpc.New(
		context.Background(), log,
		address, timeout, retriesCount,
		appId, appName, appSecret,
	)
	if err != nil {
		log.Error(op, "failed to init ping client", sl.Err(err))
		os.Exit(1)
	}

	return &GRPC{
		ping: ping,
	}
}
