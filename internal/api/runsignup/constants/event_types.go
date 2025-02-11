package constants

// ValidEventTypes lists all acceptable event_type values for RunSignup API
var ValidEventTypes = map[string]bool{
	"running_race":       true,
	"virtual_race":       true,
	"running_only":       true,
	"walking_only":       true,
	"race_walk":          true,
	"triathlon":          true,
	"duathlon":           true,
	"bike_race":          true,
	"bike_ride":          true,
	"trail_race":         true,
	"open_course_trail":  true,
	"ultra":              true,
	"obstacle_course":    true,
	"swim":               true,
	"swim_run":           true,
	"aqua_bike":          true,
}
