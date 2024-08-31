package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"imran":   20,
			"alperen": 10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)
	t.Run("returns imran's score", func(t *testing.T) {
		req := newGetScoreReq("imran")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusOK)
		assertResBody(t, res.Body.String(), "20")
	})

	t.Run("returns alperen's score", func(t *testing.T) {
		req := newGetScoreReq("alperen")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusOK)
		assertResBody(t, res.Body.String(), "10")
	})

	t.Run("return 404 on missing player", func(t *testing.T) {
		req := newGetScoreReq("zebberi")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("it returns accepted on POST", func(t *testing.T) {
		player := "imran"
		req := newPostWinReq(player)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecoirWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league tabla as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Lacriment", 20},
			{"Foo", 25},
			{"Bar", 30},
		}
		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		req := newLeagueReq()
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		got := getLeagueFromRes(t, res.Body)
		assertStatus(t, res.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, res, "application/json")
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func newGetScoreReq(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newPostWinReq(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}

func getLeagueFromRes(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}

func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func newLeagueReq() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}
