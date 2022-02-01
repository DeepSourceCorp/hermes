package rule

type Rule struct {
	Trigger Trigger `json:"trigger"`
	Action  Action  `json:"action"`
}
