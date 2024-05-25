package textsearch

import (
	"circuithouse/common"
	"testing"
)

func TestTextSearch(t *testing.T) {
	filePath := "/Users/rushiyadwade/Documents/go_dir/source/circuitsearch/sample.json"
	inMemIdx := GetInMemSearch(filePath)
	title := "The Dune"
	desc := ""

	matchedDocs := inMemIdx.Search(common.ConcatStrings(title, desc))
	expectedLength := 2
	if len(matchedDocs) != expectedLength {
		t.Errorf("expectedLength is: %d, got: %d", expectedLength, len(matchedDocs))
	}

}
