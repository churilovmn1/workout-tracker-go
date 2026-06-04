// Package broker implements a Redis-backed message queue for asynchronous
// notifications (Telegram reminders, trainer comments). Events are pushed onto a
// Redis list and drained by a Worker running in its own goroutine.
package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// queueKey is the Redis list that backs the event queue.
const queueKey = "workout:events"

// Publisher publishes events onto the broker queue.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

// RedisBroker is a Redis-backed Publisher + Subscriber. It uses a Redis list as
// a simple durable FIFO queue: Publish does LPUSH, Consume does blocking BRPOP.
type RedisBroker struct {
	client *redis.Client
}

// NewRedisBroker wraps an existing Redis client.
func NewRedisBroker(client *redis.Client) *RedisBroker {
	return &RedisBroker{client: client}
}

// Connect parses a redis:// URL, opens a client and pings it.
func Connect(ctx context.Context, url string) (*RedisBroker, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	client := redis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return NewRedisBroker(client), nil
}

// Publish marshals and enqueues an event.
func (b *RedisBroker) Publish(ctx context.Context, event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	if err := b.client.LPush(ctx, queueKey, data).Err(); err != nil {
		return fmt.Errorf("enqueue event: %w", err)
	}
	return nil
}

// Consume blocks up to timeout for the next event. It returns (nil, nil) when
// the wait times out with no event, so callers can poll in a loop.
func (b *RedisBroker) Consume(ctx context.Context, timeout time.Duration) (*Event, error) {
	res, err := b.client.BRPop(ctx, timeout, queueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // timed out, no event
		}
		return nil, fmt.Errorf("dequeue event: %w", err)
	}
	// res == [queueKey, value]
	var event Event
	if err := json.Unmarshal([]byte(res[1]), &event); err != nil {
		return nil, fmt.Errorf("unmarshal event: %w", err)
	}
	return &event, nil
}

// Close releases the underlying Redis connection.
func (b *RedisBroker) Close() error {
	return b.client.Close()
}

// NoopPublisher is used when Redis is not configured: publishing is a no-op so
// callers never have to nil-check the publisher.
type NoopPublisher struct{}

// Publish discards the event.
func (NoopPublisher) Publish(context.Context, Event) error { return nil }
