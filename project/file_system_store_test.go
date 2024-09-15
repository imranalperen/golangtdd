package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		db := strings.NewReader(`[
		{"Name": "alperen", "Wins": 10},
		{"Name": "imran", "Wins": 33}]`)

		store := FileSystemPlayerStore{db}

		got := store.GetLeague()
		want := []Player{
			{"alperen", 10},
			{"imran", 33},
		}

		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		db := strings.NewReader(`[
		{"Name": "alperen", "Wins": 10},
		{"Name": "imran", "Wins": 33}]`)

		store := FileSystemPlayerStore{db}

		got := store.GetPlayerScore("alperen")
		want := 10

		assertScoreEquals(t, got, want)
	})
}

// ates a temporary file for us to use. The "db" value we've pas
func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
