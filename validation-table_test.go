package flowchart

import (
	"testing"
)

func TestSafeValidationTableCreation(t *testing.T) {
	type creationTest struct {
		note      string
		args      []interface{}
		wantError bool
	}
	creationTests := []creationTest{
		{
			note:      "create a blank table",
			args:      []interface{}{},
			wantError: false,
		},
		{
			note:      "create a table with tags",
			args:      []interface{}{"first", true, "second", false, "third", true},
			wantError: false,
		},
		{
			note:      "create a table with badly typed tags",
			args:      []interface{}{"first", "true"},
			wantError: true,
		},
		{
			note:      "create a table with odd number of tags, want error",
			args:      []interface{}{"first", true, "second"},
			wantError: true,
		},
	}

	for _, test := range creationTests {
		table, err1 := NewValidationTable(test.args...)
		if test.wantError {
			if err1 == nil {
				t.Errorf("Expected error but none appeared for test %s", test.note)
				t.FailNow()
			} else {
				continue
			}
		} else if err1 != nil {
			t.Errorf("test: %s \n %v", test.note, err1)
			t.FailNow()
		}

		// create a table from validationString
		valStr := table.toString()
		copyTable, err2 := valStr.toTable()
		if err2 != nil {
			t.Error(err2)
			t.FailNow()
		}

		// manually add a tag to check for panics
		table.AddFlag("didn't panic", true)
		copyTable.AddFlag("didn't panic", true)
	}
}

func TestSafeValidationTableMeetRequirements(t *testing.T) {
	// check if a larger table meetsRequirements for a smaller one
	bigTable, _ := NewValidationTable("first", true, "second", true, "third", false, "fourth", false, "fifth", true)
	smallTable, _ := NewValidationTable("first", true, "second", true, "third", false)

	if !bigTable.meetsRequirementsOf(smallTable) {
		t.Errorf("large validation table didn't match smaller subset like it should")
	}

	// check that missing tags does not meet requirements
	missingFalse, _ := NewValidationTable("first", true, "second", true, "third", false, "fifth", true)
	if missingFalse.meetsRequirementsOf(bigTable) {
		t.Errorf("missing false flag still matched")

	}
	missingTrue, _ := NewValidationTable("first", true, "second", true, "third", false, "fourth", false)
	if missingTrue.meetsRequirementsOf(bigTable) {
		t.Errorf("missing true flag still matched")
	}

	// check that wrong flags does not meet requirements
	badFlag, _ := NewValidationTable("first", false)
	if badFlag.meetsRequirementsOf(bigTable) {
		t.Errorf("incorrect flag value still matched")
	}

	// check that tag-adding-order does not change string output
	scrambled, _ := NewValidationTable("fourth", false, "third", false, "first", true, "fifth", true, "second", true)
	if bigTable.toString() != scrambled.toString() {
		t.Errorf("tag order should not affect output string")

	}
}

func TestSafeCombineValidationTales(t *testing.T) {
	// combine two blank VTs
	blank1, _ := NewValidationTable()
	blank2, _ := NewValidationTable()
	combinedBlank := blank1.And(blank2)
	if combinedBlank.toString() != " " {
		t.Errorf("Expected a new blank table to be created")
	}

	// Combine populated table with blank
	pop1, _ := NewValidationTable("valA", true, "valB", true)
	combined1 := pop1.And(blank1)
	if combined1.toString() != "valA:true,valB:true" {
		t.Errorf("Did not get the expected combination of pop1 and blank1, got %s instead", combined1.toString())
	}
	if blank1.toString() != " " {
		t.Errorf("Combining tables should not change the second starting table")
	}

	combined1b := blank1.And(pop1)
	if combined1b.toString() != combined1.toString() {
		t.Errorf("Expected %s, got %s", combined1.toString(), combined1b.toString())
	}

	if blank1.toString() != " " {
		t.Errorf("Combining tables should not change the first starting table")
	}
	// combine two populated tables
	pop2, _ := NewValidationTable("valB", false, "valC", false, "valD", false)
	combined2 := pop1.And(pop2)
	if combined2.toString() != "valA:true,valB:false,valC:false,valD:false" {
		t.Errorf("Expected 'valA:true,valB:false,valC:false,valD:false', got %s", combined2.toString())
	}

	// combine them backwards
	combined2R := pop2.And(pop1)
	if combined2R.toString() != "valA:true,valB:true,valC:false,valD:false" {
		t.Errorf("Expected second table to overwrite shared values of first table, got %s", combined2R.toString())
	}

	//combine 3 in a chain
	pop3, _ := NewValidationTable("valD", true, "valE", true)
	combined3 := pop1.And(pop2).And(pop3)
	if combined3.toString() != "valA:true,valB:false,valC:false,valD:true,valE:true" {
		t.Errorf("Problem combining 3 tables in a chain, expected valA:true,valB:false,valC:false,valD:true,valE:true got %s", combined3.toString())
	}
}
