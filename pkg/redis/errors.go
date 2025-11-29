package redis

import "errors"

// Redis errors.
var (
	// ErrInvalidAddress is returned when the address is invalid or empty.
	ErrInvalidAddress = errors.New("redis: address is required")
	// ErrConnectionFailed is returned when the connection fails.
	ErrConnectionFailed = errors.New("redis: connection failed")
	// ErrKeyNotFound is returned when a key is not found.
	ErrKeyNotFound = errors.New("redis: key not found")
	// ErrLockNotObtained is returned when a lock cannot be obtained.
	ErrLockNotObtained = errors.New("redis: lock not obtained")
	// ErrOperationFailed is returned when an operation fails.
	ErrOperationFailed = errors.New("redis: operation failed")
)
