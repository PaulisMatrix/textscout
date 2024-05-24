# circuitsearch
search API over in-memory movies dataset


# Approach1: (Using database for searching)

* Store the records in a database. 
* Query the db: `select * from movies where movie_name like "%'name'%" and movie_description like "%'desc'%"`
* Furthermore to optimise this query, we can have a `secondary index` on `movie_name` or `movie_description` depending on the read query patterns the server receives.



# API structure:

* Request: Searchable by both title and desc or either field. lowercase matching id one so its case insensitive.
    * GET request: `curl -i --location 'http://localhost:8080/api/v1/search?title=kong&desc=godzilla'`
    * query params: `title` and `desc`

* Response: List of all first 5 matched movies.
    * Example:
        ```{
            "movies": [
                {
                    "backdrop_path": "/j3Z3XktmWB1VhsS8iXNcrR86PXi.jpg",
                    "genre_ids": [
                        878,
                        28,
                        12
                    ],
                    "id": 1,
                    "original_language": "en",
                    "original_title": "Godzilla x Kong: The New Empire",
                    "overview": "Following their explosive showdown, Godzilla and Kong must reunite against a colossal undiscovered threat hidden within our world, challenging their very existence – and our own.",
                    "popularity": 7832.06,
                    "poster_path": "/v4uvGFAkKuYfyKLGZnYj6l47ERQ.jpg",
                    "release_date": "2024-03-27",
                    "title": "Godzilla x Kong: The New Empire",
                    "vote_average": 7.249,
                    "vote_count": 1920
                }
            ]
        }```

# Approach2: (Using an in-memory Inverted Index for searching)

* Combine `movie_name` and `movie_description` and build the inverted index using this whole text. This is fine since its safe to say both are related since they are talking about the same movie.
* This content is passed through tokenizing + normalising + stopWordsRemoval + stemming pipeline to generate the final keywords/tokens.
* Inverted index: `keyword: []MovieIDs`. Its a map of keyword and value being list of all movies ids containing that keyword.
* `index.Add([]MovieData)` : builds the index. 
* `index.Search(query)` : searches the index and returns the movie ids containing the query keywords. That's the final result.
* Small catch: Works only for english words since other languages contexts/meanings change.


# 
