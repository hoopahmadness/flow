package flowchart

import (
	"errors"
)

type Transition struct {
	Name       string                      `json:"name"`
	NextStages map[ValidationString]string `json:"nextStages"`
}

func newTransition(name string) Transition {
	return Transition{
		Name:       name,
		NextStages: map[ValidationString]string{},
	}
}

func (t *Transition) AddStage(stage string, valTable ValidationTable) {
	t.NextStages[valTable.toString()] = stage
}

func (t Transition) getOutcome(incomingTable ValidationTable) (string, error) {
	for canonVals, status := range t.NextStages {
		canonTable, err := canonVals.toTable()
		if err != nil {
			return INVALID, err
		}
		matches := incomingTable.meetsRequirementsOf(canonTable)
		if matches {
			return status, nil
		}
	}
	return INVALID, errors.New("no outcome found given current validations")
}
