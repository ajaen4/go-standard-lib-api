package db

import (
	"github.com/ajaen4/go-standard-lib-api/pkg/api_errors"
	"github.com/ajaen4/go-standard-lib-api/pkg/encryption"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email string, pss string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			err := ErrUserAlrExist
			return User{}, &err
		}
	}

	id := len(dbStructure.Users) + 1
	pssHash, err := bcrypt.GenerateFromPassword([]byte(pss), 4)
	if err != nil {
		return User{}, err
	}

	newUser := User{
		Email:   email,
		PssHash: pssHash,
		Id:      id,
	}
	dbStructure.Users[id] = newUser
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (db *DB) UpdateUser(id int, newEmail string, newPss string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newPssHash, err := encryption.Hash(newPss)
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		err := ErrUserNotExist
		return User{}, &err
	}

	user.Email = newEmail
	user.PssHash = newPssHash
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) Login(email string, pss string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword(user.PssHash, []byte(pss))
			if err == nil {
				return user, nil
			} else {
				incPss := ErrIncorrectPss
				return User{}, &incPss
			}
		}
	}
	userNotExist := ErrUserNotExist
	return User{}, &userNotExist
}

func (db *DB) SaveRefToken(id int, refreshToken string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		err := ErrUserNotExist
		return &err
	}
	user.RefToken = refreshToken
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ValidateRefToken(refreshToken string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	var user User
	for _, User := range dbStructure.Users {
		if User.RefToken == refreshToken {
			user = User
			break
		}
	}
	if user.Id == 0 {
		err := ErrUserNotExist
		return User{}, &err
	}
	if user.RefToken != refreshToken {
		err := api_errors.UnauthErr
		err.LogMess = "Incorrect refresh token"
		return User{}, &err
	}

	return user, nil
}

func (db *DB) RevokeRefToken(refreshToken string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	for _, user := range dbStructure.Users {
		if user.RefToken == refreshToken {
			user.RefToken = ""
			dbStructure.Users[user.Id] = user
			err := db.writeDB(dbStructure)
			if err != nil {
				return err
			}
			return nil
		}
	}

	userNotExist := ErrUserNotExist
	return &userNotExist
}

func (db *DB) UserChirpyRed(userId int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[userId]
	if !ok {
		err := ErrUserNotExist
		return &err
	}

	user.IsChirpyRed = true
	dbStructure.Users[userId] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
