package ping

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ssov1 "github.com/Shuv1Wolf/jwt_protos/gen/go/sso"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api ssov1.PingClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
	appId int64,
	appName string,
	appSecret string,
) (*Client, error) {
	const op = "ping.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	// TODO: сделать защищённое соединение
	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := &Client{
		api: ssov1.NewPingClient(cc),
	}

	log.Info("checking connection to grps")
	ping, err := client.ping(int(appId))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !ping {
		log.Info("the application is not in the sso")
		log.Info("creating a new application with information from the config")
		res, err := client.api.NewApp(context.Background(), &ssov1.IsNewAppRequest{
			Id:     int64(appId),
			Name:   appName,
			Secret: appSecret,
		})
		if err != nil {
			panic(err)
		}
		log.Info(res.String())
	}

	return client, err
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func (c *Client) ping(appId int) (bool, error) {
	ping, err := c.api.Ping(context.Background(), &ssov1.IsPingRequest{
		AppId: int64(appId),
	})
	if err != nil {
		return ping.GetClient(), err
	}
	return ping.GetClient(), nil
}
