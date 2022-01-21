package event

import "github.com/deepsourcelabs/hermes/models"

type RepositoryIssueIntroduced struct {
	Payload RepositoryIssueIntroducedPayload
}

type RepositoryIssueIntroducedPayload struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	CreatedAt int    `json:"createdAt"`
	Data      struct {
		Object struct {
			ID          string             `json:"id"`
			Object      string             `json:"object"`
			Repository  models.Repository  `json:"repository"`
			Issue       models.Issue       `json:"issue"`
			Occurrences []models.Occurence `json:"occurrences"`
		} `json:"object"`
		IssueOccurrencesIntroduced int `json:"issueOccurrencesIntroduced"`
	} `json:"data"`
}
