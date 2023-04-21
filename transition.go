package flowchart

import (
	"errors"
	"fmt"
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

func (t *Transition) AddStage(stage string, valTable ValidationTable, originStages ...string) {
	if len(originStages) == 0 {
		t.NextStages[valTable.toString()] = stage
		return
	}
	for _, origin := range originStages {
		valTableWithOrigin := valTable.MakeCopy()
		originFlag := fmt.Sprintf(originStageFlag, origin)
		valTableWithOrigin.AddFlag(originFlag, true)
		t.NextStages[valTableWithOrigin.toString()] = stage
	}
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
