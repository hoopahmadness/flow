package flowchart

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ValidationTable struct {
	table map[string]bool
	tags  []string
}

type ValidationString string

func NewValidationTable(args ...interface{}) (ValidationTable, error) {
	newTable := ValidationTable{
		table: map[string]bool{},
		tags:  []string{},
	}

	if len(args) == 0 {
		return newTable, nil
	}
	if len(args)%2 != 0 {
		return ValidationTable{}, errors.New("an even number of arguments is required for NewValidationTable")
	}

	for ii := 0; ii < len(args)-1; ii += 2 {
		tag, tagOK := args[ii].(string)
		if !tagOK {
			return newTable, errors.New("didn't get type string as expected")
		}

		flag, flagOK := args[ii+1].(bool)
		if !flagOK {
			return newTable, errors.New("didn't get type boolean as expected")
		}
		newTable.AddFlag(tag, flag)
	}
	return newTable, nil
}

func (vt ValidationTable) MakeCopy() ValidationTable {
	str := vt.toString()
	table, _ := str.toTable()
	return table
}

func (vt ValidationTable) toString() ValidationString {
	if len(vt.tags) == 0 {
		return ValidationString(" ")
	}

	out := []string{}
	for _, tag := range vt.tags {
		out = append(out, fmt.Sprintf("%s:%t", tag, vt.table[tag]))
	}

	return ValidationString(strings.Join(out, ","))
}

func (vt *ValidationTable) AddFlag(tag string, flag bool) {
	// add tag to values array, sorted
	previousValue := ""
Loop:
	for index, listValue := range vt.tags {
		// for scenarios with a non-zero array of tags
		switch {
		case listValue == tag: // our tag is already in the list; do not update array
			break Loop
		case previousValue == "" && tag < listValue: // special case; tag fits at the very front
			vt.tags = append([]string{tag}, vt.tags...)
			break Loop
		case tag > listValue && index == len(vt.tags)-1: // special case where our tag fits at the very end
			vt.tags = append(vt.tags, tag)
		case previousValue < tag && tag > listValue: // tag fits between these two indices; moving on
			previousValue = listValue
		case tag < listValue: // we have passed our insertion point

			vt.tags = append(vt.tags, "")
			copy(vt.tags[index+1:], vt.tags[index:])
			vt.tags[index] = tag
			break Loop
		}
	}

	// special case where the array is 0 length
	if len(vt.tags) == 0 {
		vt.tags = []string{tag}
	}

	// update/add flag to map
	vt.table[tag] = flag
}

func (vt ValidationTable) meetsRequirementsOf(incoming ValidationTable) bool {
	for _, tag := range incoming.tags {
		ourFlag, exists := vt.table[tag]
		if !exists {
			return false
		}
		if incoming.table[tag] != ourFlag {
			return false
		}
	}
	return true
}

func (valStr ValidationString) toTable() (ValidationTable, error) {
	table, _ := NewValidationTable()
	if string(valStr) == " " {
		return table, nil
	}
	pairs := strings.Split(string(valStr), ",")

	for _, pair := range pairs {
		tokens := strings.Split(pair, ":")
		tag := tokens[0]
		flag, err := strconv.ParseBool(tokens[1])
		if err != nil {
			return table, err
		}
		table.AddFlag(tag, flag)
	}

	return table, nil
}
