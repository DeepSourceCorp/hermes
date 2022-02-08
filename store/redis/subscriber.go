package redis

import (
	"context"
	"encoding/json"
	"fmt"

	model "github.com/deepsourcelabs/hermes/subscriber"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/ksuid"
)

type subscriberStore struct {
	Conn *redis.Client
}

func NewSubscriberStore(conn *redis.Client) model.Repository {
	return &subscriberStore{
		Conn: conn,
	}
}

func (store *subscriberStore) Create(ctx context.Context, subscriber *model.Subscriber) (*model.Subscriber, error) {
	subscriber.ID = ksuid.New().String()
	raw, err := json.Marshal(subscriber)
	if err != nil {
		return nil, err
	}

	if err := store.Conn.Set(
		ctx,
		fmt.Sprintf("subscriber:%s", subscriber.ID),
		raw,
		0,
	).Err(); err != nil {
		return nil, err
	}
	return subscriber, nil
}

func (store *subscriberStore) GetByID(ctx context.Context, id string) (*model.Subscriber, error) {
	subscriber := new(model.Subscriber)

	res, err := store.Conn.Get(ctx,
		fmt.Sprintf("subscriber:%s", id),
	).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(res), subscriber); err != nil {
		return nil, err
	}
	return subscriber, nil
}
