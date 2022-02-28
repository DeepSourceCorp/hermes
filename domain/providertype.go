package domain

type ProviderType string

const (
	ProviderTypeSlack ProviderType = "slack"
	ProviderTypeJIRA  ProviderType = "jira"
)
