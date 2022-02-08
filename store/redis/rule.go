package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	model "github.com/deepsourcelabs/hermes/rule"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/ksuid"
)

type storeSubscription struct {
	Conn *redis.Client
}

func NewRuleStore(conn *redis.Client) model.Repository {
	return &storeSubscription{
		Conn: conn,
	}
}

func (s *storeSubscription) Create(ctx context.Context, rule *model.SerializableRule) (*model.SerializableRule, error) {
	rule.ID = ksuid.New().String()
	raw, err := json.Marshal(rule)
	if err != nil {
		return nil, err
	}

	pipe := s.Conn.TxPipeline()
	// Set the rule object.
	pipe.Expire(ctx, "tx_pipeline_counter", 1*time.Hour)
	if err := pipe.Set(
		ctx,
		fmt.Sprintf("rule:%s", rule.ID),
		raw,
		0,
	).Err(); err != nil {
		return nil, err
	}

	// Set the ruleID against the subcsriptionID.
	if err := pipe.LPush(
		ctx,
		fmt.Sprintf("rule-list:%s", rule.SubscriptionID),
		rule.ID,
	).Err(); err != nil {
		return nil, err
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *storeSubscription) GetByID(ctx context.Context, id string) (*model.SerializableRule, error) {
	rule := new(model.SerializableRule)
	res, err := s.Conn.Get(
		ctx,
		fmt.Sprintf("rule:%s", id),
	).Result()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(res), rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *storeSubscription) GetAll(ctx context.Context, subscriptionID string) ([]model.SerializableRule, error) {
	rules := []model.SerializableRule{}

	key := fmt.Sprintf("rule-list:%s", subscriptionID)

	// Count of subscriptions stored against the subscriberID
	size, err := s.Conn.LLen(
		ctx,
		key,
	).Result()
	if err != nil {
		return nil, err
	}

	// Get the subscription IDs for the subscription.
	ids, err := s.Conn.LRange(
		ctx,
		key,
		0,
		size,
	).Result()
	if err != nil {
		return rules, err
	}

	// Populate a slice with the rules for the subscription.
	for _, id := range ids {
		res, err := s.Conn.Get(
			ctx,
			fmt.Sprintf("rule:%s", id),
		).Result()
		if err != nil {
			return rules, err
		}

		s := new(model.SerializableRule)
		if err := json.Unmarshal([]byte(res), s); err != nil {
			return rules, err
		}
		rules = append(rules, *s)
	}

	return rules, nil
}
