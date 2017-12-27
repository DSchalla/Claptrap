package rules

type Case struct {
	Name       string
	Conditions []Condition
	Responses  []Response
}
