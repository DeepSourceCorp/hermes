package redis

import (
	"context"
	"encoding/json"

	"github.com/deepsourcelabs/hermes/subscriber"
	model "github.com/deepsourcelabs/hermes/subscriber"
	"github.com/go-redis/redis/v8"
)

type subscriberStore struct {
	Conn *redis.Client
}

func NewSubscriberStore(conn *redis.Client) subscriber.Repository {
	return &subscriberStore{
		Conn: conn,
	}
}

func (store *subscriberStore) Create(ctx context.Context, subscriber *model.Subscriber) error {
	raw, err := json.Marshal(subscriber)
	if err != nil {
		return err
	}
	if err := store.Conn.Set(ctx, subscriber.ID, raw, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (store *subscriberStore) GetByID(ctx context.Context, id string) (*model.Subscriber, error) {
	subscriber := new(model.Subscriber)
	res, err := store.Conn.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(res), subscriber); err != nil {
		return nil, err
	}
	return subscriber, nil
}
