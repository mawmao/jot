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

	s := m.header() + "\n\n"

	todos := m.current()
	if len(todos) == 0 {
		s += "  nothing yet - press a to jot one down\n"
	}
	for i, todo := range todos {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if todo.Done {
			checked = "x"
		}

		line := fmt.Sprintf("%s [%s] %s", cursor, checked, todo.Title)

		switch {
		case m.mode == modeConfirmingDelete && m.cursor == i:
			line = deleteStyle.Render(line)
		case m.mode == modeEditing && m.cursor == i:
			line = editStyle.Render(line)
		}

		s += line + "\n"
	}

	switch {
	case m.mode == modeConfirmingDelete:
		todos := m.current()
		if m.cursor < len(todos) {
			s += "\nPress 'd' again to delete\n"
		}
	case m.mode != modeBrowsing:
		s += "\n" + m.input.View() + "\n"
		s += "\nenter to save · esc to cancel\n"
	case m.showHelp:
		s += "\nh/l prev/next day · j/k move · a add · e edit · d delete\n"
		s += "w week overview · t today · ? help · q quit\n"
		s += "press ? to hide help\n"
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
	return label
}

func (m model) overviewView() string {
	end := m.overviewWeek.AddDate(0, 0, 6)
	s := fmt.Sprintf("Week of %s - %s, %d\n\n",
		m.overviewWeek.Format("Jan 2"),
		end.Format("Jan 2"),
		m.overviewWeek.Year())

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

		line := fmt.Sprintf("%s%s %-12s  %s",
			cursor, marker, day.Format("Mon Jan 2"), count)

		if m.overviewCursor == i {
			line = editStyle.Render(line)
		}

		s += line + "\n"
	}

	s += "\nh/l prev/next week · j/k move · enter open day · esc back\n"
	return s
}
