package db

import "slices"

func (db *DB) GetChirps(sortBy string) ([]Chirp, error) {
	chirpsById, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := []Chirp{}
	for _, chirp := range chirpsById.Chirps {
		chirps = append(chirps, chirp)
	}

	if sortBy == "desc" {
		slices.Reverse(chirps)
	}

	return chirps, nil
}

func (db *DB) GetChirpsByAuthId(authorId int, sortBy string) ([]Chirp, error) {
	chirpsById, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := []Chirp{}
	for _, chirp := range chirpsById.Chirps {
		if chirp.AuthorId == authorId {
			chirps = append(chirps, chirp)
		}
	}

	if sortBy == "desc" {
		slices.Reverse(chirps)
	}

	return chirps, nil
}

func (db *DB) GetChirp(chirpID int) (Chirp, error) {
	chirpsById, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	for dbChirpID, dbChirp := range chirpsById.Chirps {
		if dbChirpID == chirpID {
			return dbChirp, nil
		}
	}
	chirpNotFound := ErrChirpNotFound
	return Chirp{}, &chirpNotFound
}

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		Body:     body,
		Id:       id,
		AuthorId: authorId,
	}
	dbStructure.Chirps[id] = newChirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) DeleteChirp(userId int, chirpId int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	if chirpId != len(dbStructure.Chirps) {
		err := ErrIncorrectChirpId
		return &err
	}

	if dbStructure.Chirps[chirpId].AuthorId != userId {
		err := ErrIncorrectAuthorId
		return &err
	}

	delete(dbStructure.Chirps, chirpId)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
