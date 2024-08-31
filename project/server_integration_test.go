package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := InMemoryPlayerStore{}
	server := PlayerServer{&store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinReq(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinReq(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinReq(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreReq(player))
	assertStatus(t, response.Code, http.StatusOK)

	assertResBody(t, response.Body.String(), "3")
}
