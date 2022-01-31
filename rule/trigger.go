package rule

import (
	"errors"
	"regexp"

	"github.com/tidwall/gjson"
	"gopkg.in/Knetic/govaluate.v2"
)

type RuleTrigger struct {
	// RuleTrigger is a string that contains an evaluable rule.  For example,
	// `[event.repository.id] == "xyz" && [repositoryId] == "abc"`
	RuleExpression string `json:"rule"`
	accessors      []string
}

func (c *RuleTrigger) Evaluate(payload map[string]interface{}) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(c.RuleExpression)
	if err != nil {
		return false, err
	}
	result, err := expression.Evaluate(payload)
	if err != nil {
		return false, err
	}
	switch v := result.(type) {
	case bool:
		return v, nil
	}
	return false, errors.New("rule could not be evaluated to boolean")
}

func (c *RuleTrigger) extractAccessors() []string {
	r := regexp.MustCompile(`\[([\w\d-.]*)\]`)
	matches := r.FindAllStringSubmatch(c.RuleExpression, -1)
	accessors := []string{}
	for _, v := range matches {
		accessors = append(accessors, v[1])
	}
	return accessors
}

func (c *RuleTrigger) MakeParams(eventJSON []byte) map[string]interface{} {
	params := map[string]interface{}{}
	for _, v := range c.accessors {
		params[v] = gjson.GetBytes(eventJSON, v).Value()
	}
	return params
}
