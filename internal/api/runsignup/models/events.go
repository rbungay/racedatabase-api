package models

import "github.com/rbungay/racedatabase-api/internal/api/runsignup/constants"

type Event struct {
	ID          int    `json:"race_id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	ExternalURL string `json:"external_race_url"` 
	LogoURL     string `json:"logo_url"`          
	Category    constants.EventCategory `json:"category"`
}
