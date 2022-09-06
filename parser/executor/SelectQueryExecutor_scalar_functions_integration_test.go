package executor

import (
	"fmt"
	"goselect/parser"
	"goselect/parser/context"
	"math"
	"os"
	"testing"
)

func TestResultsWithProjections1(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select name, now() from ../resources/TestResultsWithProjections/single", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("TestResultsWithProjections_A.txt"), context.EmptyValue},
	}
	assertMatch(t, expected, queryResults, 1)
}

func TestResultsWithProjections2(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), base64(name) from ../resources/TestResultsWithProjections/single", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.txt"), context.StringValue("VGVzdFJlc3VsdHNXaXRoUHJvamVjdGlvbnNfQS50eHQ=")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsInNestedDirectories(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select isdir, lower(name), path from ../resources/TestResultsWithProjections/ order by 1 desc, 2 asc", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.BooleanValue(true), context.StringValue("empty"), context.StringValue("../resources/TestResultsWithProjections/empty")},
		{context.BooleanValue(true), context.StringValue("multi"), context.StringValue("../resources/TestResultsWithProjections/multi")},
		{context.BooleanValue(true), context.StringValue("single"), context.StringValue("../resources/TestResultsWithProjections/single")},
		{context.BooleanValue(false), context.StringValue("empty.log"), context.StringValue("../resources/TestResultsWithProjections/empty/Empty.log")},
		{context.BooleanValue(false), context.StringValue("testresultswithprojections_a.log"), context.StringValue("../resources/TestResultsWithProjections/multi/TestResultsWithProjections_A.log")},
		{context.BooleanValue(false), context.StringValue("testresultswithprojections_a.txt"), context.StringValue("../resources/TestResultsWithProjections/single/TestResultsWithProjections_A.txt")},
		{context.BooleanValue(false), context.StringValue("testresultswithprojections_b.log"), context.StringValue("../resources/TestResultsWithProjections/multi/TestResultsWithProjections_B.log")},
		{context.BooleanValue(false), context.StringValue("testresultswithprojections_c.txt"), context.StringValue("../resources/TestResultsWithProjections/multi/TestResultsWithProjections_C.txt")},
		{context.BooleanValue(false), context.StringValue("testresultswithprojections_d.txt"), context.StringValue("../resources/TestResultsWithProjections/multi/TestResultsWithProjections_D.txt")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsInNestedDirectoriesWithOptionToTraverseNestedDirectoriesAsFalse(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select isdir, lower(name), path from ../resources/TestResultsWithProjections/ order by 1 desc, 2 asc", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions().DisableNestedTraversal()).Execute()
	expected := [][]context.Value{
		{context.BooleanValue(true), context.StringValue("empty"), context.StringValue("../resources/TestResultsWithProjections/empty")},
		{context.BooleanValue(true), context.StringValue("multi"), context.StringValue("../resources/TestResultsWithProjections/multi")},
		{context.BooleanValue(true), context.StringValue("single"), context.StringValue("../resources/TestResultsWithProjections/single")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsInNestedDirectoriesWithOptionToIgnoreTraversalOfDirectories(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select isdir, lower(name), path from ../resources/TestResultsWithProjections/ order by 1 desc, 2 asc", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}

	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions().EnableNestedTraversal().DirectoriesToIgnoreTraversal([]string{"multi", "empty"})).Execute()
	expected := [][]context.Value{
		{context.BooleanValue(true), context.StringValue("empty"), context.StringValue("../resources/TestResultsWithProjections/empty")},
		{context.BooleanValue(true), context.StringValue("multi"), context.StringValue("../resources/TestResultsWithProjections/multi")},
		{context.BooleanValue(true), context.StringValue("single"), context.StringValue("../resources/TestResultsWithProjections/single")},
		{context.BooleanValue(false), context.StringValue("testresultswithprojections_a.txt"), context.StringValue("../resources/TestResultsWithProjections/single/TestResultsWithProjections_A.txt")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsInCaseInsensitiveManner(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("SELECT LOWER(NAME), BASE64(NAME) FROM ../resources/TestResultsWithProjections/single", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.txt"), context.StringValue("VGVzdFJlc3VsdHNXaXRoUHJvamVjdGlvbnNfQS50eHQ=")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjections3(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), ext from ../resources/TestResultsWithProjections/single", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.txt"), context.StringValue(".txt")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsWithConcatWsFunction(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), concatWs(lower(name), uid, gid, '#') from ../resources/TestResultsWithProjections/single", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	uid, gid := os.Getuid(), os.Getgid()

	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.txt"), context.StringValue(fmt.Sprintf("testresultswithprojections_a.txt#%v#%v", uid, gid))},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsWithSubstringFunction(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), substr(lower(name), 15) from ../resources/TestResultsWithProjections/single", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()

	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.txt"), context.StringValue("projections_a.txt")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsWithDayDifference(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select daydiff(now(), now()) from ../resources/TestResultsWithProjections/single", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()

	result := queryResults.atIndex(0).AllAttributes()[0]
	asFloat64, _ := result.GetNumericAsFloat64()

	if math.Round(asFloat64) != float64(0) {
		t.Fatalf("Expected day difference of 2 current times to be equal to zero but received %v and round resulted in %v", asFloat64, math.Round(asFloat64))
	}
}

func TestResultsWithProjectionsWithHourDifference(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select hourdiff(now(), now()) from ../resources/TestResultsWithProjections/single", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()

	result := queryResults.atIndex(0).AllAttributes()[0]
	asFloat64, _ := result.GetNumericAsFloat64()

	if math.Round(asFloat64) != float64(0) {
		t.Fatalf("Expected hour difference of 2 current times to be equal to zero but received %v and round resulted in %v", asFloat64, math.Round(asFloat64))
	}
}

func TestResultsWithProjectionsAndLimit1(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), ext from ../resources/TestResultsWithProjections/multi limit 3", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	if queryResults.Count() != 3 {
		t.Fatalf("Expected result Count to be %v, received %v", 3, queryResults.Count())
	}
}

func TestResultsWithProjectionsAndLimit2(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), ext from ../resources/TestResultsWithProjections/multi limit 0", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	if queryResults.Count() != 0 {
		t.Fatalf("Expected result Count to be %v, received %v", 3, queryResults.Count())
	}
}

func TestResultsWithProjectionsAndLimit3(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name) from ../resources/TestResultsWithProjections/multi order by 1 limit 2", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log")},
		{context.StringValue("testresultswithprojections_b.log")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsOrderBy1(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select upper(lower(name)), ext from ../resources/TestResultsWithProjections/multi order by 1 desc", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("TESTRESULTSWITHPROJECTIONS_D.TXT"), context.StringValue(".txt")},
		{context.StringValue("TESTRESULTSWITHPROJECTIONS_C.TXT"), context.StringValue(".txt")},
		{context.StringValue("TESTRESULTSWITHPROJECTIONS_B.LOG"), context.StringValue(".log")},
		{context.StringValue("TESTRESULTSWITHPROJECTIONS_A.LOG"), context.StringValue(".log")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsOrderBy2(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), ext from ../resources/TestResultsWithProjections/multi order by 2, 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log"), context.StringValue(".log")},
		{context.StringValue("testresultswithprojections_b.log"), context.StringValue(".log")},
		{context.StringValue("testresultswithprojections_c.txt"), context.StringValue(".txt")},
		{context.StringValue("testresultswithprojections_d.txt"), context.StringValue(".txt")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsIncludingConcatFunction(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), concat(lower(name), '-FILE') from ../resources/TestResultsWithProjections/multi order by 2, 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log"), context.StringValue("testresultswithprojections_a.log-FILE")},
		{context.StringValue("testresultswithprojections_b.log"), context.StringValue("testresultswithprojections_b.log-FILE")},
		{context.StringValue("testresultswithprojections_c.txt"), context.StringValue("testresultswithprojections_c.txt-FILE")},
		{context.StringValue("testresultswithprojections_d.txt"), context.StringValue("testresultswithprojections_d.txt-FILE")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsIncludingConcatWsFunction(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), concatws(lower(name), 'FILE', '@') from ../resources/TestResultsWithProjections/multi order by 2, 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log"), context.StringValue("testresultswithprojections_a.log@FILE")},
		{context.StringValue("testresultswithprojections_b.log"), context.StringValue("testresultswithprojections_b.log@FILE")},
		{context.StringValue("testresultswithprojections_c.txt"), context.StringValue("testresultswithprojections_c.txt@FILE")},
		{context.StringValue("testresultswithprojections_d.txt"), context.StringValue("testresultswithprojections_d.txt@FILE")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsIncludingContainsFunction(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), contains(lower(name), 'log') from ../resources/TestResultsWithProjections/multi order by 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log"), context.BooleanValue(true)},
		{context.StringValue("testresultswithprojections_b.log"), context.BooleanValue(true)},
		{context.StringValue("testresultswithprojections_c.txt"), context.BooleanValue(false)},
		{context.StringValue("testresultswithprojections_d.txt"), context.BooleanValue(false)},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsIncludingAddFunction(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), add(len(name), 4) from ../resources/TestResultsWithProjections/multi order by 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log"), context.Float64Value(36)},
		{context.StringValue("testresultswithprojections_b.log"), context.Float64Value(36)},
		{context.StringValue("testresultswithprojections_c.txt"), context.Float64Value(36)},
		{context.StringValue("testresultswithprojections_d.txt"), context.Float64Value(36)},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsIncludingSubtractFunction(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), sub(len(name), 2) from ../resources/TestResultsWithProjections/multi order by 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log"), context.Float64Value(30)},
		{context.StringValue("testresultswithprojections_b.log"), context.Float64Value(30)},
		{context.StringValue("testresultswithprojections_c.txt"), context.Float64Value(30)},
		{context.StringValue("testresultswithprojections_d.txt"), context.Float64Value(30)},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsIncludingMultiplyFunction(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), mul(len(name), 2) from ../resources/TestResultsWithProjections/multi order by 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log"), context.Float64Value(64)},
		{context.StringValue("testresultswithprojections_b.log"), context.Float64Value(64)},
		{context.StringValue("testresultswithprojections_c.txt"), context.Float64Value(64)},
		{context.StringValue("testresultswithprojections_d.txt"), context.Float64Value(64)},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsIncludingDivideFunction(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), div(len(name), 2) from ../resources/TestResultsWithProjections/multi order by 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log"), context.Float64Value(16)},
		{context.StringValue("testresultswithprojections_b.log"), context.Float64Value(16)},
		{context.StringValue("testresultswithprojections_c.txt"), context.Float64Value(16)},
		{context.StringValue("testresultswithprojections_d.txt"), context.Float64Value(16)},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsIncludingNegativeValueInAddSubMulDivFunction(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select add(len(name), -2), sub(len(name), -2), mul(len(name), -2), div(len(name), -2) from ../resources/TestResultsWithProjections/multi order by 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.Float64Value(30), context.Float64Value(34), context.Float64Value(-64), context.Float64Value(-16)},
		{context.Float64Value(30), context.Float64Value(34), context.Float64Value(-64), context.Float64Value(-16)},
		{context.Float64Value(30), context.Float64Value(34), context.Float64Value(-64), context.Float64Value(-16)},
		{context.Float64Value(30), context.Float64Value(34), context.Float64Value(-64), context.Float64Value(-16)},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsWithIdentity(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), identity(add(1,2)) from ../resources/TestResultsWithProjections/multi order by 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log"), context.Float64Value(3.0)},
		{context.StringValue("testresultswithprojections_b.log"), context.Float64Value(3.0)},
		{context.StringValue("testresultswithprojections_c.txt"), context.Float64Value(3.0)},
		{context.StringValue("testresultswithprojections_d.txt"), context.Float64Value(3.0)},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsWithBaseName(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), lower(basename) from ../resources/TestResultsWithProjections/multi order by 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log"), context.StringValue("testresultswithprojections_a")},
		{context.StringValue("testresultswithprojections_b.log"), context.StringValue("testresultswithprojections_b")},
		{context.StringValue("testresultswithprojections_c.txt"), context.StringValue("testresultswithprojections_c")},
		{context.StringValue("testresultswithprojections_d.txt"), context.StringValue("testresultswithprojections_d")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsWithFormatSize(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select lower(name), fmtsize from ../resources/TestResultsWithProjections/multi order by 1", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	queryResults, _ := NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	expected := [][]context.Value{
		{context.StringValue("testresultswithprojections_a.log"), context.StringValue("58 B")},
		{context.StringValue("testresultswithprojections_b.log"), context.StringValue("58 B")},
		{context.StringValue("testresultswithprojections_c.txt"), context.StringValue("58 B")},
		{context.StringValue("testresultswithprojections_d.txt"), context.StringValue("58 B")},
	}
	assertMatch(t, expected, queryResults)
}

func TestResultsWithProjectionsWithoutProperParametersToAFunction(t *testing.T) {
	newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
	aParser, err := parser.NewParser("select name, lower() from ../resources/TestResultsWithProjections/single", newContext)
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	selectQuery, err := aParser.Parse()
	if err != nil {
		t.Fatalf("error is %v", err)
	}
	_, err = NewSelectQueryExecutor(selectQuery, newContext, NewDefaultOptions()).Execute()
	if err == nil {
		t.Fatalf("Expected an error on running a query with lower() without any parameter")
	}
}

func assertMatch(t *testing.T, expected [][]context.Value, queryResults *EvaluatingRows, skipAttributeIndices ...int) {
	contains := func(slice []int, value int) bool {
		for _, v := range slice {
			if value == v {
				return true
			}
		}
		return false
	}
	if uint32(len(expected)) != queryResults.Count() {
		t.Fatalf("Expected length of the query results to be %v, received %v", len(expected), queryResults.Count())
	}
	for rowIndex, row := range expected {
		if len(row) != queryResults.atIndex(rowIndex).TotalAttributes() {
			t.Fatalf("Expected length of the rowAttributes in row index %v to be %v, received %v", rowIndex, len(row), queryResults.atIndex(rowIndex).TotalAttributes())
		}
		rowAttributes := queryResults.atIndex(rowIndex).AllAttributes()
		for attributeIndex, attributeValue := range row {
			if !contains(skipAttributeIndices, attributeIndex) && rowAttributes[attributeIndex].CompareTo(attributeValue) != 0 {
				t.Fatalf("Expected %v to match %v at row index %v, attribute index %v",
					attributeValue,
					rowAttributes[attributeIndex],
					rowIndex,
					attributeIndex,
				)
			}
		}
	}
}
