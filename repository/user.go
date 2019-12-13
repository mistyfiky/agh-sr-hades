package repository

import "database/sql"

type User struct {
	Username string
	Password string
}

func FindUserByUsername(username string) *User {
	var result User
	row := db.QueryRow("SELECT `username`, `password` FROM `users` WHERE `username` = ?", username)
	switch err := row.Scan(&result.Username, &result.Password); err {
	case nil:
		return &result
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}
