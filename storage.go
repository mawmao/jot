package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const dateFormat = "2006-01-02"

type Store map[string][]Todo

func jotDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "jot")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func storePath() (string, error) {
	dir, err := jotDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "store.json"), nil
}

func loadStore() Store {
	path, err := storePath()
	if err != nil {
		return Store{}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Store{}
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil || s == nil {
		return Store{}
	}
	return s
}

func saveStore(s Store) error {
	path, err := storePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// logicalDay returns the calendar date that t belongs to, given a
// day boundary at dayStart. With dayStart = 4h, 03:00 belongs to the
// previous date; 04:00 begins the new one.
func logicalDay(t time.Time, dayStart time.Duration) time.Time {

	shifted := t.Add(-dayStart)
	y, m, d := shifted.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func startOfWeek(t time.Time) time.Time {
	wd := int(t.Weekday())
	if wd == 0 { // Sunday
		wd = 7
	}
	return t.AddDate(0, 0, -(wd - 1))
}

func dayIndexInWeek(day, weekStart time.Time) int {
	dy, dm, dd := day.Date()
	wy, wm, wd := weekStart.Date()
	d := time.Date(dy, dm, dd, 0, 0, 0, 0, time.UTC)
	w := time.Date(wy, wm, wd, 0, 0, 0, 0, time.UTC)
	return int(d.Sub(w).Hours() / 24)
}

func dayKey(t time.Time) string {
	return t.Format(dateFormat)
}
