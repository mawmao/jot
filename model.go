package main

import (
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"fmt"
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
	statusErr      string
}

func initialModel(dayStart time.Duration, status string) model {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.Placeholder = "Jot something down..."
	ti.CharLimit = 200

	return model{
		store:     loadStore(),
		day:       logicalDay(time.Now(), dayStart),
		dayStart:  dayStart,
		input:     ti,
		showHelp:  false,
		statusErr: status,
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

func (m *model) clampCursor() {
	n := len(m.current())
	if n == 0 || m.cursor < 0 {
		m.cursor = 0
		return
	}
	if m.cursor >= n {
		m.cursor = n - 1
	}
}

func (m *model) persist() {
	if err := saveStore(m.store); err != nil {
		m.statusErr = fmt.Sprintf("save failed: %v", err)
		return
	}
	m.statusErr = ""
}

func (m *model) startAdding() tea.Cmd {
	m.mode = modeAdding
	m.input.Prompt = "> "
	m.input.SetValue("")
	return m.input.Focus()
}

func (m *model) startEditing() tea.Cmd {
	todos := m.current()
	if len(todos) == 0 {
		return nil
	}
	m.clampCursor()
	m.mode = modeEditing
	m.input.Prompt = "edit> "
	m.input.SetValue(todos[m.cursor].Title)
	m.input.CursorEnd()
	return m.input.Focus()
}
