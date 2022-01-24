package models

type Occurence struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	Title       string `json:"title"`
	Path        string `json:"path"`
	BeginLine   int    `json:"beginLine"`
	BeginColumn int    `json:"beginColumn"`
	EndLine     int    `json:"endLine"`
	EndColumn   int    `json:"endColumn"`
}
