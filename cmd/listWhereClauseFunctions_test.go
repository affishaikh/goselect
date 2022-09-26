package cmd

import (
	"bytes"
	"goselect/parser/context"
	"strings"
	"testing"
)

func TestWhereClauseFunctionsUnsorted(t *testing.T) {
	rootCmd.SetArgs([]string{"listWhereClauseFunctions"})
	buffer := new(bytes.Buffer)
	rootCmd.SetOut(buffer)

	_ = rootCmd.Execute()
	contents := buffer.String()

	aliasesByFunction := context.NewFunctions().AllFunctionsWithAliasesHavingTag("where")
	for function, aliases := range aliasesByFunction {
		if !strings.Contains(contents, function) {
			t.Fatalf("Expected function %v to be contained in the where clause supported functions but was not, received %v", function, contents)
		}
		for _, alias := range aliases {
			if !strings.Contains(contents, alias) {
				t.Fatalf("Expected alias %v to be contained in the where clause supported functions but was not, received %v", alias, contents)
			}
		}
	}
}

func TestWhereClauseFunctionsSorted(t *testing.T) {
	rootCmd.SetArgs([]string{"listWhereClauseFunctions", "--sorted", "true"})
	buffer := new(bytes.Buffer)
	rootCmd.SetOut(buffer)

	_ = rootCmd.Execute()
	contents := buffer.String()

	aliasesByFunction := context.NewFunctions().AllFunctionsWithAliasesHavingTag("where")
	for function, aliases := range aliasesByFunction {
		if !strings.Contains(contents, function) {
			t.Fatalf("Expected function %v to be contained in the where clause supported functions but was not, received %v", function, contents)
		}
		for _, alias := range aliases {
			if !strings.Contains(contents, alias) {
				t.Fatalf("Expected alias %v to be contained in the where clause supported functions but was not, received %v", alias, contents)
			}
		}
	}
}

func TestWhereClauseFunctionsUnSortedWithInvalidValue(t *testing.T) {
	rootCmd.SetArgs([]string{"listWhereClauseFunctions", "--sorted", "unknown"})
	buffer := new(bytes.Buffer)
	rootCmd.SetOut(buffer)

	_ = rootCmd.Execute()
	contents := buffer.String()

	aliasesByFunction := context.NewFunctions().AllFunctionsWithAliasesHavingTag("where")
	for function, aliases := range aliasesByFunction {
		if !strings.Contains(contents, function) {
			t.Fatalf("Expected function %v to be contained in the where clause supported functions but was not, received %v", function, contents)
		}
		for _, alias := range aliases {
			if !strings.Contains(contents, alias) {
				t.Fatalf("Expected alias %v to be contained in the where clause supported functions but was not, received %v", alias, contents)
			}
		}
	}
}
