package flowchart

import (
	"errors"
	"fmt"
)

type Transition struct {
	Name       string                      `json:"name"`
	NextStages map[ValidationString]string `json:"nextStages"`
}

func NewTransition(name string) Transition {
	return Transition{
		Name:       name,
		NextStages: map[ValidationString]string{},
	}
}

func (t *Transition) AddStage(originStage *Stage, nextSteps ...interface{}) error {
	if originStage == nil {
		return fmt.Errorf("Unable to add stage with nil origin")
	}
	if len(nextSteps)%2 != 0 {
		return fmt.Errorf("Pairs of validation tables and destination stages are required for next steps")
	}
	originStage.addTransition(t.Name)
	for ii := 0; ii < len(nextSteps); ii += 2 {
		valTable, OK := nextSteps[ii].(ValidationTable)
		if !OK {
			return fmt.Errorf("Expected a valudation table, got %T", nextSteps[ii])
		}
		nextStage, OK := nextSteps[ii+1].(Stage)
		if !OK {
			return fmt.Errorf("Expected a destination stage, got %T", nextSteps[ii+1])
		}
		originFlag := fmt.Sprintf(originStageFlag, originStage.Name)
		valTable.AddFlag(originFlag, true)
		t.NextStages[valTable.toString()] = nextStage.Name
	}
	return nil
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
