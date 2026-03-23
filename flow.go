package flowchart

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type Flowable interface {
	GetStatus(logger FlowLogger) (string, error)
	SetStatus(newStatus string, actionInfo map[string]any, logger FlowLogger) error
	GetContext(logger FlowLogger) (ValidationTable, error)
}

type UnfinishedFlow[Asset Flowable] struct {
	Stages      map[string]Stage
	Transitions map[string]Transition
}
type Flow[Asset Flowable] struct {
	stages      map[string]Stage
	transitions map[string]Transition
}

func (f UnfinishedFlow[Asset]) Finish() Flow[Asset] {
	newFlow := Flow[Asset]{
		stages:      f.Stages,
		transitions: f.Transitions,
	}
	return newFlow
}

func NewFlow[Asset Flowable]() UnfinishedFlow[Asset] {
	return UnfinishedFlow[Asset]{
		Stages:      map[string]Stage{},
		Transitions: map[string]Transition{},
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

func (f Flow[Asset]) TakeAction(asset Asset, action string, actionInfo map[string]any, parentLogger FlowLogger) (string, error) {
	flowTypeStr := fmt.Sprintf("%T", f)
	strOpen := strings.Split(flowTypeStr, "[")[1]
	strClosed := strings.Split(strOpen, "]")[0]
	logger := parentLogger.New("caller", "flow.TakeAction",
		"flowType", strClosed,
		"action", action,
	)
	// check if asset is a pointer
	if !isPointer(asset) {
		logger.Log("Please pass a pointer to your asset in TakeAction()")
		return INVALID, fmt.Errorf("please pass a pointer to your asset in TakeAction()")
	}

	if actionInfo == nil {
		actionInfo = map[string]any{}
	}
	actionInfo["action"] = action

	// check if action is part of our flow
	tran, OK := f.transitions[action]
	if !OK {
		return INVALID, fmt.Errorf("given action '%s' is not valid for this flow", action)
	}

	// get current stage and validations
	status, err := asset.GetStatus(logger)
	if err != nil {
		return INVALID, err
	}
	validations, err := asset.GetContext(logger)
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
		if innerErr := asset.SetStatus(newStatus, actionInfo, logger); innerErr != nil {
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
