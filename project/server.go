package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerServer struct {
	store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	score := p.store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	// Fprintf accepts only string
	fmt.Fprintf(w, "%v", score)
}

type PlayerStore interface {
	GetPlayerScore(name string) int
}
