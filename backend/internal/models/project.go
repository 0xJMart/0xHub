package models

// Project represents a project in the hub
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Icon        string `json:"icon,omitempty"`
	Category    string `json:"category,omitempty"`
	Status      string `json:"status,omitempty"`
}
