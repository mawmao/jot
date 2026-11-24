package main

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case modeAdding:
		return m.updateAdding(msg)
	case modeEditing:
		return m.updateEditing(msg)
	case modeConfirmingDelete:
		return m.updateConfirmingDelete(msg)
	case modeOverview:
		return m.updateOverview(msg)
	}
	return m.updateBrowsing(msg)
}

func (m model) updateBrowsing(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "?":
			m.showHelp = !m.showHelp

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.current())-1 {
				m.cursor++
			}

		case "enter", "space":
			todos := m.current()
			if len(todos) > 0 {
				m.clampCursor()
				todos[m.cursor].Done = !todos[m.cursor].Done
				_ = saveStore(m.store)
			}

		case "a", "n":
			m.mode = modeAdding
			m.input.Prompt = "> "
			m.input.SetValue("")
			return m, m.input.Focus()

		case "e":
			todos := m.current()
			if len(todos) > 0 {
				m.mode = modeEditing
				m.input.Prompt = "edit> "
				m.input.SetValue(todos[m.cursor].Title)
				return m, m.input.Focus()
			}

		case "d", "x":
			if len(m.current()) > 0 {
				m.mode = modeConfirmingDelete
			}

		case "left", "h", "[":
			m.day = m.day.AddDate(0, 0, -1)
			m.cursor = 0
			m.clampCursor()

		case "right", "l", "]":
			m.day = m.day.AddDate(0, 0, 1)
			m.cursor = 0
			m.clampCursor()

		case "w":
			m.mode = modeOverview
			m.overviewWeek = startOfWeek(m.day)
			m.overviewCursor = dayIndexInWeek(m.day, m.overviewWeek)
			return m, nil

		case "t":
			m.day = logicalDay(time.Now(), m.dayStart)
			m.cursor = 0
			m.clampCursor()
		}
	}
	return m, nil
}

func (m model) updateAdding(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.mode = modeBrowsing
			m.input.Blur()
			return m, nil

		case "enter":
			title := strings.TrimSpace(m.input.Value())
			if title != "" {
				todos := append(m.current(), Todo{Title: title})
				m.setCurrent(todos)
				m.cursor = len(todos) - 1
				_ = saveStore(m.store)
			}
			m.mode = modeBrowsing
			m.input.Blur()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) updateEditing(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.mode = modeBrowsing
			m.input.Blur()
			return m, nil

		case "enter":
			title := strings.TrimSpace(m.input.Value())
			if title != "" {
				todos := m.current()
				todos[m.cursor].Title = title
				_ = saveStore(m.store)
			}
			m.mode = modeBrowsing
			m.input.Blur()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) updateConfirmingDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {

		case "d", "enter":
			todos := m.current()
			if m.cursor < len(todos) {
				todos = append(todos[:m.cursor], todos[m.cursor+1:]...)
				m.setCurrent(todos)
				if m.cursor >= len(todos) && m.cursor > 0 {
					m.cursor--
				}
				_ = saveStore(m.store)
			}
			m.mode = modeBrowsing

		case "n", "esc":
			m.mode = modeBrowsing
		}
	}
	return m, nil
}

func (m model) updateOverview(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc", "w":
			m.mode = modeBrowsing
			return m, nil

		case "up", "k":
			if m.overviewCursor > 0 {
				m.overviewCursor--
			}

		case "down", "j":
			if m.overviewCursor < 6 {
				m.overviewCursor++
			}

		case "left", "h", "[":
			m.overviewWeek = m.overviewWeek.AddDate(0, 0, -7)

		case "right", "l", "]":
			m.overviewWeek = m.overviewWeek.AddDate(0, 0, 7)

		case "t":
			now := logicalDay(time.Now(), m.dayStart)
			m.overviewWeek = startOfWeek(now)
			m.overviewCursor = dayIndexInWeek(now, m.overviewWeek)

		case "enter":
			m.day = m.overviewWeek.AddDate(0, 0, m.overviewCursor)
			m.cursor = 0
			m.clampCursor()
			m.mode = modeBrowsing
		}
	}
	return m, nil
}
