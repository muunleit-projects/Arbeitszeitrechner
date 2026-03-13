package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	azr "github.com/muunleit-projects/Arbeitszeitrechner"
)

// Styles — refined color palette for a premium TUI feel.
var (
	// Accent colors
	purple    = lipgloss.Color("#7C3AED")
	violet    = lipgloss.Color("#8B5CF6")
	indigo    = lipgloss.Color("#6366F1")
	slate     = lipgloss.Color("#94A3B8")
	dimSlate  = lipgloss.Color("#64748B")
	surface   = lipgloss.Color("#1E1B4B")
	red       = lipgloss.Color("#F87171")
	green     = lipgloss.Color("#34D399")
	white     = lipgloss.Color("#F8FAFC")
	dimWhite  = lipgloss.Color("#CBD5E1")

	// Title bar with gradient-like effect
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(purple).
			Padding(0, 2).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(violet).
			Italic(true)

	// Error display
	errorStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(red)

	// Table output
	tableStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(indigo).
			Padding(1, 2)

	// Help / hint text
	helpStyle = lipgloss.NewStyle().
			Foreground(dimSlate).
			Italic(true).
			MarginTop(1)

	// Prompt label
	promptStyle = lipgloss.NewStyle().
			Foreground(slate).
			Bold(true)

	// Status / success
	successStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	// Outer container
	containerStyle = lipgloss.NewStyle().
			Padding(1, 2)
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
	azr         azr.Zeitpunkt
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "hh:mm"
	ti.Focus()
	ti.CharLimit = 5
	ti.SetWidth(10)

	// Customize text input styles for a polished look
	styles := ti.Styles()
	styles.Focused.Prompt = lipgloss.NewStyle().Foreground(violet).Bold(true)
	styles.Focused.Text = lipgloss.NewStyle().Foreground(white)
	styles.Focused.Placeholder = lipgloss.NewStyle().Foreground(dimSlate)
	ti.SetStyles(styles)

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
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

		switch m.state {
		case stateInput:
			if msg.String() == "enter" {
				checkin := m.textInput.Value()
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

func (m model) View() tea.View {
	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render("⏱  Arbeitszeitrechner"))
	s.WriteString("\n")

	switch m.state {
	case stateInput:
		s.WriteString(subtitleStyle.Render("Arbeitszeit berechnen"))
		s.WriteString("\n\n")
		s.WriteString(promptStyle.Render("Wann hast du eingecheckt?"))
		s.WriteString("\n\n")
		s.WriteString("  " + m.textInput.View())
		s.WriteString("\n")

		if m.err != nil {
			s.WriteString("\n")
			s.WriteString(errorStyle.Render(fmt.Sprintf("✗ Fehler: %v", m.err)))
			s.WriteString("\n")
		}

		s.WriteString("\n")
		s.WriteString(helpStyle.Render("  Format: hh:mm  •  Enter bestätigen  •  Esc beenden"))

	case stateDisplay:
		s.WriteString(subtitleStyle.Render("Deine Arbeitszeiten"))
		s.WriteString("\n\n")
		s.WriteString(tableStyle.Render(m.tableString))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render(
			"  " + successStyle.Render("r") + " Neustart  •  " +
				successStyle.Render("q") + " Beenden",
		))
	}

	s.WriteString("\n")

	content := containerStyle.Render(s.String())

	return tea.NewView(content)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
