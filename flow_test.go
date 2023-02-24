package flowchart

import (
	"errors"
	"testing"
)

const (
	// stages
	stageEgg         = "egg"
	stageCaterpillar = "caterpillar"
	stageCocoon      = "cocoon"
	stageButterfly   = "butterfly"
	stageMoth        = "moth"
	stageEaten       = "eaten"

	// actions
	actionHatch  = "hatch"
	actionGrow   = "grow"
	actionEmerge = "emerge"
	actionSeen   = "seen"
)

type Butterfly struct {
	color     string
	lifeStage string
}

// These three functions will be passed into the NewFlow() function, allowing it to interact with our asset.
func getButterflyStatus(asset interface{}) (string, error) {
	bug, OK := asset.(*Butterfly)
	if !OK {
		return "", errors.New("")
	}
	return bug.lifeStage, nil
}

func setButterflyStatus(asset interface{}, status string) error {
	bug, OK := asset.(*Butterfly)
	if !OK {
		return errors.New("")
	}
	bug.lifeStage = status
	return nil
}

// For a Butterfly, these are the important values to validate. This function should validate any conditional that will be
// used by any one of your Transitions. A Transition might check one or more of these values; any extra flags won't cause trouble.
func getButterflyContext(asset interface{}) (ValidationTable, error) {
	bug, OK := asset.(*Butterfly)
	if !OK {
		return ValidationTable{}, errors.New("")
	}

	greenTag := "isGreen"
	isGreen := bug.color == "green"

	brownTag := "isBrown"
	isBrown := bug.color == "brown"

	unrelatedTag := "isAdult"
	isAdult := bug.lifeStage == stageButterfly || bug.lifeStage == stageMoth

	return NewValidationTable(greenTag, isGreen, brownTag, isBrown, unrelatedTag, isAdult)
}

// As an example, we will generate a flow for a butterfly
// A butterfly starts as an egg, and must progress to the next stage with the appropriate action
// At almost any point it can be seen and eaten by a bird, but this only happens to non-green ones
// coccoons are also safe from being eaten
// Some butterflies (the brown ones) are secretly moths! So when they EMERGE they are moths, not butterflies
// Make sure your strings match up! Using constants is best
func generateFlow() Flow {
	// generate a flow object with setters and getters for butterfly struct
	tempButterflyFlow := NewFlow(getButterflyStatus, setButterflyStatus, getButterflyContext)

	// Generate all stages, adding transiton names as we go
	eggStage := NewStage(stageEgg)
	eggStage.AddTransition(actionHatch)
	eggStage.AddTransition(actionSeen)

	catStage := NewStage(stageCaterpillar)
	catStage.AddTransition(actionGrow)
	catStage.AddTransition(actionSeen)

	cocoonStage := NewStage(stageCocoon) // cocoons can't be seen, so we don't add the SEEN transition
	cocoonStage.AddTransition(actionEmerge)

	butterflyStage := NewStage(stageButterfly) // butterflies can't grow any more, they just eventually get seen and eaten
	butterflyStage.AddTransition(actionSeen)

	mothStage := NewStage(stageMoth)
	mothStage.AddTransition(actionSeen)

	// generate some validation tables
	blankTable, _ := NewValidationTable()                    // for all those actions that don't need extra validation
	seenValidator, _ := NewValidationTable("isGreen", false) // you can only be seen if you aren't green
	mothValidator, _ := NewValidationTable("isBrown", true)  // only brown bugs turn into moths
	// if you want to make tweaks to an existing table instead of starting over, make a copy
	mothInvalid := mothValidator.MakeCopy()
	mothInvalid.AddFlag("isBrown", false)

	// generate all transitions, adding optional validation tables as we go
	hatchTran := newTransition(actionHatch)
	hatchTran.AddStage(stageCaterpillar, blankTable) // there's no validation needed to hatch

	growTran := newTransition(actionGrow)
	growTran.AddStage(stageCocoon, blankTable)

	emergeTran := newTransition(actionEmerge)
	emergeTran.AddStage(stageButterfly, mothInvalid)
	emergeTran.AddStage(stageMoth, mothValidator) // since EMERGE has two possible outcomes, we use validation tables to choose between them

	seenTran := newTransition(actionSeen)        // there's only one possible outcome but we still use a validation table
	seenTran.AddStage(stageEaten, seenValidator) // attempting to SEEN a green bug will return an error

	// add all the stages and transitions to the flow
	tempButterflyFlow.AddStages(eggStage, catStage, cocoonStage, butterflyStage, mothStage)
	tempButterflyFlow.AddTransitions(hatchTran, growTran, emergeTran, seenTran)

	// you can't use a flow until you Finish it
	butterflyFlow := tempButterflyFlow.Finish()

	return butterflyFlow
}

type butterflyTest struct {
	action    string
	result    string
	wantError bool
}

func runButterflyTests(bug *Butterfly, testBatch []butterflyTest, t *testing.T) {
	flow := generateFlow()

	for _, test := range testBatch {
		change, err := flow.CheckRequest(bug, test.action)
		if err != nil && !test.wantError {
			t.Error(err)
			t.Fail()
		}
		if err == nil && test.wantError {
			t.Error("Expected error but got none")
			t.Fail()
		}
		if change != test.result {
			t.Errorf("Wanted %s, got %s.", test.result, change)
			t.FailNow()
		}
		if bug.lifeStage != change && !test.wantError {
			t.Errorf("Asset status not correctly being updated: currently %s, should be %s.", bug.lifeStage, change)
			t.FailNow()
		}
	}

}

func TestSafeButterfliesHappy(t *testing.T) {
	// test happy path, show terminal stage (can't progress further)
	happyPath := []butterflyTest{
		{
			action:    actionHatch,
			result:    stageCaterpillar,
			wantError: false,
		},
		{
			action:    actionGrow,
			result:    stageCocoon,
			wantError: false,
		},
		{
			action:    actionEmerge,
			result:    stageButterfly,
			wantError: false,
		},
		{
			action:    actionEmerge,
			result:    INVALID,
			wantError: true,
		},
	}
	Rodney := Butterfly{
		color:     "yellow",
		lifeStage: stageEgg,
	}

	runButterflyTests(&Rodney, happyPath, t)
}

func TestSafeButterfliesBranch(t *testing.T) {
	// one bug becomes a butterfly, the other a moth
	butterflyPath := []butterflyTest{
		{
			action:    actionEmerge,
			result:    stageButterfly,
			wantError: false,
		},
	}
	Janice := Butterfly{
		color:     "green",
		lifeStage: stageCocoon,
	}

	runButterflyTests(&Janice, butterflyPath, t)

	mothPath := []butterflyTest{
		{
			action:    actionEmerge,
			result:    stageMoth,
			wantError: false,
		},
	}
	Sandy := Butterfly{
		color:     "brown",
		lifeStage: stageCocoon,
	}

	runButterflyTests(&Sandy, mothPath, t)
}

func TestSafeButterfliesGetEaten(t *testing.T) {
	// Quincy is safe as a cocoon but doesn't make it as a butterfly
	eatenPath := []butterflyTest{
		{
			action:    actionSeen,
			result:    INVALID,
			wantError: true,
		},
		{
			action:    actionEmerge,
			result:    stageButterfly,
			wantError: false,
		},
		{
			action:    actionSeen,
			result:    stageEaten,
			wantError: false,
		},
	}
	Quincy := Butterfly{
		color:     "red",
		lifeStage: stageCocoon,
	}
	runButterflyTests(&Quincy, eatenPath, t)

}
