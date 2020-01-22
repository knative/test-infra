package sender

import (
	"knative.dev/test-infra/tools/wathola/internal/config"
	"math/rand"
	"time"
)

// New creates new Sender
func New() Sender {
	config.ReadIfPresent()
	return &sender{
		active:  true,
		counter: 0,
	}
}

// NewEventID creates new event ID
func NewEventID() string {
	return randString(16)
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func randStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randString(length int) string {
	return randStringWithCharset(length, charset)
}
