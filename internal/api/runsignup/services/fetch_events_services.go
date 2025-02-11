package services

import (
	"fmt"
)


func FetchEventsByState(state string) (string, error) {
	fmt.Println("Fetching events for state:", state)
	return "Event list for " + state, nil
}
