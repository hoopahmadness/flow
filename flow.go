package flowchart

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type UnfinishedFlow struct {
	Stages        map[string]Stage
	Transitions   map[string]Transition
	statusGetter  func(interface{}) (string, error)
	statusSetter  func(interface{}, string) error
	contextGetter func(interface{}) (ValidationTable, error)
}
type Flow struct {
	stages        map[string]Stage
	transitions   map[string]Transition
	statusGetter  func(interface{}) (string, error)
	statusSetter  func(interface{}, string) error
	contextGetter func(interface{}) (ValidationTable, error)
}

func (f UnfinishedFlow) Finish() Flow {
	newFlow := Flow{
		stages:        f.Stages,
		transitions:   f.Transitions,
		statusGetter:  f.statusGetter,
		statusSetter:  f.statusSetter,
		contextGetter: f.contextGetter,
	}
	return newFlow
}

func NewFlow(
	statusGetter func(interface{}) (string, error),
	statusSetter func(interface{}, string) error,
	contextGetter func(interface{}) (ValidationTable, error),
) UnfinishedFlow {
	return UnfinishedFlow{
		Stages:        map[string]Stage{},
		Transitions:   map[string]Transition{},
		statusGetter:  statusGetter,
		statusSetter:  statusSetter,
		contextGetter: contextGetter,
	}
}

func (f *UnfinishedFlow) AddStages(stages ...Stage) {
	for _, stage := range stages {
		f.Stages[stage.Name] = stage
	}
}

func (f *UnfinishedFlow) AddTransitions(transitions ...Transition) {
	for _, transition := range transitions {
		f.Transitions[transition.Name] = transition
	}
}

func (f Flow) TakeAction(asset interface{}, action string) (string, error) {
	// check if asset is a pointer
	if !isPointer(asset) {
		return INVALID, fmt.Errorf("please pass a pointer to your asset in CheckRequest()")
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
	return reflect.ValueOf(asset).Type().Kind() == reflect.Ptr
}
