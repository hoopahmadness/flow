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
