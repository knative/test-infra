package config

import (
	"github.com/sirupsen/logrus"
	"time"
)

// ReceiverTeardownConfig holds config receiver teardown
type ReceiverTeardownConfig struct {
	Duration time.Duration
}

// ReceiverProgressConfig holds config receiver progress reporting
type ReceiverProgressConfig struct {
	Duration time.Duration
}

// ReceiverConfig hold configuration for receiver
type ReceiverConfig struct {
	Teardown ReceiverTeardownConfig
	Progress ReceiverProgressConfig
	Port     int
}

// SenderConfig hold configuration for sender
type SenderConfig struct {
	Address  string
	Interval time.Duration
	Cooldown time.Duration
}

// ForwarderConfig holds configuration for forwarder
type ForwarderConfig struct {
	Target string
	Port   int
}

// ReadinessConfig holds a readiness configuration
type ReadinessConfig struct {
	Enabled bool
	URI     string
	Message string
	Status  int
}

// Config hold complete configuration
type Config struct {
	Sender    SenderConfig
	Forwarder ForwarderConfig
	Receiver  ReceiverConfig
	Readiness ReadinessConfig
	LogLevel  logrus.Level
}
