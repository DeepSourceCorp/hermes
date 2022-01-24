package models

type Analyzer struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Shortcode string `json:"shortcode"`
	Name      string `json:"name"`
	LogoURL   string `json:"logoUrl"`
}
