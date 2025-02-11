package models

import "github.com/rbungay/racedatabase-api/internal/api/runsignup/constants"

type Event struct {
	ID        int    `json:"race_id"`
	Name      string `json:"name"`
	StartDate string `json:"next_date"`
	EndDate   string `json:"next_end_date"`
	URL       string `json:"url"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zipcode   string `json:"zipcode"`
	EventType string        `json:"event_type"`  
	Category  constants.EventCategory `json:"category"` 
}
