package rule

import "context"

type RuleAction interface {
	Do(func(context.Context, interface{}) interface{}, error)
}
