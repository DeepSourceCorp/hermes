package models

type Account struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	Login       string `json:"login"`
	VcsProvider string `json:"vcsProvider"`
	VcsURL      string `json:"vcsUrl"`
	Type        string `json:"type"`
	AvatarURL   string `json:"avatarUrl"`
}
