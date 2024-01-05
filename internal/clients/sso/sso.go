package sso

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ssov1 "github.com/Shuv1Wolf/jwt_protos/gen/go/sso"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	apiAuth ssov1.AuthClient
	apiPing ssov1.PingClient
	log     *slog.Logger
}

func New(ctx context.Context,
	log *slog.Logger,
	address string,
	timeout time.Duration,
	retriesCount int,
	appId int64,
	appName string,
	appSecret string,
) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	authClient := ssov1.NewAuthClient(cc)
	pingClient := ssov1.NewPingClient(cc)

	client := &Client{
		apiAuth: authClient,
		apiPing: pingClient,
	}

	log.Info("checking connection to grps")
	ping, err := client.ping(int(appId))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !ping {
		log.Info("the application is not in the sso")
		log.Info("creating a new application with information from the config")
		res, err := client.apiPing.NewApp(context.Background(), &ssov1.IsNewAppRequest{
			Id:     int64(appId),
			Name:   appName,
			Secret: appSecret,
		})
		if err != nil {
			panic(err)
		}
		log.Info(res.String())
	}

	return &Client{
		apiAuth: authClient,
		apiPing: pingClient,
	}, nil
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) ping(appId int) (bool, error) {
	ping, err := c.apiPing.Ping(context.Background(), &ssov1.IsPingRequest{
		AppId: int64(appId),
	})
	if err != nil {
		return ping.GetClient(), err
	}
	return ping.GetClient(), nil
}

func (c *Client) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "grpc.IsAdmin"

	resp, err := c.apiAuth.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: userID,
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.IsAdmin, nil
}
