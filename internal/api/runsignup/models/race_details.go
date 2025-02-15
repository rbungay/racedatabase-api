package models


type RaceDetails struct {
	ID         int             `json:"race_id"`
	Name       string          `json:"name"`
	URL        string          `json:"url"`
	ExternalURL string         `json:"external_race_url"`
	LogoURL    string          `json:"logo_url"`
	Timezone   string          `json:"timezone"`
	Events     []EventDetails  `json:"events"`
}


type EventDetails struct {
	EventID    int                  `json:"event_id"`
	Name       string               `json:"name"`
	StartTime  string               `json:"start_time"`
	EndTime    string               `json:"end_time"`
	EventType  string               `json:"event_type"`
	Distance   string               `json:"distance"`
	RegOpens   string               `json:"registration_opens"`
	Category   string               `json:"category"`
	RegPeriods []RegistrationPeriod `json:"registration_periods"`
}


type RegistrationPeriod struct {
	Opens    string `json:"registration_opens"`
	Closes   string `json:"registration_closes"`
	Fee      string `json:"race_fee"`
	ProcFee  string `json:"processing_fee"`
}
