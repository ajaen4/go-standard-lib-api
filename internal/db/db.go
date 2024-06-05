package db

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/ajaen4/go-standard-lib-api/pkg/api_errors"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	PssHash     []byte `json:"pss_hash"`
	RefToken    string `json:"refresh_token"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

var ErrChirpNotFound = api_errors.ClientErr{
	HttpCode: http.StatusBadRequest,
	Message:  "chirp id not found",
}
var ErrUserAlrExist = api_errors.ClientErr{
	HttpCode: http.StatusBadRequest,
	Message:  "user already exists",
}
var ErrUserNotExist = api_errors.ClientErr{
	HttpCode: http.StatusBadRequest,
	Message:  "user doesn't exist",
}
var ErrIncorrectPss = api_errors.ClientErr{
	HttpCode: http.StatusBadRequest,
	Message:  "incorrect password",
}
var ErrIncorrectChirpId = api_errors.ClientErr{
	HttpCode: http.StatusBadRequest,
	Message:  "incorrect chirp id",
}
var ErrIncorrectAuthorId = api_errors.ClientErr{
	HttpCode: http.StatusBadRequest,
	Message:  "incorrect author id",
}

func NewDB(path string) (*DB, error) {
	isDebug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	log.Println("debug:", *isDebug)

	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	if *isDebug {
		err := db.RemoveDB()
		if err != nil {
			return &DB{}, err
		}
	}

	err := db.ensureDB()
	if err != nil {
		return &DB{}, err
	}

	return db, nil
}

func (db *DB) RemoveDB() error {
	if _, err := os.Stat(db.path); err == nil {
		errR := os.Remove(db.path)
		if errR != nil {
			return errors.New(fmt.Sprintf("Failed to delete database: %v", errR))
		}
	}
	return nil
}

func (db *DB) ensureDB() error {
	_, errS := os.Stat(db.path)
	if errS != nil && os.IsNotExist(errS) {
		errW := db.writeDB(DBStructure{
			Chirps: map[int]Chirp{},
			Users:  map[int]User{},
		})
		if errW != nil {
			return errW
		}
	} else if errS != nil {
		return errS
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	fileContent, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	var chirpsById DBStructure
	err = json.Unmarshal(fileContent, &chirpsById)
	if err != nil {
		return DBStructure{}, err
	}
	return chirpsById, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	jsonContent, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	errW := os.WriteFile(db.path, jsonContent, 0644)
	if errW != nil {
		return errW
	}
	return nil
}
