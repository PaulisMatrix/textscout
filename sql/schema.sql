-- Define the movies table

CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    adult BOOLEAN, 
    backdrop_path VARCHAR(100), 
    genre_ids INTEGER[],
    movie_id INTEGER NOT NULL,
    movie_language VARCHAR(100),
    movie_original_title TEXT,
    movie_overview TEXT,
    popularity DOUBLE PRECISION,
    poster_path VARCHAR(100), 
    release_date VARCHAR(100),
    movie_title TEXT NOT NULL,
    video BOOLEAN,
    vote_average DOUBLE PRECISION,
    vote_count BIGINT
);

INSERT INTO movies(adult, backdrop_path, genre_ids, movie_id, movie_language, movie_original_title, movie_overview,
popularity, poster_path, release_date, movie_title, video, vote_average, vote_count) values(
    FALSE, '/j3Z3XktmWB1VhsS8iXNcrR86PXi.jpg', '{878, 28, 12}', 823464, 'en', 'Godzilla x Kong: The New Empire',
    'Following their explosive showdown, Godzilla and Kong must reunite against a colossal undiscovered threat hidden within our world, challenging their very existence – and our own.',
    7832.06, '/v4uvGFAkKuYfyKLGZnYj6l47ERQ.jpg', '2024-03-27', 'Godzilla x Kong: The New Empire', FALSE, 7.249, 1920
    );


select * from movies where lower(movie_title) like lower('%Kong%');