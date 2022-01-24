package models

type Issue struct {
	ID               string   `json:"id"`
	Object           string   `json:"object"`
	Shortcode        string   `json:"shortcode"`
	Title            string   `json:"title"`
	ShortDescription string   `json:"shortDescription"`
	Category         string   `json:"category"`
	AutofixAvailable bool     `json:"autofixAvailable"`
	IsRecommended    bool     `json:"isRecommended"`
	Analyzer         Analyzer `json:"analyzer"`
}
