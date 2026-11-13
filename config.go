package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	DayStart string `json:"day_start"` // "HH:MM", e.g. "04:00"
}

func configPath() (string, error) {
	dir, err := jotDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func loadConfig() Config {
	cfg := Config{DayStart: "00:00"}
	path, err := configPath()
	if err != nil {
		return cfg
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, &cfg)
	if cfg.DayStart == "" {
		cfg.DayStart = "00:00"
	}
	return cfg
}

// this turns "HH:MM" into a duration since midnight.
func parseDayStart(s string) (time.Duration, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, fmt.Errorf("want HH:MM, got %q", s)
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute, nil
}

func formatClock(d time.Duration) string {
	m := int(d.Minutes())
	return fmt.Sprintf("%02d:%02d", m/60, m%60)
}
