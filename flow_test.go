package flowchart

import (
	"fmt"
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

	// simplified action for an alternate flow where one action is meant to be used across various stages
	actionAge = "age" // replaces hatch, grow, and emerge but distinct from seen
)

type Butterfly struct {
	color     string
	lifeStage string
	cocoonAge int
}

// These three functions let the Butterfly struct implement the Flowable interface
func (bug *Butterfly) GetStatus() (string, error) {
	return bug.lifeStage, nil
}

func (bug *Butterfly) SetStatus(status, action string) error {
	if bug.lifeStage == stageCocoon && action == actionAge {
		bug.cocoonAge++
	}
	bug.lifeStage = status
	return nil
}

// For a Butterfly, these are the important values to validate. This function should validate any conditional that will be
// used by any one of your Transitions. A Transition might check one or more of these values; any extra flags won't cause trouble.
func (bug *Butterfly) GetContext() (ValidationTable, error) {
	greenTag := "isGreen"
	isGreen := bug.color == "green"

	brownTag := "isBrown"
	isBrown := bug.color == "brown"

	finishedTag := "isFinishedMetamorphosing"
	isFinished := bug.cocoonAge > 0 && bug.lifeStage == stageCocoon

	unrelatedTag := "isAdult"
	isAdult := bug.lifeStage == stageButterfly || bug.lifeStage == stageMoth

	return NewValidationTable(
		greenTag, isGreen,
		brownTag, isBrown,
		finishedTag, isFinished,
		unrelatedTag, isAdult,
	)
}

// As an example, we will generate a flow for a butterfly
// A butterfly starts as an egg, and must progress to the next stage with the appropriate action
// At almost any point it can be seen and eaten by a bird, but this only happens to non-green ones
// coccoons are also safe from being eaten
// Some butterflies (the brown ones) are secretly moths! So when they EMERGE they are moths, not butterflies
func generateGranularFlow() Flow[*Butterfly] {
	// generate a flow object with setters and getters for butterfly struct
	tempButterflyFlow := NewFlow[*Butterfly]()

	// Generate all stages
	eggStage := NewStage(stageEgg)

	catStage := NewStage(stageCaterpillar)

	cocoonStage := NewStage(stageCocoon)

	butterflyStage := NewStage(stageButterfly) // butterflies and moths can't grow any more, they just eventually get seen and eaten

	mothStage := NewStage(stageMoth)

	eatenStage := NewStage(stageEaten)

	// generate some validation tables
	blankTable, _ := NewValidationTable()                    // for all those actions that don't need extra validation
	seenValidator, _ := NewValidationTable("isGreen", false) // you can only be seen if you aren't green
	mothValidator, _ := NewValidationTable("isBrown", true)  // only brown bugs turn into moths
	// if you want to make tweaks to an existing table instead of starting over, make a copy
	mothInvalid := mothValidator.MakeCopy()
	mothInvalid.AddFlag("isBrown", false)

	// generate all transitions, adding validation tables as we go
	hatchTran := NewTransition(actionHatch)
	err := hatchTran.AddStage(&eggStage, blankTable, catStage)
	if err != nil {
		fmt.Println(err)
	}

	growTran := NewTransition(actionGrow)
	err = growTran.AddStage(&catStage, blankTable, cocoonStage)
	if err != nil {
		fmt.Println(err)
	}

	emergeTran := NewTransition(actionEmerge)
	err = emergeTran.AddStage(&cocoonStage, mothInvalid, butterflyStage, mothValidator, mothStage)
	if err != nil {
		fmt.Println(err)
	}

	seenTran := NewTransition(actionSeen) // there's only one possible outcome but we still use a validation table
	err = seenTran.AddStage(&eggStage, seenValidator, eatenStage)
	if err != nil {
		fmt.Println(err)
	}
	err = seenTran.AddStage(&catStage, seenValidator, eatenStage)
	if err != nil {
		fmt.Println(err)
	}
	// cocoons can't be seen, so we don't add the Cocoon stage
	err = seenTran.AddStage(&butterflyStage, seenValidator, eatenStage)
	if err != nil {
		fmt.Println(err)
	}
	err = seenTran.AddStage(&mothStage, seenValidator, eatenStage)
	if err != nil {
		fmt.Println(err)
	}

	// add all the stages and transitions to the flow
	tempButterflyFlow.AddStages(eggStage, catStage, cocoonStage, butterflyStage, mothStage)
	tempButterflyFlow.AddTransitions(hatchTran, growTran, emergeTran, seenTran)

	// you can't use a flow until you Finish it
	butterflyFlow := tempButterflyFlow.Finish()

	return butterflyFlow
}

// Same flow as above, but all growing actions are replaced with the simplified Age action
func generateSimpleFlow() Flow[*Butterfly] {
	// generate a flow object with setters and getters for butterfly struct
	tempButterflyFlow := NewFlow[*Butterfly]()

	// Generate all stages
	eggStage := NewStage(stageEgg)

	catStage := NewStage(stageCaterpillar)

	cocoonStage := NewStage(stageCocoon)

	butterflyStage := NewStage(stageButterfly) // butterflies and moths can't grow any more, they just eventually get seen and eaten

	mothStage := NewStage(stageMoth)

	eatenStage := NewStage(stageEaten)

	// generate some validation tables
	blankTable, _ := NewValidationTable()                    // for all those actions that don't need extra validation
	seenValidator, _ := NewValidationTable("isGreen", false) // you can only be seen if you aren't green
	mothValidator, _ := NewValidationTable(
		"isFinishedMetamorphosing", true,
		"isBrown", true,
	) // only brown bugs turn into moths
	// if you want to make tweaks to an existing table instead of starting over, make a copy
	mothInvalid := mothValidator.MakeCopy()
	mothInvalid.AddFlag("isBrown", false)

	needsMoreTimeValidator, _ := NewValidationTable("isFinishedMetamorphosing", false)

	// generate all transitions, adding validation tables as we go
	ageTran := NewTransition(actionAge)
	err := ageTran.AddStage(&eggStage, blankTable, catStage)
	if err != nil {
		fmt.Println(err)
	}
	err = ageTran.AddStage(&catStage, blankTable, cocoonStage)
	if err != nil {
		fmt.Println(err)
	}
	err = ageTran.AddStage(&cocoonStage,
		mothInvalid, butterflyStage,
		mothValidator, mothStage,
		needsMoreTimeValidator, cocoonStage,
	)
	if err != nil {
		fmt.Println(err)
	}

	seenTran := NewTransition(actionSeen) // there's only one possible outcome but we still use a validation table
	err = seenTran.AddStage(&eggStage, seenValidator, eatenStage)
	if err != nil {
		fmt.Println(err)
	}
	err = seenTran.AddStage(&catStage, seenValidator, eatenStage)
	if err != nil {
		fmt.Println(err)
	}
	// cocoons can't be seen, so we don't add the Cocoon stage
	err = seenTran.AddStage(&butterflyStage, seenValidator, eatenStage)
	if err != nil {
		fmt.Println(err)
	}
	err = seenTran.AddStage(&mothStage, seenValidator, eatenStage)
	if err != nil {
		fmt.Println(err)
	}

	// add all the stages and transitions to the flow
	tempButterflyFlow.AddStages(eggStage, catStage, cocoonStage, butterflyStage, mothStage)
	tempButterflyFlow.AddTransitions(ageTran, seenTran)

	// you can't use a flow until you Finish it
	butterflyFlow := tempButterflyFlow.Finish()

	return butterflyFlow
}

type butterflyTest struct {
	action    string
	result    string
	wantError bool
}

func runButterflyTests(bug *Butterfly, testBatch []butterflyTest, generateFlow func() Flow[*Butterfly], t *testing.T) {
	flow := generateFlow()

	for _, test := range testBatch {
		change, err := flow.TakeAction(bug, test.action)
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
	happyPathGranular := []butterflyTest{
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
	happyPathSimplified := []butterflyTest{
		{
			action:    actionAge,
			result:    stageCaterpillar,
			wantError: false,
		},
		{
			action:    actionAge,
			result:    stageCocoon,
			wantError: false,
		},
		{
			action:    actionAge,
			result:    stageCocoon,
			wantError: false,
		},
		{
			action:    actionAge,
			result:    stageButterfly,
			wantError: false,
		},
		{
			action:    actionAge,
			result:    INVALID,
			wantError: true,
		},
	}
	Rodney := Butterfly{
		color:     "yellow",
		lifeStage: stageEgg,
	}
	Riley := Butterfly{
		color:     "yellow",
		lifeStage: stageEgg,
	}

	runButterflyTests(&Rodney, happyPathGranular, generateGranularFlow, t)
	runButterflyTests(&Riley, happyPathSimplified, generateSimpleFlow, t)
}

func TestSafeButterfliesBranch(t *testing.T) {
	// one bug becomes a butterfly, the other a moth
	butterflyPathGranular := []butterflyTest{
		{
			action:    actionEmerge,
			result:    stageButterfly,
			wantError: false,
		},
	}
	butterflyPathSimplified := []butterflyTest{
		{
			action:    actionAge,
			result:    stageButterfly,
			wantError: false,
		},
	}
	Janice := Butterfly{
		color:     "green",
		lifeStage: stageCocoon,
		cocoonAge: 1,
	}
	Jennifer := Butterfly{
		color:     "green",
		lifeStage: stageCocoon,
		cocoonAge: 1,
	}

	runButterflyTests(&Janice, butterflyPathGranular, generateGranularFlow, t)
	runButterflyTests(&Jennifer, butterflyPathSimplified, generateSimpleFlow, t)

	mothPathGranular := []butterflyTest{
		{
			action:    actionEmerge,
			result:    stageMoth,
			wantError: false,
		},
	}
	mothPathSimplified := []butterflyTest{
		{
			action:    actionAge,
			result:    stageMoth,
			wantError: false,
		},
	}
	Sandy := Butterfly{
		color:     "brown",
		lifeStage: stageCocoon,
		cocoonAge: 1,
	}
	Sidney := Butterfly{
		color:     "brown",
		lifeStage: stageCocoon,
		cocoonAge: 1,
	}

	runButterflyTests(&Sandy, mothPathGranular, generateGranularFlow, t)
	runButterflyTests(&Sidney, mothPathSimplified, generateSimpleFlow, t)
}

func TestSafeButterfliesGetEaten(t *testing.T) {
	// Quincy and Quinton are safe as cocoons but don't make it as butterflies
	eatenPathGranular := []butterflyTest{
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
	eatenPathSimplified := []butterflyTest{
		{
			action:    actionSeen,
			result:    INVALID,
			wantError: true,
		},
		{
			action:    actionAge,
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
		cocoonAge: 1,
	}
	Quinton := Butterfly{
		color:     "red",
		lifeStage: stageCocoon,
		cocoonAge: 1,
	}
	runButterflyTests(&Quincy, eatenPathGranular, generateGranularFlow, t)
	runButterflyTests(&Quinton, eatenPathSimplified, generateSimpleFlow, t)

}
