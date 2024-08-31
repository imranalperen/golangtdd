package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "integration player"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinReq(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinReq(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinReq(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreReq(player))
		assertStatus(t, response.Code, http.StatusOK)

		assertResBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		res := httptest.NewRecorder()
		server.ServeHTTP(res, newLeagueReq())
		assertStatus(t, res.Code, http.StatusOK)

		got := getLeagueFromRes(t, res.Body)
		want := []Player{
			{"integration player", 3},
		}
		assertLeague(t, got, want)
	})
}
