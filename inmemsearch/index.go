package inmemsearch

import (
	"textscout/common"
)

// choosen values for bitmap array lengths
// generally this length depends on the max movie_id currently present since
// that is how we are calculating its index in the array then its bitmap value
const MAX_BIT_LEN int = 31
const MAX_BITMAP_LEN int = 22581

// inverted index to map each word to all the document IDs it occurs in
type IndexMap struct {
	DocFreq     int
	PostingList []int
}
type Index map[string]*IndexMap

func (idx Index) Add(docs []Document) {
	for _, doc := range docs {
		for _, token := range analyze(common.ConcatStrings(doc.MovieTitle, doc.Overview)) {
			indexMap, ok := idx[token]
			if !ok {
				// init Index for each new token
				idx[token] = &IndexMap{
					DocFreq:     1,
					PostingList: []int{doc.ID},
				}
				continue
			}

			// avoids adding the same ID twice if the word is repeated more than once in the same sentence.
			curIds := indexMap.PostingList
			if len(curIds) != 0 && curIds[len(curIds)-1] == doc.ID {
				// increment frequency
				indexMap.DocFreq++
				continue
			}

			// add the new docID and increment the frequency
			indexMap.PostingList = append(curIds, doc.ID)
			indexMap.DocFreq++
		}
	}

}

func (idx Index) SearchIntersection(query string) []int {
	docIDs := make([]int, 0)

	for _, token := range analyze(query) {
		// get the docIDs list from inverted index for each token
		// find the common IDs from all such list
		indexMap, ok := idx[token]
		if !ok {
			// token doesn't exist, do we just return or return the found docIDs
			continue
		}
		if len(docIDs) == 0 {
			// init
			docIDs = indexMap.PostingList
			continue
		}
		docIDs = idx.intersection(docIDs, indexMap.PostingList)
	}

	return docIDs
}

func (idx Index) WordFreq(word string) int {
	// return the word count/freq in the whole document
	freq, ok := idx[word]
	if !ok {
		return 0
	}
	return freq.DocFreq
}

// a=[0, 4] intersection b=[0,1] -> c=[0]
func (idx Index) intersection(seta, setb []int) []int {
	intersection := make([]int, 0)
	bitmap := make([]int, MAX_BITMAP_LEN)

	// compute the bitmap first
	for _, item := range seta {
		// index to know where should we put the integer
		// mask to calculate the actual mask
		// example : https://go.dev/play/p/Wv-9Wxn-4Lj
		index, mask := divmod(item, MAX_BIT_LEN)
		bitMask := int(1 << mask)
		bitmap[index] |= bitMask
	}

	// calculate the intersection
	for _, item := range setb {
		index, mask := divmod(item, MAX_BIT_LEN)
		source := int(1 << mask)
		target := bitmap[index]

		if source&target != 0 {
			intersection = append(intersection, item)
		}
	}

	return intersection

}

func (idx Index) SearchUnion(query string) []int {
	docIDs := make([]int, 0)

	for _, token := range analyze(query) {
		// get the docIDs list from inverted index for each token
		// find the common IDs from all such list
		indexMap, ok := idx[token]
		if !ok {
			// token doesn't exist, do we just return or return the found docIDs
			continue
		}
		if len(docIDs) == 0 {
			// init
			docIDs = indexMap.PostingList
			continue
		}
		docIDs = idx.union(docIDs, indexMap.PostingList)
	}

	return docIDs
}

// a=[0, 4] union b=[0,1] -> c=[0, 1, 4]
func (idx Index) union(seta, setb []int) []int {
	union := make([]int, 0)
	bitmap := make([]int, MAX_BITMAP_LEN)

	// compute the bitmap first
	for _, item := range seta {
		index, mask := divmod(item, MAX_BIT_LEN)
		bitMask := int(1 << mask)
		bitmap[index] |= bitMask
		union = append(union, item)
	}

	// calculate the union
	for _, item := range setb {
		index, mask := divmod(item, MAX_BIT_LEN)
		source := int(1 << mask)
		target := bitmap[index]

		// dont add the same id twice
		if source&target == 0 {
			union = append(union, item)
		}
	}

	return union
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}
