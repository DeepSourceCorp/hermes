package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deepsourcelabs/hermes/subscription"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/ksuid"
)

type subscriptionStore struct {
	Conn *redis.Client
}

func NewSubscriptionStore(Conn *redis.Client) subscription.Repository {
	return &subscriptionStore{
		Conn,
	}
}

func (store *subscriptionStore) Create(ctx context.Context, subscription *subscription.Subscription) (*subscription.Subscription, error) {
	subscription.ID = ksuid.New().String()
	raw, err := json.Marshal(subscription)
	if err != nil {
		return nil, err
	}

	// Start a Redis MULTI
	pipe := store.Conn.TxPipeline()
	pipe.Expire(ctx, "tx_pipeline_counter", 1*time.Hour)

	// Set the subscription object.
	if err := pipe.Set(
		ctx,
		fmt.Sprintf("subscription:%s", subscription.ID),
		raw,
		0,
	).Err(); err != nil {
		return nil, err
	}

	// Map the subscription against the subscriber.
	if err := pipe.LPush(
		ctx,
		fmt.Sprintf("subscription-list:%s", subscription.SubscriberID),
		subscription.ID,
	).Err(); err != nil {
		return nil, err
	}

	// Redis EXEC
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (store *subscriptionStore) GetByID(ctx context.Context, id string) (*subscription.Subscription, error) {
	subscription := new(subscription.Subscription)

	res, err := store.Conn.Get(
		ctx,
		fmt.Sprintf("subscription:%s", id),
	).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(res), subscription); err != nil {
		return nil, err
	}
	return subscription, nil
}

func (store *subscriptionStore) GetAll(ctx context.Context, subscriberID string) ([]subscription.Subscription, error) {
	subscriptions := []subscription.Subscription{}

	key := fmt.Sprintf("subscription-list:%s", subscriberID)

	// Count of subscriptions stored against the subscriberID
	size, err := store.Conn.LLen(
		ctx,
		key,
	).Result()
	if err != nil {
		return nil, err
	}

	// Get the subscription IDs for the subscriber
	ids, err := store.Conn.LRange(
		ctx,
		key,
		0,
		size,
	).Result()
	if err != nil {
		return subscriptions, err
	}

	// Populate a slice with the subscriptions for the subscriber
	for _, id := range ids {
		res, err := store.Conn.Get(
			ctx,
			fmt.Sprintf("subscription:%s", id),
		).Result()
		if err != nil {
			return subscriptions, err
		}

		s := new(subscription.Subscription)
		if err := json.Unmarshal([]byte(res), s); err != nil {
			return subscriptions, err
		}
		subscriptions = append(subscriptions, *s)
	}

	return subscriptions, nil
}
