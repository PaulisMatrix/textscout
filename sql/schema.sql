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