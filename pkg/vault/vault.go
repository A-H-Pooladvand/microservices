package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/a-h-pooladvand/microservices/pkg/observability"
	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

// Client wraps Vault client with observability support.
type Client struct {
	config  Config
	conn    *vault.Client
	tracer  *observability.Tracer
	metrics *OperationMetrics
	logger  *zap.Logger
}

// OperationMetrics holds metrics for Vault operations.
type OperationMetrics struct {
	RequestsTotal   metric.Int64Counter
	RequestDuration metric.Float64Histogram
	ErrorsTotal     metric.Int64Counter
}

// New creates a new Vault client with observability.
func New(cfg Config, logger *zap.Logger) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	c := vault.DefaultConfig()
	c.Address = cfg.Address
	c.Timeout = cfg.Timeout
	c.MaxRetries = cfg.MaxRetries
	c.MinRetryWait = cfg.RetryWaitMin
	c.MaxRetryWait = cfg.RetryWaitMax

	if cfg.TLSEnabled {
		tlsConfig := &vault.TLSConfig{
			Insecure: cfg.TLSSkipVerify,
		}
		if cfg.CACert != "" {
			tlsConfig.CACert = cfg.CACert
		}
		if cfg.ClientCert != "" && cfg.ClientKey != "" {
			tlsConfig.ClientCert = cfg.ClientCert
			tlsConfig.ClientKey = cfg.ClientKey
		}
		if err := c.ConfigureTLS(tlsConfig); err != nil {
			return nil, fmt.Errorf("configure tls: %w", err)
		}
	}

	vc, err := vault.NewClient(c)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	if cfg.Namespace != "" {
		vc.SetNamespace(cfg.Namespace)
	}

	// Initialize observability
	tracer := observability.NewTracer("vault")
	meter := observability.NewMeter("vault")

	requestsTotal, err := meter.Counter("vault_requests_total", "Total number of Vault requests")
	if err != nil {
		return nil, fmt.Errorf("create requests counter: %w", err)
	}

	requestDuration, err := meter.Histogram("vault_request_duration_seconds", "Vault request duration in seconds", "s")
	if err != nil {
		return nil, fmt.Errorf("create request duration histogram: %w", err)
	}

	errorsTotal, err := meter.Counter("vault_errors_total", "Total number of Vault errors")
	if err != nil {
		return nil, fmt.Errorf("create errors counter: %w", err)
	}

	if logger == nil {
		logger = zap.NewNop()
	}

	client := &Client{
		config: cfg,
		conn:   vc,
		tracer: tracer,
		metrics: &OperationMetrics{
			RequestsTotal:   requestsTotal,
			RequestDuration: requestDuration,
			ErrorsTotal:     errorsTotal,
		},
		logger: logger,
	}

	// Login if credentials are provided
	if cfg.RoleID != "" && cfg.SecretID != "" {
		if _, err := client.login(context.Background()); err != nil {
			return nil, fmt.Errorf("vault login failed: %w", err)
		}
	}

	return client, nil
}

func (c *Client) login(ctx context.Context) (*vault.Secret, error) {
	ctx, span := c.tracer.StartWithAttributes(ctx, "vault.Login",
		attribute.String("vault.method", "approle"),
	)
	defer span.End()

	startTime := time.Now()

	secretID := &auth.SecretID{
		FromString: c.config.SecretID,
	}

	appRoleAuth, err := auth.NewAppRoleAuth(c.config.RoleID, secretID)
	if err != nil {
		observability.SetSpanError(span, err)
		c.recordMetrics(ctx, "login", startTime, err)
		return nil, err
	}

	secret, err := c.conn.Auth().Login(ctx, appRoleAuth)
	if err != nil {
		observability.SetSpanError(span, err)
		c.recordMetrics(ctx, "login", startTime, err)
		return nil, err
	}

	c.recordMetrics(ctx, "login", startTime, nil)
	observability.SetSpanOK(span)
	return secret, nil
}

// Get retrieves a secret from Vault with tracing.
func (c *Client) Get(ctx context.Context, path string) (*vault.KVSecret, error) {
	ctx, span := c.tracer.StartWithAttributes(ctx, "vault.Get",
		attribute.String("vault.path", path),
		attribute.String("vault.mount", c.config.MountPath),
	)
	defer span.End()

	startTime := time.Now()

	secret, err := c.conn.KVv2(c.config.MountPath).Get(ctx, path)
	if err != nil {
		observability.SetSpanError(span, err)
		c.recordMetrics(ctx, "get", startTime, err)
		c.logger.Error("failed to get secret",
			zap.String("path", path),
			zap.Error(err),
		)
		return nil, err
	}

	c.recordMetrics(ctx, "get", startTime, nil)
	observability.SetSpanOK(span)
	return secret, nil
}

// Put stores a secret in Vault with tracing.
func (c *Client) Put(ctx context.Context, path string, data map[string]interface{}) (*vault.KVSecret, error) {
	ctx, span := c.tracer.StartWithAttributes(ctx, "vault.Put",
		attribute.String("vault.path", path),
		attribute.String("vault.mount", c.config.MountPath),
	)
	defer span.End()

	startTime := time.Now()

	secret, err := c.conn.KVv2(c.config.MountPath).Put(ctx, path, data)
	if err != nil {
		observability.SetSpanError(span, err)
		c.recordMetrics(ctx, "put", startTime, err)
		c.logger.Error("failed to put secret",
			zap.String("path", path),
			zap.Error(err),
		)
		return nil, err
	}

	c.recordMetrics(ctx, "put", startTime, nil)
	observability.SetSpanOK(span)
	return secret, nil
}

// Delete deletes a secret from Vault with tracing.
func (c *Client) Delete(ctx context.Context, path string) error {
	ctx, span := c.tracer.StartWithAttributes(ctx, "vault.Delete",
		attribute.String("vault.path", path),
		attribute.String("vault.mount", c.config.MountPath),
	)
	defer span.End()

	startTime := time.Now()

	err := c.conn.KVv2(c.config.MountPath).Delete(ctx, path)
	if err != nil {
		observability.SetSpanError(span, err)
		c.recordMetrics(ctx, "delete", startTime, err)
		c.logger.Error("failed to delete secret",
			zap.String("path", path),
			zap.Error(err),
		)
		return err
	}

	c.recordMetrics(ctx, "delete", startTime, nil)
	observability.SetSpanOK(span)
	return nil
}

// Parse retrieves and parses a secret into a struct.
func (c *Client) Parse(ctx context.Context, path string, v any) error {
	secret, err := c.Get(ctx, path)
	if err != nil {
		return err
	}

	b, err := json.Marshal(secret.Data)
	if err != nil {
		return fmt.Errorf("marshal secret data: %w", err)
	}

	return json.Unmarshal(b, v)
}

// Health checks the health of the Vault server.
func (c *Client) Health(ctx context.Context) error {
	ctx, span := c.tracer.StartWithAttributes(ctx, "vault.Health")
	defer span.End()

	startTime := time.Now()

	health, err := c.conn.Sys().Health()
	if err != nil {
		observability.SetSpanError(span, err)
		c.recordMetrics(ctx, "health", startTime, err)
		return err
	}

	if !health.Initialized || health.Sealed {
		err := fmt.Errorf("vault is not healthy: initialized=%v sealed=%v", health.Initialized, health.Sealed)
		observability.SetSpanError(span, err)
		return err
	}

	c.recordMetrics(ctx, "health", startTime, nil)
	observability.SetSpanOK(span)
	return nil
}

func (c *Client) recordMetrics(ctx context.Context, operation string, startTime time.Time, err error) {
	duration := time.Since(startTime).Seconds()
	attrs := metric.WithAttributes(attribute.String("operation", operation))

	c.metrics.RequestsTotal.Add(ctx, 1, attrs)
	c.metrics.RequestDuration.Record(ctx, duration, attrs)

	if err != nil {
		c.metrics.ErrorsTotal.Add(ctx, 1, attrs)
	}
}
