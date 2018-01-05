package rules

type Case struct {
	Name         string
	Conditions   []Condition
	Responses    []Response
	ResponseFunc func(event Event, rh ResponseHandler) bool
}
