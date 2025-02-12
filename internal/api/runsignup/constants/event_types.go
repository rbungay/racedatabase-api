package constants


type EventCategory string




const (
	CategoryRuns      EventCategory = "Runs"
	CategoryWalks     EventCategory = "Walks"
	CategoryObstacle  EventCategory = "Obstacle"
	CategoryBike      EventCategory = "Bike"
	CategorySwim      EventCategory = "Swim"
	CategoryTriathlon EventCategory = "Triathlon"
	CategoryOther     EventCategory = "Other"
)


var EventTypeToCategory = map[string]EventCategory{
	"running_race":      CategoryRuns,
	"virtual_race":      CategoryRuns,
	"running_only":      CategoryRuns,
	"trail_race":        CategoryRuns,
	"ultra":             CategoryRuns,
	"open_course_trail": CategoryRuns,
	"walking_only":      CategoryWalks,
	"race_walk":         CategoryWalks,
	"obstacle_course":   CategoryObstacle,
	"bike_race":         CategoryBike,
	"bike_ride":         CategoryBike,
	"swim":              CategorySwim,
	"aqua_bike":         CategoryOther,
	"duathlon":          CategoryOther,
	"swim_run":          CategoryOther,
	"triathlon":         CategoryTriathlon,
}


var ValidEventTypes = map[string]bool{
	"running_race":      true,
	"virtual_race":      true,
	"running_only":      true,
	"trail_race":        true,
	"ultra":             true,
	"open_course_trail": true,
	"walking_only":      true,
	"race_walk":         true,
	"obstacle_course":   true,
	"bike_race":         true,
	"bike_ride":         true,
	"swim":              true,
	"aqua_bike":         true,
	"duathlon":          true,
	"swim_run":          true,
	"triathlon":         true,
}
