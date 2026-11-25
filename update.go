package main

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
		if ws.Width > 4 {
			// account for "> " / "edit> " prompt + cursor padding
			m.input.SetWidth(ws.Width - 6)
		}
		return m, nil
	}

	switch m.mode {
	case modeAdding, modeEditing:
		return m.updateTextInput(msg)
	case modeConfirmingDelete:
		return m.updateConfirmingDelete(msg)
	case modeConfirmingMove:
		return m.updateConfirmingMove(msg)
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
			m.markDone()

		case "a", "n":
			return m, m.startAdding()

		case "e":
			return m, m.startEditing()

		case "d", "x":
			if len(m.current()) > 0 {
				m.mode = modeConfirmingDelete
			}

		case "left", "h", "[":
			if !m.noDay {
				m.day = m.day.AddDate(0, 0, -1)
				m.cursor = 0
				m.sortMode = sortManual
				m.clampCursor()
			} else {
				m.markDone()
			}

		case "right", "l", "]":
			if !m.noDay {
				m.day = m.day.AddDate(0, 0, 1)
				m.cursor = 0
				m.sortMode = sortManual
				m.clampCursor()
			} else {
				m.markDone()
			}

		case "t":
			m.noDay = false
			m.day = logicalDay(time.Now(), m.dayStart)
			m.cursor = 0
			m.sortMode = sortManual
			m.clampCursor()

		case "w":
			if m.noDay {
				m.noDay = false
				m.cursor = 0
				m.sortMode = sortManual
				m.clampCursor()
			} else {
				m.mode = modeOverview
				m.overviewWeek = startOfWeek(m.day)
				m.overviewCursor = dayIndexInWeek(m.day, m.overviewWeek)
			}
			return m, nil

		case "m":
			if len(m.current()) > 0 {
				m.mode = modeConfirmingMove
			}

		case "K":
			m.moveItem(-1)

		case "J":
			m.moveItem(1)

		case "s":
			m.cycleSort()
		}
	}
	return m, nil
}

func (m model) updateTextInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.mode = modeBrowsing
			m.input.Blur()
			return m, nil

		case "enter", "tab":
			title := strings.TrimSpace(m.input.Value())
			if title != "" {
				if m.mode == modeAdding {
					todos := append(m.current(), Todo{Title: title})
					m.setCurrent(todos)
					m.cursor = len(todos) - 1
					m.persist()
				} else {
					todos := m.current()
					if m.cursor < len(todos) {
						todos[m.cursor].Title = title
						m.persist()
					}
				}
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
				m.persist()
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

		case "esc":
			m.noDay = false
			m.mode = modeBrowsing
			return m, nil

		case "w":
			m.noDay = true
			m.cursor = 0
			m.sortMode = sortManual
			m.clampCursor()
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

		case "enter", "space":
			m.day = m.overviewWeek.AddDate(0, 0, m.overviewCursor)
			m.cursor = 0
			m.clampCursor()
			m.mode = modeBrowsing
		}
	}
	return m, nil
}

func (m model) updateConfirmingMove(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "m", "enter":
			m.moveCurrentItem()
			m.mode = modeBrowsing

		case "n", "esc":
			m.mode = modeBrowsing
		}
	}
	return m, nil
}

func (m model) markDone() {
	todos := m.current()
	if len(todos) > 0 {
		m.clampCursor()
		todos[m.cursor].Done = !todos[m.cursor].Done
		m.persist()
	}
}
