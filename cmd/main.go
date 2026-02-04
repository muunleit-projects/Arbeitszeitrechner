package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	azr "github.com/muunleit-projects/Arbeitszeitrechner"
)

// Styles.
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true)

	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)
)

type state int

const (
	stateInput state = iota
	stateDisplay
)

type model struct {
	state       state
	textInput   textinput.Model
	tableString string
	err         error
	azr         azr.Zeitpunkt // Use the interface/struct from the package if needed, or just helpers
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "hh:mm"
	ti.Focus()
	ti.CharLimit = 5
	ti.Width = 10

	// Initialize the azr object
	// We can ignore the error here for now as NewArbeitszeitrechner only errors on bad options
	// and we are not passing any options that would fail (defaults).
	// In a real app we might handle this better.
	z, _ := azr.NewArbeitszeitrechner()

	return model{
		state:     stateInput,
		textInput: ti,
		azr:       z,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

		switch m.state {
		case stateInput:
			switch msg.Type {
			case tea.KeyEnter:
				checkin := m.textInput.Value()
				// Validate/Calculate
				s, err := m.azr.TabelleString(checkin)
				if err != nil {
					m.err = err

					return m, nil
				}

				m.tableString = s
				m.state = stateDisplay
				m.err = nil

				return m, nil
			}

		case stateDisplay:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "r":
				m.state = stateInput
				m.textInput.Reset()
				m.tableString = ""

				return m, nil
			}
		}
	}

	if m.state == stateInput {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString("\n")
	s.WriteString(titleStyle.Render("Arbeitszeitrechner"))
	s.WriteString("\n\n")

	switch m.state {
	case stateInput:
		s.WriteString("Wann hast du eingecheckt?\n\n")
		s.WriteString(m.textInput.View())
		s.WriteString("\n\n")

		if m.err != nil {
			s.WriteString(errorStyle.Render(fmt.Sprintf("Fehler: %v", m.err)))
			s.WriteString("\n")
		}

		helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		s.WriteString(helpStyle.Render("(Format: hh:mm, Enter zum Bestätigen, Esc zum Beenden)"))

	case stateDisplay:
		s.WriteString(tableStyle.Render(m.tableString))
		s.WriteString("\n\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("(q zum Beenden, r zum Neustart)"))
	}

	s.WriteString("\n")

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
