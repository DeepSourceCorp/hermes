package redis

import (
	"context"
	"encoding/json"

	model "github.com/deepsourcelabs/hermes/rule"
	"github.com/go-redis/redis/v8"
)

type ruleStore struct {
	Conn *redis.Client
}

type action struct {
	Name       string
	TemplateID string
}

type ruleObj struct {
	Trigger model.Trigger `json:"trigger"`
	Action  action        `json:"action"`
}

func NewRuleStore(conn *redis.Client) model.Repository {
	return &ruleStore{
		Conn: conn,
	}
}

func (store *ruleStore) Create(ctx context.Context, rule *model.Rule) error {
	a := action{
		Name:       rule.Action.Name(),
		TemplateID: rule.Action.TemplateID(),
	}

	ro := ruleObj{
		Trigger: rule.Trigger,
		Action:  a,
	}

	raw, err := json.Marshal(ro)
	if err != nil {
		return err
	}
	if err := store.Conn.Set(ctx, rule.Trigger.RuleExpression, raw, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (store *ruleStore) GetByID(ctx context.Context, subscriberID, subscriptionID, id string) (*model.Rule, error) {
	r := new(model.Rule)
	res, err := store.Conn.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(res), r); err != nil {
		return nil, err
	}
	return r, nil
}
