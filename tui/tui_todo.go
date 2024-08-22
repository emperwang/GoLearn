package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

type model struct {
	todos    []string
	cursor   int
	selected map[int]struct{}
}

func (m *model) Init() tea.Cmd {

	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.todos)-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m *model) View() string {
	s := "todo list: \n\n"

	for i, choice := range m.todos {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\n Press q to quit. \n"
	return s
}

var initModel = model{
	todos:    []string{"cleaning", "wash clothes", "write a blog"},
	selected: make(map[int]struct{}),
}

func TuiStarter() {
	cmd := tea.NewProgram(&initModel)
	if _, err := cmd.Run(); err != nil {
		log.Errorf("tui started failed. %v", err)
	}
}
