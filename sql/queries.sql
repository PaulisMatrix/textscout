-- name: AddMovie :exec
INSERT INTO movies(
    adult, backdrop_path, genre_ids, movie_id, movie_language, movie_original_title, movie_overview,
    popularity, poster_path, release_date, movie_title, video, vote_average, vote_count) 
    VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
);

-- name: GetMovieByTitle :many
SELECT * FROM movies WHERE LOWER(movie_title) LIKE LOWER('%' || $1 || '%') LIMIT 5;

-- name: GetMovieByDesc :many
SELECT * FROM movies WHERE LOWER(movie_overview) LIKE LOWER('%' || $1 || '%') LIMIT 5;

-- name: GetMovieByTitleAndDesc :many
SELECT * FROM movies WHERE LOWER(movie_title) LIKE LOWER('%' || $1 || '%') AND LOWER(movie_overview) LIKE LOWER('%' || $2 || '%') LIMIT 5;