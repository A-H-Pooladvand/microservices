package rabbitmq

import "errors"

// RabbitMQ errors.
var (
	// ErrInvalidAddress is returned when the address is invalid.
	ErrInvalidAddress = errors.New("rabbitmq: invalid address")
	// ErrConnectionFailed is returned when the connection fails.
	ErrConnectionFailed = errors.New("rabbitmq: connection failed")
	// ErrChannelClosed is returned when the channel is closed.
	ErrChannelClosed = errors.New("rabbitmq: channel closed")
	// ErrPublishFailed is returned when publishing fails.
	ErrPublishFailed = errors.New("rabbitmq: publish failed")
	// ErrConsumeFailed is returned when consuming fails.
	ErrConsumeFailed = errors.New("rabbitmq: consume failed")
	// ErrQueueDeclareFailed is returned when queue declaration fails.
	ErrQueueDeclareFailed = errors.New("rabbitmq: queue declaration failed")
	// ErrExchangeDeclareFailed is returned when exchange declaration fails.
	ErrExchangeDeclareFailed = errors.New("rabbitmq: exchange declaration failed")
)
