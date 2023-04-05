package config

import (
	"strings"
	"time"
)

type NatsConfig struct {
	// NatsHost     string `env:"NATS_HOST" envDefault:"nats"`
	// NatsPort     uint16 `env:"NATS_PORT" envDefault:"4222"`
	NatsAddresses string `envconfig:"NATS_ADDRESSES" required:"true"`
	NatsUser      string `envconfig:"NATS_USER" default:"nast"`
	NatsPassword  string `envconfig:"NATS_PASSWORD" default:"password"`

	NatsConnectionRetryOnFailed bool          `envconfig:"NATS_CONNECTION_RETRY" default:"true"`
	NatsConnectionRetryCount    uint16        `envconfig:"NATS_CONNECTION_RETRY_COUNT" default:"30"`
	NatsConnectionRetryTimeout  time.Duration `envconfig:"NATS_CONNECTION_RETRY_TIMEOUT" default:"15s"`

	NatsFlushTimeOut time.Duration `envconfig:"NATS_FLUSH_TIMEOUT" default:"15s"`

	NatsWorkersPerConsumer uint16 `envconfig:"NATS_WORKER_PER_CONSUMER" default:"5"`

	NatsSubscriptionRetry        bool          `envconfig:"NATS_WORKER_SUBSCRIPTION_RETRY" default:"true"`
	NatsSubscriptionRetryCount   uint16        `envconfig:"NATS_WORKER_SUBSCRIPTION_RETRY_COUNT" default:"3"`
	NatsSubscriptionRetryTimeout time.Duration `envconfig:"NATS_WORKER_SUBSCRIPTION_RETRY_TIMEOUT" default:"3s"`

	// calculated variables
	nastAddresses []string
}

func (c *NatsConfig) GetNatsAddresses() []string {
	return c.nastAddresses
}

func (c *NatsConfig) GetNatsJoinedAddresses() string {
	return c.NatsAddresses
}

// func (c *NatsConfig) GetNatsHost() string {
//	return c.NatsHost
// }
//
// func (c *NatsConfig) GetNatsPort() uint16 {
//	return c.NatsPort
// }

func (c *NatsConfig) GetNatsUser() string {
	return c.NatsUser
}

func (c *NatsConfig) GetNatsPassword() string {
	return c.NatsPassword
}

func (c *NatsConfig) IsRetryOnConnectionFailed() bool {
	return c.NatsConnectionRetryOnFailed
}

func (c *NatsConfig) GetNatsConnectionRetryCount() uint16 {
	return c.NatsConnectionRetryCount
}

func (c *NatsConfig) GetNatsConnectionRetryTimeout() time.Duration {
	return c.NatsConnectionRetryTimeout
}

func (c *NatsConfig) GetFlushTimeout() time.Duration {
	return c.NatsFlushTimeOut
}

func (c *NatsConfig) GetWorkersCountPerConsumer() uint16 {
	return c.NatsWorkersPerConsumer
}

// Prepare variables to static configuration
func (c *NatsConfig) Prepare() error {
	endpoints := strings.Split(c.NatsAddresses, ",")
	length := len(endpoints)
	if length < 1 {
		return nil
	}
	c.nastAddresses = endpoints

	return nil
}
