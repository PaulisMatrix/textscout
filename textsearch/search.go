package textsearch

type InMemSearch struct {
	idx       Index
	movieDocs []Document
}

func prepareIndex(filePath string) (Index, []Document, error) {
	// build the inverted index by reading the json from this filepath
	docs, err := loadMovies(filePath)
	if err != nil {
		return nil, []Document{}, err
	}

	// create the in-memory inverted index
	index := make(Index)
	index.Add(docs)
	return index, docs, nil
}

func GetInMemSearch(filePath string) *InMemSearch {
	index, mdocs, err := prepareIndex(filePath)
	if err != nil {
		panic(err.Error())
	}
	return &InMemSearch{
		idx:       index,
		movieDocs: mdocs,
	}
}

func (im *InMemSearch) Search(query string) []Document {
	// search the index given some query
	// query being the movie_title or movie_description
	docIDs := im.idx.Search(query)
	docs := make([]Document, 0)

	for _, id := range docIDs {
		md := im.movieDocs[id]
		docs = append(docs, md)
	}

	return docs
}
