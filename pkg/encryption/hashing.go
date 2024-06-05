package encryption

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func Hash(str string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), 4)
	if err != nil {
		return []byte{}, errors.New("internal error")
	}
	return hash, nil
}
