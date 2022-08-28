package source

import (
	"goselect/parser/tokenizer"
	"os/user"
	"testing"
)

func TestCreatesANewSourceFromCurrentDirectory1(t *testing.T) {
	tokens := tokenizer.NewEmptyTokens()
	tokens.Add(tokenizer.NewToken(tokenizer.RawString, "."))

	source, _ := NewSource(tokens.Iterator())
	if source.Directory != "." {
		t.Fatalf("Expected Directory path to be %v, received %v", ".", source.Directory)
	}
}

func TestCreatesANewSourceFromCurrentDirectory2(t *testing.T) {
	tokens := tokenizer.NewEmptyTokens()
	tokens.Add(tokenizer.NewToken(tokenizer.RawString, "."))
	tokens.Add(tokenizer.NewToken(tokenizer.RawString, "where"))

	source, _ := NewSource(tokens.Iterator())
	if source.Directory != "." {
		t.Fatalf("Expected Directory path to be %v, received %v", ".", source.Directory)
	}
}

func TestCreatesANewSourceWithHomeDirectorySymbol1(t *testing.T) {
	tokens := tokenizer.NewEmptyTokens()
	tokens.Add(tokenizer.NewToken(tokenizer.RawString, "~"))

	source, _ := NewSource(tokens.Iterator())
	expectedPath := homeDirectory()

	if source.Directory != expectedPath {
		t.Fatalf("Expected Directory path to be %v, received %v", expectedPath, source.Directory)
	}
}

func TestCreatesANewSourceWithHomeDirectorySymbol2(t *testing.T) {
	tokens := tokenizer.NewEmptyTokens()
	tokens.Add(tokenizer.NewToken(tokenizer.RawString, "~/apps"))

	source, _ := NewSource(tokens.Iterator())
	expectedPath := homeDirectory() + "/apps"

	if source.Directory != expectedPath {
		t.Fatalf("Expected Directory path to be %v, received %v", expectedPath, source.Directory)
	}
}

func TestThrowsAnErrorWithoutAnyTokens(t *testing.T) {
	tokens := tokenizer.NewEmptyTokens()
	_, err := NewSource(tokens.Iterator())

	if err == nil {
		t.Fatalf("Expected error to be non-nil when creating a source without any tokens")
	}
}

func homeDirectory() string {
	currentUser, err := user.Current()
	if err == nil {
		return currentUser.HomeDir
	}
	return ""
}
