package flowchart

import (
	"testing"
)

func TestSafeMatch(t *testing.T) {
	stage1 := NewStage("stage1")
	stage2 := NewStage("stage2")
	stage3 := NewStage("stage3")
	stage4 := NewStage("stage4")

	blankVal, _ := NewValidationTable()
	tableA, _ := NewValidationTable("ConstraintA", true)
	tableAf, _ := NewValidationTable("ConstraintA", false)

	tableB, _ := NewValidationTable("ConstraintB", true)
	tableBf, _ := NewValidationTable("ConstraintB", false)

	tableC, _ := NewValidationTable("ConstraintC", true)
	tableCf, _ := NewValidationTable("ConstraintC", false)

	tranBlankFirst := NewTransition("transBlankFirst")
	tranBlankFirst.AddStage(&stage1,
		blankVal, stage2,
		tableA, stage3,
		tableB, stage3,
		tableC, stage3,
		tableAf, stage3,
		tableBf, stage3,
		tableCf, stage3,
	)
	transComplexFirst := NewTransition("transComplexFirst")
	transComplexFirst.AddStage(&stage1,
		tableA.And(tableB).And(tableC), stage1,
		tableA.And(tableB), stage2,
		tableA, stage3,
		blankVal, stage4,
	)
	type MatchingTest struct {
		description   string
		transition    Transition
		inputTable    ValidationTable
		expectedMatch string
		expectedErr   error
	}

	testSuite := []MatchingTest{
		{
			description:   "Blank outcome comes first and always matches",
			transition:    tranBlankFirst,
			inputTable:    tableA,
			expectedMatch: stage2.Name,
			expectedErr:   nil,
		},
		{
			description:   "Complex table matches over blank table because it's first",
			transition:    transComplexFirst,
			inputTable:    tableA.And(tableB).And(tableC),
			expectedMatch: stage1.Name,
			expectedErr:   nil,
		},
		{
			description:   "Try all combos, set 1",
			transition:    transComplexFirst,
			inputTable:    tableA,
			expectedMatch: stage3.Name,
			expectedErr:   nil,
		},
		{
			description:   "Try all combos, set 2",
			transition:    transComplexFirst,
			inputTable:    tableAf,
			expectedMatch: stage4.Name,
			expectedErr:   nil,
		},
		{
			description:   "Try all combos, set 3",
			transition:    transComplexFirst,
			inputTable:    tableB,
			expectedMatch: stage4.Name,
			expectedErr:   nil,
		},
		{
			description:   "Try all combos, set 4",
			transition:    transComplexFirst,
			inputTable:    tableA.And(tableBf),
			expectedMatch: stage3.Name,
			expectedErr:   nil,
		},
	}

	for _, testCase := range testSuite {
		for ii := 0; ii < 100; ii++ {
			outcome, err := testCase.transition.getOutcome(testCase.inputTable)
			if err != testCase.expectedErr {
				t.Errorf("Test failed for %s", testCase.description)
				t.Errorf("Expected errors did not match: wanted %v, got %v", testCase.expectedErr, err)
				break
			}
			if outcome != testCase.expectedMatch {
				t.Errorf("Test failed for %s", testCase.description)
				t.Errorf("Expected matched stage did not match: wanted %s, got %s", testCase.expectedMatch, outcome)
				break
			}
		}
	}
}
