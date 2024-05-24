package textsearch

type InMemSearch struct {
	idx Index
}

func prepareIndex(filePath string) (Index, error) {
	// build the inverted index by reading the json from this filepath
	docs, err := loadMovies(filePath)
	if err != nil {
		return nil, err
	}

	// create the in-memory inverted index
	index := make(Index)
	index.Add(docs)
	return index, nil
}
func GetInMemSearch(filePath string) *InMemSearch {
	index, err := prepareIndex(filePath)
	if err != nil {
		panic(err.Error())
	}
	return &InMemSearch{
		idx: index,
	}
}

func (im *InMemSearch) Search(query string) {
	// search the index given some query
	// query being the movie_title or movie_description
	im.idx.Search(query)

}
