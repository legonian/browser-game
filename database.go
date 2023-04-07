package main

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	usersLock sync.RWMutex
	users     map[string]*User
}

type User struct {
	Username string
	Hash     []byte
}

var ErrUserExistsValid = errors.New("user already exists and hash is matched")
var ErrUserExistsInvalid = errors.New("user already exists and hash is not matched")

func NewDB() (*DB, error) {
	return &DB{
		users: make(map[string]*User),
	}, nil
}

// CreateUser will create new user in database. If user already exists then hash
// will be checked and returned different error.
func (db *DB) CreateUser(username, password string) error {
	db.usersLock.RLock()
	user, exists := db.users[username]
	db.usersLock.RUnlock()

	if exists {
		err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
		if err != nil {
			return ErrUserExistsInvalid
		}

		return ErrUserExistsValid
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return ErrUserExistsInvalid
	}

	db.usersLock.Lock()
	db.users[username] = &User{
		Username: username,
		Hash:     hash,
	}
	db.usersLock.Unlock()

	return nil
}
