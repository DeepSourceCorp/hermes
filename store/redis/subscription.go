package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deepsourcelabs/hermes/subscription"
	model "github.com/deepsourcelabs/hermes/subscription"
	"github.com/go-redis/redis/v8"
)

type subscriptionStore struct {
	Conn *redis.Client
}

func NewSubscriptionStore(Conn *redis.Client) subscription.Repository {
	return &subscriptionStore{
		Conn,
	}
}

func (store *subscriptionStore) Create(ctx context.Context, subscription *model.Subscription) (*model.Subscription, error) {
	raw, err := json.Marshal(subscription)
	if err != nil {
		return nil, err
	}

	fmt.Println("thenga")
	fmt.Println(subscription)

	pipe := store.Conn.TxPipeline()
	pipe.Expire(ctx, "tx_pipeline_counter", 1*time.Hour)
	key1 := fmt.Sprintf("subscription:%s:%s", subscription.SubscriberID, subscription.ID)
	if err := pipe.Set(ctx, key1, raw, 0).Err(); err != nil {
		return nil, err
	}

	key2 := fmt.Sprintf("subscription-list:%s", subscription.SubscriberID)
	if err := pipe.LPush(ctx, key2, subscription.ID).Err(); err != nil {
		return nil, err
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (store *subscriptionStore) GetByID(ctx context.Context, subscriberID, id string) (*model.Subscription, error) {
	key := fmt.Sprintf("subscription:%s:%s", subscriberID, id)

	subscription := new(model.Subscription)
	res, err := store.Conn.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(res), subscription); err != nil {
		return nil, err
	}
	return subscription, nil
}

func (store *subscriptionStore) GetAll(ctx context.Context, subscriberID string) ([]model.Subscription, error) {
	subscriptions := []model.Subscription{}

	key1 := fmt.Sprintf("subscription-list:%s", subscriberID)
	sLen, err := store.Conn.LLen(ctx, key1).Result()
	if err != nil {
		return nil, err
	}

	sIDs, err := store.Conn.LRange(ctx, key1, 0, sLen).Result()
	if err != nil {
		return subscriptions, err
	}

	for _, id := range sIDs {
		key2 := fmt.Sprintf("subscription:%s:%s", subscriberID, id)
		res, err := store.Conn.Get(ctx, key2).Result()
		if err != nil {
			return subscriptions, err
		}
		sub := new(model.Subscription)
		if err := json.Unmarshal([]byte(res), sub); err != nil {
			return subscriptions, err
		}
		fmt.Println(sub)
		subscriptions = append(subscriptions, *sub)
	}

	return subscriptions, nil
}
