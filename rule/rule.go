package rule

type Rule struct {
	Trigger RuleTrigger `json:"trigger"`
	Action  RuleAction  `json:"action"`
}
