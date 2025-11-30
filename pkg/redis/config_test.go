package redis_test

import (
	"testing"
	"time"

	"github.com/a-h-pooladvand/microservices/pkg/redis"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     redis.Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: redis.Config{
				Address: "localhost:6379",
			},
			wantErr: false,
		},
		{
			name: "missing address",
			cfg: redis.Config{
				Address: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewConfig_Defaults(t *testing.T) {
	cfg := redis.NewConfig("localhost:6379", "user", "pass")

	if cfg.DB != 0 {
		t.Errorf("default DB = %v, want %v", cfg.DB, 0)
	}

	if cfg.MaxRetries != 3 {
		t.Errorf("default MaxRetries = %v, want %v", cfg.MaxRetries, 3)
	}

	if cfg.PoolSize != 10 {
		t.Errorf("default PoolSize = %v, want %v", cfg.PoolSize, 10)
	}

	if cfg.MinIdleConns != 3 {
		t.Errorf("default MinIdleConns = %v, want %v", cfg.MinIdleConns, 3)
	}

	if cfg.DialTimeout != 5*time.Second {
		t.Errorf("default DialTimeout = %v, want %v", cfg.DialTimeout, 5*time.Second)
	}

	if cfg.ReadTimeout != 3*time.Second {
		t.Errorf("default ReadTimeout = %v, want %v", cfg.ReadTimeout, 3*time.Second)
	}

	if cfg.WriteTimeout != 3*time.Second {
		t.Errorf("default WriteTimeout = %v, want %v", cfg.WriteTimeout, 3*time.Second)
	}
}

func TestNewConfig_WithOptions(t *testing.T) {
	cfg := redis.NewConfig(
		"localhost:6379",
		"user",
		"pass",
		redis.WithDB(1),
		redis.WithMaxRetries(5),
		redis.WithPoolSize(20),
		redis.WithMinIdleConns(5),
		redis.WithDialTimeout(10*time.Second),
		redis.WithReadTimeout(5*time.Second),
		redis.WithWriteTimeout(5*time.Second),
	)

	if cfg.DB != 1 {
		t.Errorf("DB = %v, want %v", cfg.DB, 1)
	}

	if cfg.MaxRetries != 5 {
		t.Errorf("MaxRetries = %v, want %v", cfg.MaxRetries, 5)
	}

	if cfg.PoolSize != 20 {
		t.Errorf("PoolSize = %v, want %v", cfg.PoolSize, 20)
	}

	if cfg.MinIdleConns != 5 {
		t.Errorf("MinIdleConns = %v, want %v", cfg.MinIdleConns, 5)
	}

	if cfg.DialTimeout != 10*time.Second {
		t.Errorf("DialTimeout = %v, want %v", cfg.DialTimeout, 10*time.Second)
	}

	if cfg.ReadTimeout != 5*time.Second {
		t.Errorf("ReadTimeout = %v, want %v", cfg.ReadTimeout, 5*time.Second)
	}

	if cfg.WriteTimeout != 5*time.Second {
		t.Errorf("WriteTimeout = %v, want %v", cfg.WriteTimeout, 5*time.Second)
	}
}

func TestConfig_Address(t *testing.T) {
	cfg := redis.NewConfig("redis.example.com:6379", "admin", "secret")

	if cfg.Address != "redis.example.com:6379" {
		t.Errorf("Address = %v, want %v", cfg.Address, "redis.example.com:6379")
	}

	if cfg.User != "admin" {
		t.Errorf("User = %v, want %v", cfg.User, "admin")
	}

	if cfg.Password != "secret" {
		t.Errorf("Password = %v, want %v", cfg.Password, "secret")
	}
}
