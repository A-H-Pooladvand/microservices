package etcd

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"po/configs"
	"po/pkg/log"
	"time"
)

func Provide(lc fx.Lifecycle, config *configs.Etcd) (client *clientv3.Client, err error) {
	tls := config.TLS()

	client, err = clientv3.New(clientv3.Config{
		Endpoints:            config.Endpoints(),
		AutoSyncInterval:     config.AutoSyncInterval(),
		DialTimeout:          config.DialTimeout(),
		DialKeepAliveTime:    config.DialKeepAliveTime(),
		DialKeepAliveTimeout: config.DialKeepAliveTimeout(),
		MaxCallSendMsgSize:   4 << 10,
		MaxCallRecvMsgSize:   4 << 10,
		TLS:                  tls,
		//Username:              config.Username,
		//Password:              config.Password,
		RejectOldCluster:      true,
		Logger:                log.Logger(),
		MaxUnaryRetries:       5,
		BackoffWaitBetween:    time.Second,
		BackoffJitterFraction: 0.2,
		DialOptions: []grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(tls)),
		},
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return err
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, err
}
