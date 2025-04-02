package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jharlan-hash/gospell/internal/api"
	"github.com/jharlan-hash/gospell/internal/definition"
	"github.com/jharlan-hash/gospell/internal/tts"
	"github.com/jharlan-hash/gospell/internal/wpm"
	"github.com/muesli/reflow/wordwrap"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pborman/getopt"
	"google.golang.org/api/option"
)

func main() {
	credentialFlag := getopt.StringLong("credentials", 'c', "", "Path to Google Cloud credentials JSON file (optional)")
	helpFlag := getopt.BoolLong("help", 'h', "display help")

	getopt.Parse()

	if *helpFlag {
		getopt.Usage()
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	model := initialModel(*credentialFlag, ctx)

	p := tea.NewProgram(&model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type wordMessage struct {
	word       string
	definition string
}

type correctMessage struct{}
type incorrectMessage struct{}

type model struct {
	textInput       textinput.Model
	streak          int
	correction      string
	definition      string
	credentialPath  string
	word            string
	initialTime     time.Time
	finalTime       time.Time
	width           int
	height          int
	definitionState *definition.State
	ttsState        *tts.TTS
	borderColor     lipgloss.Color
}

// initialModel initializes the model with a text input field and a random word.
func initialModel(credentialPath string, ctx context.Context) model {
	ti := textinput.New()
	ti.Placeholder = "spell spoken word..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	state := &definition.State{}

	ttsState := &tts.TTS{}
	ttsState.Ctx = ctx

	// Get a random word and its definition.
	word := api.GetAcceptableWord()
	return model{
		textInput:       ti,
		credentialPath:  credentialPath,
		correction:      "\n",
		word:            word,
		definitionState: state,
		definition:      state.GetDefinition(word),
		ttsState:        ttsState,
	}
}

func (m *model) Init() tea.Cmd {
	if m.credentialPath != "" {
		// User provided custom credentials file
		client, err := texttospeech.NewClient(m.ttsState.Ctx, option.WithCredentialsFile(m.credentialPath))
		if err != nil {
			tea.ExitAltScreen()
			log.Fatal("Bad credentials file - make sure the path is correct\n" + err.Error())
		}
		m.ttsState.Client = client
	} else {
		tea.ExitAltScreen()
		log.Fatal("Please provide a Google Cloud credentials file.")
	}

	m.ttsState.Word = m.word
	go m.ttsState.SayWord()
	return textinput.Blink
}

// Command to generate a new word.
func getNewWord(m *model) tea.Cmd {
	return func() tea.Msg {
		word := api.GetAcceptableWord()
		def := m.definitionState.GetDefinition(word)

		// Play the word audio in a goroutine to avoid blocking.
		m.ttsState.Word = word
		go m.ttsState.SayWord()

		return wordMessage{
			word:       word,
			definition: def,
		}
	}
}

// Update function handling messages
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case wordMessage:
		// Update model with new word.
		m.word = msg.word
		m.definition = msg.definition
		return m, nil

	case tea.KeyMsg:
		if m.textInput.Value() == "" && msg.Type != tea.KeyCtrlR {
			m.initialTime = time.Now() // start timer on first key press.
		}

		m.finalTime = time.Now() // update timer on every key press.

		switch msg.Type {
		case tea.KeyEnter: // submit word while ignoring empty input.
			return m.submitWord()
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyCtrlD: // exit.
			return m, tea.Quit
		case tea.KeyCtrlR: // repeat word.
			m.ttsState.Word = m.word
			go m.ttsState.SayWord()
			return m, nil
		case tea.KeyDown:
			// If the user presses down, we want to get the next definition.
			m.definition = m.definitionState.NextDefinition()
		case tea.KeyUp:
			// If the user presses up, we want to get the previous definition.
			m.definition = m.definitionState.PrevDefinition()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height


	case correctMessage:
         correctColor := lipgloss.Color("#66ac5a") // og
		m.streak++
		m.definition = wordwrap.String(m.definition, 100)
		m.borderColor = correctColor // Set border color to green for correct answer

		m.correction = ""
		return m, getNewWord(m)

	case incorrectMessage:
         incorrectColor := lipgloss.Color("#ED4337") // og
		m.streak = 0
		m.definition = wordwrap.String(m.definition, 100)
		m.borderColor = incorrectColor // Set border color to red for incorrect answer

		m.correction = fmt.Sprintf("Correct spelling: %s", m.word)
		return m, getNewWord(m)
	}

	// Handle text input updates.
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// submitWord checks the user's input against the correct word.
// If the input is correct, it returns a correctMessage.
// If the input is incorrect, it returns an incorrectMessage.
// It also resets the text input field.
func (m *model) submitWord() (tea.Model, tea.Cmd) {
	if m.textInput.Value() == "" {
		return m, nil
	}

	userInput := m.textInput.Value()
	m.textInput.Reset()

	if userInput == m.word { // Correct answer.
		return m, func() tea.Msg { return correctMessage{} }
	} 
	
	// Incorrect answer.
	return m, func() tea.Msg { return incorrectMessage{} }


}

func (m model) View() string {
     foregroundColor := lipgloss.Color(1)
     backgroundColor := lipgloss.Color(1)

	// Create a container style for the main content
	inputContainer := lipgloss.NewStyle().
		Padding(1, 2).
		Margin(1).
        Foreground(lipgloss.Color(foregroundColor)).
        Background(lipgloss.Color(backgroundColor)).
        BorderForeground(lipgloss.Color(foregroundColor)).
        BorderBackground(lipgloss.Color(backgroundColor)).
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Center, lipgloss.Center)

	width := 75

	// Center the input field within the container
	inputView := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(m.borderColor).
        Foreground(lipgloss.Color(foregroundColor)).
        Background(lipgloss.Color(backgroundColor)).
        BorderBackground(lipgloss.Color(backgroundColor)).
		Render(m.textInput.View())

	// Center the definition but keep it within the container's width
	definitionText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(width).
		Render(m.definition)

	correctionText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(width).
		Render(m.correction)

	// Combine all elements with the container style
	content := inputContainer.Render(
		inputView + "\n" +
			definitionText + "\n" +
			correctionText,
	)

	// Style for the status bar at the bottom
	renderString := fmt.Sprintf(
		"Gospell: Press 'ESC' / 'CtrlC' to exit, 'CtrlR' to repeat word, ↑/↓ to navigate definitions | Current WPM: %d | Streak: %d",
		wpm.CalculateWpm(m.textInput.Value(), m.initialTime, m.finalTime),
		m.streak,
	)

	statusBar := lipgloss.NewStyle().
		Background(lipgloss.Color("#cfd6f1")).
		Foreground(lipgloss.Color("#1e1e2d")).
		Width(m.width). // Make it full width
		Align(lipgloss.Left).
		Render(renderString)

	// Use Place for the main content, positioning it in the center
	mainContent := lipgloss.Place(
		m.width,
		m.height-lipgloss.Height(statusBar), // leave room for status bar
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	// Build the final view with the status bar at the bottom
	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		lipgloss.PlaceHorizontal(
			m.width,
			lipgloss.Left,
			statusBar,
		),
	)
}
