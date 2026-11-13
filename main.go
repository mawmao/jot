package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	dayStartFlag := flag.String("day-start", "", "day boundary as HH:MM (overrides config)")
	flag.Parse()

	cfg := loadConfig()
	if *dayStartFlag != "" {
		cfg.DayStart = *dayStartFlag
	}
	dayStart, err := parseDayStart(cfg.DayStart)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid day-start: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(dayStart))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
