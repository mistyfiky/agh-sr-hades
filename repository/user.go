package repository

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func init() {
	defer func() {
		if result := recover(); nil != result {
			err, _ := result.(error)
			log.Fatal(err)
		}
	}()
	drop := "DROP TABLE IF EXISTS `users`"
	if _, err := db.Exec(drop); nil != err {
		panic(err)
	}
	create := "CREATE TABLE `users` (" +
		"`username` VARCHAR(255) NOT NULL, " +
		"`password` VARCHAR(255) NOT NULL, " +
		"PRIMARY KEY (`username`))"
	if _, err := db.Exec(create); nil != err {
		panic(err)
	}
	SaveUser(NewUser("user", "pass"))
}

type user struct {
	username string
	password string
}

func NewUser(username string, plainPassword string) *user {
	return &user{username: username, password: hashPassword(plainPassword)}
}

func (user *user) GetUsername() string {
	return user.username
}

func (user *user) GetPassword() string {
	return user.password
}

func hashPassword(plainPassword string) string {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	if nil != err {
		panic(err)
	}
	return string(password[:])
}

func (user *user) IsAuthenticatedBy(plainPassword string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(user.password), []byte(plainPassword))
}

func FindUserByUsername(username string) *user {
	var result user
	query := "SELECT `username`, `password` FROM `users` WHERE `username` = ?"
	row := db.QueryRow(query, username)
	switch err := row.Scan(&result.username, &result.password); err {
	case nil:
		return &result
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

func SaveUser(user *user) {
	query := "INSERT INTO `users` (`username`, `password`) VALUES (?, ?)"
	if _, err := db.Exec(query, user.username, user.password); nil != err {
		panic(err)
	}
}
