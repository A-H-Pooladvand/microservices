package configs

import (
	"crypto/tls"
	"crypto/x509"
	"go.uber.org/zap"
	"os"
	"po/pkg/vault"
	"time"
)

type Etcd struct {
	Host        string `env:"ETCD_HOST" json:"host"`
	Port        string `env:"ETCD_PORT" json:"port"`
	Username    string `env:"ETCD_USERNAME" json:"username"`
	Password    string `env:"ETCD_PASSWORD" json:"password"`
	Timeout     int    `env:"ETCD_TIMEOUT" json:"timeout"`
	TLSCertFile string `env:"ETCD_TLS_CERT_FILE" json:"tls_cert_file"`
	TLSKeyFile  string `env:"ETCD_TLS_KEY_FILE" json:"tls_key_file"`
}

func (e Etcd) Endpoints() []string {
	return []string{
		e.Host + ":" + e.Port,
	}
}

func (e Etcd) DialTimeout() time.Duration {
	return time.Duration(e.Timeout) * time.Second
}

func (e Etcd) AutoSyncInterval() time.Duration {
	return time.Minute
}

func (e Etcd) DialKeepAliveTime() time.Duration {
	return time.Minute
}

func (e Etcd) DialKeepAliveTimeout() time.Duration {
	return time.Second * 20
}

func (e Etcd) TLS() *tls.Config {
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client-key.pem")

	if err != nil {
		zap.L().Fatal("failed to load cert key pair", zap.Error(err))
	}

	caCert, err := os.ReadFile("certs/ca.pem")
	if err != nil {
		zap.L().Fatal("failed to load cert file caCert", zap.Error(err))
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		zap.L().Fatal("failed to append CA Cert")
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}
}

func NewEtcd(client *vault.Client) (*Etcd, error) {
	c := &Etcd{}

	err := Parse("etcd", c, client)

	return c, err
}
