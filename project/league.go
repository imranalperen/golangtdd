package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem at parsing league, %v", err)
	}

	return league, err
}
