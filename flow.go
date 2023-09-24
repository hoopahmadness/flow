package flowchart

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type UnfinishedFlow[Asset any] struct {
	Stages        map[string]Stage
	Transitions   map[string]Transition
	statusGetter  func(*Asset) (string, error) // perhaps generics would be useful here!
	statusSetter  func(*Asset, string) error
	contextGetter func(*Asset) (ValidationTable, error)
}
type Flow[Asset any] struct {
	stages        map[string]Stage
	transitions   map[string]Transition
	statusGetter  func(*Asset) (string, error)
	statusSetter  func(*Asset, string) error
	contextGetter func(*Asset) (ValidationTable, error)
}

func (f UnfinishedFlow[Asset]) Finish() Flow[Asset] {
	newFlow := Flow[Asset]{
		stages:        f.Stages,
		transitions:   f.Transitions,
		statusGetter:  f.statusGetter,
		statusSetter:  f.statusSetter,
		contextGetter: f.contextGetter,
	}
	return newFlow
}

func NewFlow[Asset any](
	statusGetter func(*Asset) (string, error),
	statusSetter func(*Asset, string) error,
	contextGetter func(*Asset) (ValidationTable, error),
) UnfinishedFlow[Asset] {
	return UnfinishedFlow[Asset]{
		Stages:        map[string]Stage{},
		Transitions:   map[string]Transition{},
		statusGetter:  statusGetter,
		statusSetter:  statusSetter,
		contextGetter: contextGetter,
	}
}

func (f *UnfinishedFlow[Asset]) AddStages(stages ...Stage) {
	for _, stage := range stages {
		f.Stages[stage.Name] = stage
	}
}

func (f *UnfinishedFlow[Asset]) AddTransitions(transitions ...Transition) {
	for _, transition := range transitions {
		f.Transitions[transition.Name] = transition
	}
}

func (f Flow[Asset]) TakeAction(asset *Asset, action string) (string, error) {
	// check if asset is a pointer
	if !isPointer(asset) {
		return INVALID, fmt.Errorf("please pass a pointer to your asset in TakeAction()")
	}

	// check if action is part of our flow
	tran, OK := f.transitions[action]
	if !OK {
		return INVALID, fmt.Errorf("given action '%s' is not valid for this flow", action)
	}

	// get current stage and validations
	status, err := f.statusGetter(asset)
	if err != nil {
		return INVALID, err
	}
	validations, err := f.contextGetter(asset)
	if err != nil {
		return INVALID, err
	}

	// add origin stage flag to our validations
	validations.AddFlag(fmt.Sprintf(originStageFlag, status), true)

	// check if current stage is part of our flow
	stage, OK := f.stages[status]
	if !OK {
		return INVALID, fmt.Errorf("calculated status '%s' is not valid for this flow", status)
	}

	// check if transition is valid for that stage
	if !contains(stage.Transitions, action) {
		return INVALID, fmt.Errorf("given action '%s' is not allowed for the status %s", action, stage)
	}

	newStatus, err := tran.getOutcome(validations)
	if err == nil {
		if innerErr := f.statusSetter(asset, newStatus); innerErr != nil {
			return INVALID, errors.Wrap(innerErr, "call to f.statusSetter failed")
		}
	}

	return newStatus, err

}

func contains(list []string, single string) bool {
	for _, text := range list {
		if text == single {
			return true
		}
	}
	return false
}

func isPointer(asset interface{}) bool {
	return reflect.ValueOf(asset).Type().Kind() == reflect.Pointer
}
