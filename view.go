// ./view.go
package main

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	deleteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			Bold(true)

	editStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)
)

func (m model) View() tea.View {
	if m.mode == modeOverview {
		v := tea.NewView(m.overviewView())
		v.AltScreen = true
		return v
	}

	todos := m.current()

	reserved := 2
	scrolled := false
	if m.height > 0 && len(todos) > 0 {
		avail := max(1, m.height-reserved)
		if len(todos) > avail {
			scrolled = true
			reserved += 2
		}
	}

	switch {
	case m.mode == modeConfirmingDelete && m.cursor < len(todos):
		reserved += 2
	case m.mode != modeBrowsing:
		reserved += 3
	case m.showHelp:
		reserved += 4
	}
	if m.statusErr != "" {
		reserved += 2
	}

	start, end := m.visibleRange(len(todos), reserved)

	s := truncateLine(m.header(), m.width) + "\n\n"

	if len(todos) == 0 {
		s += truncateLine("  nothing yet - press a to jot one down", m.width) + "\n"
	}
	for i := start; i < end; i++ {
		todo := todos[i]
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if todo.Done {
			checked = "x"
		}

		line := truncateLine(fmt.Sprintf("%s [%s] %s", cursor, checked, todo.Title), m.width)

		switch {
		case m.mode == modeConfirmingDelete && m.cursor == i:
			line = deleteStyle.Render(line)
		case m.mode == modeEditing && m.cursor == i:
			line = editStyle.Render(line)
		}

		s += line + "\n"
	}

	if scrolled {
		s += truncateLine(fmt.Sprintf("\n  item %d/%d", m.cursor+1, len(todos)), m.width) + "\n"
	}

	switch {
	case m.mode == modeConfirmingDelete:
		if m.cursor < len(todos) {
			s += "\n" + truncateLine("Press 'd' again to delete", m.width) + "\n"
		}
	case m.mode != modeBrowsing:
		s += "\n" + m.input.View() + "\n"
		s += "\n" + truncateLine("enter to save · esc to cancel", m.width) + "\n"
	case m.showHelp:
		s += "\n" + truncateLine("h/l day · j/k move · J/K reorder · s sort", m.width) + "\n"
		s += truncateLine("a add · e edit · d delete · w week · t today · q quit", m.width) + "\n"
		s += truncateLine("press ? to hide help", m.width) + "\n"
	}

	if m.statusErr != "" {
		s += "\n" + deleteStyle.Render(truncateLine("! "+m.statusErr, m.width)) + "\n"
	}

	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func (m model) header() string {
	today := logicalDay(time.Now(), m.dayStart)
	base := m.day.Format("Mon Jan 2")

	var label string
	switch {
	case m.day.Equal(today):
		label = "Today · " + base
	case m.day.Equal(today.AddDate(0, 0, -1)):
		label = "Yesterday · " + base
	case m.day.Equal(today.AddDate(0, 0, 1)):
		label = "Tomorrow · " + base
	default:
		label = base
	}

	if m.dayStart != 0 {
		label += "  (day starts " + formatClock(m.dayStart) + ")"
	}

	if sLabel := m.sortMode.label(); sLabel != "" {
		label += " [" + sLabel + "]"
	}
	return label
}

func (m model) overviewView() string {
	end := m.overviewWeek.AddDate(0, 0, 6)
	header := fmt.Sprintf("Week of %s - %s, %d",
		m.overviewWeek.Format("Jan 2"),
		end.Format("Jan 2"),
		m.overviewWeek.Year())
	s := truncateLine(header, m.width) + "\n\n"

	today := logicalDay(time.Now(), m.dayStart)

	for i := 0; i < 7; i++ {
		day := m.overviewWeek.AddDate(0, 0, i)
		todos := m.store[dayKey(day)]

		done := 0
		for _, t := range todos {
			if t.Done {
				done++
			}
		}
		total := len(todos)

		cursor := " "
		if m.overviewCursor == i {
			cursor = ">"
		}

		marker := " "
		if day.Equal(today) {
			marker = "*"
		}

		var count string
		switch {
		case total == 0:
			count = "-"
		case done == total:
			count = fmt.Sprintf("%d (all done)", total)
		case done > 0:
			count = fmt.Sprintf("%d (%d done)", total, done)
		default:
			count = fmt.Sprintf("%d", total)
		}

		line := truncateLine(fmt.Sprintf("%s%s %-12s  %s",
			cursor, marker, day.Format("Mon Jan 2"), count), m.width)

		if m.overviewCursor == i {
			line = editStyle.Render(line)
		}

		s += line + "\n"
	}

	s += "\n" + truncateLine("h/l prev/next week · j/k move · enter open day · esc back", m.width) + "\n"
	return s
}

func truncateLine(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) <= maxWidth {
		return s
	}
	if maxWidth == 1 {
		return "…"
	}
	return string(runes[:maxWidth-1]) + "…"
}

func (m model) visibleRange(total, reservedLines int) (start, end int) {
	if total == 0 {
		return 0, 0
	}
	if m.height <= 0 {
		return 0, total
	}
	avail := max(1, m.height-reservedLines)
	if total <= avail {
		return 0, total
	}
	start = m.cursor - avail/2
	if start < 0 {
		start = 0
	}
	end = start + avail
	if end > total {
		end = total
		start = max(0, end-avail)
	}
	return start, end
}
