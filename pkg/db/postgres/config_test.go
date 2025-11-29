package postgres_test

import (
	"testing"
	"time"

	"github.com/a-h-pooladvand/microservices/pkg/db/postgres"
)

func TestConfig_DSN(t *testing.T) {
	cfg := postgres.NewConfig(
		"localhost",
		"5432",
		"testuser",
		"testpass",
		"testdb",
		10,
	)

	dsn := cfg.DSN()

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable TimeZone=Asia/Tehran"
	if dsn != expected {
		t.Errorf("DSN() = %v, want %v", dsn, expected)
	}
}

func TestConfig_DSN_WithOptions(t *testing.T) {
	cfg := postgres.NewConfig(
		"localhost",
		"5432",
		"testuser",
		"testpass",
		"testdb",
		10,
		postgres.WithSSLMode("require"),
		postgres.WithTimeZone("UTC"),
		postgres.WithMaxOpenConns(50),
		postgres.WithMaxIdleConns(20),
		postgres.WithConnMaxLifetime(10*time.Minute),
	)

	if cfg.SSLMode != "require" {
		t.Errorf("SSLMode = %v, want %v", cfg.SSLMode, "require")
	}

	if cfg.TimeZone != "UTC" {
		t.Errorf("TimeZone = %v, want %v", cfg.TimeZone, "UTC")
	}

	if cfg.MaxOpenConns != 50 {
		t.Errorf("MaxOpenConns = %v, want %v", cfg.MaxOpenConns, 50)
	}

	if cfg.MaxIdleConns != 20 {
		t.Errorf("MaxIdleConns = %v, want %v", cfg.MaxIdleConns, 20)
	}

	if cfg.ConnMaxLifetime != 10*time.Minute {
		t.Errorf("ConnMaxLifetime = %v, want %v", cfg.ConnMaxLifetime, 10*time.Minute)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     postgres.Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: postgres.Config{
				Host:     "localhost",
				Port:     "5432",
				Username: "user",
				DB:       "testdb",
			},
			wantErr: false,
		},
		{
			name: "missing host",
			cfg: postgres.Config{
				Port:     "5432",
				Username: "user",
				DB:       "testdb",
			},
			wantErr: true,
		},
		{
			name: "missing port",
			cfg: postgres.Config{
				Host:     "localhost",
				Username: "user",
				DB:       "testdb",
			},
			wantErr: true,
		},
		{
			name: "missing username",
			cfg: postgres.Config{
				Host: "localhost",
				Port: "5432",
				DB:   "testdb",
			},
			wantErr: true,
		},
		{
			name: "missing db",
			cfg: postgres.Config{
				Host:     "localhost",
				Port:     "5432",
				Username: "user",
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
	cfg := postgres.NewConfig(
		"localhost",
		"5432",
		"testuser",
		"testpass",
		"testdb",
		10,
	)

	if cfg.MaxOpenConns != 25 {
		t.Errorf("default MaxOpenConns = %v, want %v", cfg.MaxOpenConns, 25)
	}

	if cfg.MaxIdleConns != 10 {
		t.Errorf("default MaxIdleConns = %v, want %v", cfg.MaxIdleConns, 10)
	}

	if cfg.ConnMaxLifetime != 5*time.Minute {
		t.Errorf("default ConnMaxLifetime = %v, want %v", cfg.ConnMaxLifetime, 5*time.Minute)
	}

	if cfg.SSLMode != "disable" {
		t.Errorf("default SSLMode = %v, want %v", cfg.SSLMode, "disable")
	}

	if cfg.TimeZone != "Asia/Tehran" {
		t.Errorf("default TimeZone = %v, want %v", cfg.TimeZone, "Asia/Tehran")
	}
}
