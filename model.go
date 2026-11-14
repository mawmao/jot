package main

import (
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type mode int

const (
	modeBrowsing mode = iota
	modeAdding
	modeEditing
	modeConfirmingDelete
	modeOverview
)

type Todo struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type model struct {
	store          Store
	day            time.Time
	dayStart       time.Duration
	cursor         int
	input          textinput.Model
	mode           mode
	showHelp       bool
	overviewWeek   time.Time
	overviewCursor int
}

func initialModel(dayStart time.Duration) model {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.Placeholder = "Jot something down..."
	ti.CharLimit = 200

	return model{
		store:    loadStore(),
		day:      logicalDay(time.Now(), dayStart),
		dayStart: dayStart,
		input:    ti,
		showHelp: false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) current() []Todo {
	return m.store[dayKey(m.day)]
}

func (m *model) setCurrent(todos []Todo) {
	key := dayKey(m.day)
	if len(todos) == 0 {
		delete(m.store, key)
		return
	}
	m.store[key] = todos
}
