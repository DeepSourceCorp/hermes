package queue

import (
	"context"
	"time"

	"github.com/deepsourcelabs/hermes/backoff"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn *amqp.Connection
}

type RabbitMQOpts struct {
	BrokerURL        string
	MinRetryDuration time.Duration
	MaxRetryDuration time.Duration
	MaxRetryAttempts int64
}

func (c *RabbitMQ) Dial(ctx context.Context, opts RabbitMQOpts) (err error) {
	b := backoff.NewFibonacci(false, opts.MinRetryDuration, opts.MaxRetryDuration)
	retryCount := 0
	for {
		if conn, err := amqp.Dial(opts.BrokerURL); err == nil {
			log.Info("connected to RabbitMQ")
			c.conn = conn
			return nil
		}
		log.Warnf("failed to connect to RabbitMQ, Error: %v", err)

		retryDuration := b.Duration()

		log.Warnf("retrying connection in %d", retryDuration)
		time.Sleep(retryDuration)
		retryCount++
	}
}

func (c *RabbitMQ) Close(ctx context.Context) error {
	return c.conn.Close()
}
