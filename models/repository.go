package models

type Repository struct {
	ID              string  `json:"id"`
	Object          string  `json:"object"`
	Name            string  `json:"name"`
	VcsProvider     string  `json:"vcsProvider"`
	VcsURL          string  `json:"vcsUrl"`
	DefaultBranch   string  `json:"defaultBranch"`
	LatestCommitOid string  `json:"latestCommitOid"`
	IsPrivate       bool    `json:"isPrivate"`
	IsActivated     bool    `json:"isActivated"`
	Account         Account `json:"account"`
}
