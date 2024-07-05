package inmemsearch

import (
	"testing"
	"textscout/common"
)

func TestTextSearchANDOperation(t *testing.T) {
	filePath := "/Users/rushiyadwade/Documents/go_dir/source/textscout/sample.json"
	inMemIdx := GetInMemSearch(filePath)
	title := "Kong"
	desc := "Godzilla"

	matchedDocs := inMemIdx.Intersection(common.ConcatStrings(title, desc))
	expectedLength := 1
	if len(matchedDocs) != expectedLength {
		t.Errorf("expectedLength is: %d, got: %d", expectedLength, len(matchedDocs))
	}

}

func TestTextSearchOROperation(t *testing.T) {
	filePath := "/Users/rushiyadwade/Documents/go_dir/source/textscout/sample.json"
	inMemIdx := GetInMemSearch(filePath)
	title := "Kong"
	desc := "Godzilla"

	matchedDocs := inMemIdx.Union(common.ConcatStrings(title, desc))
	expectedLength := 3
	if len(matchedDocs) != expectedLength {
		t.Errorf("expectedLength is: %d, got: %d", expectedLength, len(matchedDocs))
	}

}
