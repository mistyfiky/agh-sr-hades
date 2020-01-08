package repository

import (
	"database/sql"
	"log"
)

func init() {
	defer func() {
		if result := recover(); nil != result {
			err, _ := result.(error)
			log.Fatal(err)
		}
	}()
	drop := "DROP TABLE IF EXISTS `users_movies`"
	if _, err := db.Exec(drop); nil != err {
		panic(err)
	}
	create := "CREATE TABLE `users_movies` (" +
		"`username` VARCHAR(255) NOT NULL, " +
		"`movie_id` VARCHAR(255) NOT NULL, " +
		"PRIMARY KEY (`username`, `movie_id`))"
	if _, err := db.Exec(create); nil != err {
		panic(err)
	}
}

type userMovie struct {
	username string
	movieId  string
}

func NewUserMovie(username, movieId string) *userMovie {
	return &userMovie{username: username, movieId: movieId}
}

func (userMovie *userMovie) GetUsername() string {
	return userMovie.username
}

func (userMovie *userMovie) GetMovieId() string {
	return userMovie.movieId
}

func FindUserMoviesByUsername(username string) []userMovie {
	var result []userMovie
	query := "SELECT `username`, `movie_id` FROM `users_movies` WHERE `username` = ?"
	rows, err := db.Query(query, username)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var tmp userMovie
		if err := rows.Scan(&tmp.username, &tmp.movieId); err != nil {
			panic(err)
		}
		result = append(result, tmp)
	}
	return result
}

func FindUserMovieByUsernameAndMovieId(username, movieId string) *userMovie {
	var result userMovie
	query := "SELECT `username`, `movie_id` FROM `users_movies` WHERE `username` = ? AND  `movie_id` = ?"
	row := db.QueryRow(query, username, movieId)
	switch err := row.Scan(&result.username, &result.movieId); err {
	case nil:
		return &result
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

func SaveUserMovie(userMovie *userMovie) {
	query := "INSERT INTO `users_movies` (`username`, `movie_id`) VALUES (?, ?)"
	if _, err := db.Exec(query, userMovie.username, userMovie.movieId); nil != err {
		panic(err)
	}
}
