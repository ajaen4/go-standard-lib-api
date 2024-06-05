package handlers

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ajaen4/go-standard-lib-api/pkg/api_errors"
)

func TestProcessWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"no words are processed", "Hello new world", "Hello new world"},
		{"some words are processed", "kerfuffle new world", "**** new world"},
		{"all words are processed", "kerfuffle sharbert fornax", "**** **** ****"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessWords(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessWords(%s) = %s; want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPostChirpReq_validate(t *testing.T) {

	tests := []struct {
		name        string
		method      string
		path        string
		payload     io.Reader
		expectedRes PostChirpReq
		expectedErr *api_errors.ClientErr
	}{
		{
			"empty request",
			"POST",
			"/api/chirps",
			strings.NewReader(`{}`),
			PostChirpReq{},
			&api_errors.ClientErr{
				HttpCode: http.StatusBadRequest,
				Message:  "Invalid body parameters",
				Errors: map[string]string{
					"body": "invalid body",
				},
			},
		},
		{
			"empty Body",
			"POST",
			"/api/chirps",
			strings.NewReader(`{"Body": ""}`),
			PostChirpReq{},
			&api_errors.ClientErr{
				HttpCode: http.StatusBadRequest,
				Message:  "Invalid body parameters",
				Errors: map[string]string{
					"body": "invalid body",
				},
			},
		},
		{
			"correct Body",
			"POST",
			"/api/chirps",
			strings.NewReader(`{"Body": "correct body"}`),
			PostChirpReq{
				Body: "correct body",
			},
			nil,
		},
		{
			"Body too long",
			"POST",
			"/api/chirps",
			strings.NewReader(`{"Body": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam."}`),
			PostChirpReq{},
			&api_errors.ClientErr{
				HttpCode: http.StatusBadRequest,
				Message:  "Invalid body parameters",
				Errors: map[string]string{
					"body": "invalid body",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, tt.payload)
			if err != nil {
				t.Fatal(err)
			}

			chirpReq := PostChirpReq{}
			resultErr := chirpReq.validate(req)
			if !compareErrors(resultErr, tt.expectedErr) {
				t.Errorf("Error returned, got %v want %v", *resultErr, *tt.expectedErr)
			}

			if tt.expectedErr == nil && chirpReq != tt.expectedRes {
				t.Errorf("Got %v want %v", chirpReq, tt.expectedRes)
			}
		})
	}
}

func compareErrors(err1, err2 *api_errors.ClientErr) bool {
	if err1 == nil && err2 == nil {
		return true
	}
	if err1 == nil && err2 != nil || err1 != nil && err2 == nil {
		return false
	}

	if err1.HttpCode != err2.HttpCode || err1.Message != err2.Message {
		return false
	}

	for k1, v1 := range err1.Errors {
		v2, ok := err2.Errors[k1]
		if !ok {
			return false
		}
		if v1 != v2 {
			return false
		}
	}
	return true
}
