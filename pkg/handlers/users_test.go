package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ajaen4/go-standard-lib-api/internal/db"
)

var apiCfg *ApiConfig

func TestMain(m *testing.M) {
	db, err := db.NewDB("./database.json")
	if err != nil {
		log.Fatalf("Error initializing DB: %s", err)
	}

	apiCfg = &ApiConfig{
		DB:        db,
		JwtSecret: os.Getenv("JWT_SECRET"),
		PolkaKey:  os.Getenv("POLKA_KEY"),
	}

	code := m.Run()

	db.RemoveDB()
	os.Exit(code)
}

func TestPostUser(t *testing.T) {
	userReq := &UserReq{
		Email:    "test@email.com",
		Password: "testPassword",
	}
	jsonReq, err := json.Marshal(userReq)
	if err != nil {
		t.Fatal(err)
	}

	readerReq := strings.NewReader(string(jsonReq))
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/api/users", readerReq)
	if err != nil {
		t.Fatal(err)
	}

	err1 := apiCfg.PostUser(w, r)
	if err != nil {
		t.Fatal(err1)
	}

	response := UserResp{}
	expected := UserResp{
		Id:          1,
		Email:       userReq.Email,
		IsChirpyRed: false,
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response != expected {
		t.Errorf(
			"handler returned wrong response: got %v want %v",
			response,
			expected,
		)
	}
}
