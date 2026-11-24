package main

import (
	"sort"
	"strings"
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
	modeConfirmingMove
	modeOverview
)

type sortMode int

const (
	sortManual sortMode = iota
	sortUndoneFirst
	sortDoneFirst
	sortAlphaAsc
	sortAlphaDesc
	sortModeCount
)

func (s sortMode) label() string {
	switch s {
	case sortUndoneFirst:
		return "undone"
	case sortDoneFirst:
		return "done"
	case sortAlphaAsc:
		return "a-z"
	case sortAlphaDesc:
		return "z-a"
	default:
		return ""
	}
}

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
	noDay          bool
	showHelp       bool
	overviewWeek   time.Time
	overviewCursor int
	width          int
	height         int
	sortMode       sortMode
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
	return m.store[m.activeKey()]
}

func (m *model) setCurrent(todos []Todo) {
	key := m.activeKey()
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

func (m *model) moveItem(delta int) {
	todos := m.current()
	if len(todos) < 2 {
		return
	}
	m.clampCursor()
	target := m.cursor + delta
	if target < 0 || target >= len(todos) {
		return
	}
	todos[m.cursor], todos[target] = todos[target], todos[m.cursor]
	m.cursor = target
	m.sortMode = sortManual
	m.persist()
}

func (m *model) cycleSort() {
	m.sortMode = (m.sortMode + 1) % sortModeCount
	if m.sortMode == sortManual {
		return
	}

	todos := m.current()
	if len(todos) < 2 {
		return
	}
	m.clampCursor()

	type indexedTodo struct {
		todo Todo
		orig int
	}
	wrapped := make([]indexedTodo, len(todos))
	for i, t := range todos {
		wrapped[i] = indexedTodo{todo: t, orig: i}
	}

	sort.SliceStable(wrapped, func(i, j int) bool {
		a, b := wrapped[i].todo, wrapped[j].todo
		switch m.sortMode {
		case sortUndoneFirst:
			if a.Done != b.Done {
				return !a.Done && b.Done
			}
			return false
		case sortDoneFirst:
			if a.Done != b.Done {
				return a.Done && !b.Done
			}
			return false
		case sortAlphaAsc:
			return strings.ToLower(a.Title) < strings.ToLower(b.Title)
		case sortAlphaDesc:
			return strings.ToLower(a.Title) > strings.ToLower(b.Title)
		default:
			return false
		}
	})

	hoveredOrig := m.cursor
	for i, w := range wrapped {
		todos[i] = w.todo
		if w.orig == hoveredOrig {
			m.cursor = i
		}
	}
	m.persist()
}

func (m model) activeKey() string {
	if m.noDay {
		return noDayKey
	}
	return dayKey(m.day)
}

func (m *model) moveCurrentItem() {
	todos := m.current()
	if len(todos) == 0 || m.cursor >= len(todos) {
		return
	}
	item := todos[m.cursor]
	todos = append(todos[:m.cursor], todos[m.cursor+1:]...)
	m.setCurrent(todos)
	m.clampCursor()

	destKey := noDayKey
	if m.noDay {
		destKey = dayKey(m.day)
	}
	m.store[destKey] = append(m.store[destKey], item)
	m.persist()
}
