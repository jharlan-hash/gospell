package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jharlan-hash/gospell/internal/api"
	"github.com/jharlan-hash/gospell/internal/definition"
	"github.com/jharlan-hash/gospell/internal/tts"
	"github.com/jharlan-hash/gospell/internal/wpm"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/muesli/reflow/wordwrap"
	"github.com/pborman/getopt"
	"github.com/tjarratt/babble"
	"google.golang.org/api/option"
)

var debug = true

func main() {
	credentialFlag := getopt.StringLong("credentials", 'c', "", "Path to Google Cloud credentials JSON file")
	numDefinitionsFlag := getopt.IntLong("definitions", 'd', 1, "Number of definitions to display")
	helpFlag := getopt.BoolLong("help", 'h', "display help")

	getopt.Parse()


	if *helpFlag {
		getopt.Usage()
		os.Exit(0)
	}

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

	model := initialModel(*credentialFlag, *numDefinitionsFlag)
    model.ctx = ctx

	p := tea.NewProgram(&model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type wordMessage struct {
	word       string
	definition string
}

type correctMessage struct{}
type incorrectMessage struct{}

type model struct {
	textInput      textinput.Model
	client         *texttospeech.Client
	ctx            context.Context
	err            error
	cache          map[string]struct{}
	babbler        babble.Babbler
	definition     string
	numDefinitions int
	credentialPath string
	word           string
	initialTime    time.Time
	finalTime      time.Time
	streak         int
	correction     string
	width          int
	height         int
}

func initialModel(credentialPath string, numDefinitions int) model {
	ti := textinput.New()
	ti.Placeholder = "spell spoken word..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput:      ti,
		client:         nil,
		ctx:            nil,
		err:            nil,
		cache:          make(map[string]struct{}),
		babbler:        babble.NewBabbler(),
		definition:     "",
		credentialPath: credentialPath,
		word:           "",
		streak:         0,
		correction:     "\n",
		numDefinitions: numDefinitions,
	}
}

func (m *model) Init() tea.Cmd {
	if m.credentialPath != "" {
		var err error

		client, err := texttospeech.NewClient(m.ctx, option.WithCredentialsFile(m.credentialPath))
		if err != nil {
			log.Fatal("Bad credentials file")
		}

		m.client = client
	} else {
		log.Fatal("No credentials file provided")
	}

	m.cache = definition.LoadCache()
	m.babbler.Count = 1

	m.word = api.GetAcceptableWord(m.babbler, m.cache)
	res := definition.GetResponse(m.word)
	m.definition = definition.GetDefinition(res, m.numDefinitions)
	go tts.SayWord(m.ctx, *m.client, m.word)

	return textinput.Blink
}

// Command to generate a new word
func getNewWord(m *model) tea.Cmd {
	return func() tea.Msg {
		word := api.GetAcceptableWord(m.babbler, m.cache)
		res := definition.GetResponse(word)
		def := definition.GetDefinition(res, m.numDefinitions)

		// Play the word audio in a goroutine to avoid blocking
		go tts.SayWord(m.ctx, *m.client, word)

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
		// Update model with new word
		m.word = msg.word
		m.definition = msg.definition
		return m, nil

	case tea.KeyMsg:
		if m.textInput.Value() == "" && msg.Type != tea.KeyCtrlR {
			m.initialTime = time.Now()
		}

        m.finalTime = time.Now()
		switch msg.Type {
		case tea.KeyEnter:
            // ignoring a blank input
            if (m.textInput.Value() == "") {
                return m, nil
            }

			userInput := m.textInput.Value()
            m.textInput.Reset()

			if userInput == m.word {
				// Correct answer
				return m, func() tea.Msg { return correctMessage{} }
			} else {
				// Incorrect answer
				return m, func() tea.Msg { return incorrectMessage{} }
			}
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyCtrlD:
			definition.SaveCache(&m.cache)
			os.Remove("temp.wav")
			return m, tea.Quit
		case tea.KeyCtrlR: // repeat word
			go api.PlayWav("temp.wav")
			return m, nil
		case tea.KeyCtrlH: // say what's in the text box
			go tts.SayWord(m.ctx, *m.client, m.textInput.Value())
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case correctMessage:
		m.streak++
		lines := strings.Split(m.definition, "\n")
		for i := range lines {
			lines[i] = color.GreenString(lines[i])
		}
		m.definition = strings.Join(lines, "\n")

		m.correction = "\n"
		return m, getNewWord(m)

	case incorrectMessage:
		m.streak = 0
		m.definition = color.RedString(m.definition)
		m.correction = fmt.Sprintf("\nCorrect spelling: %s", m.word)
		return m, getNewWord(m)
	}

	// Handle text input updates
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	// Create a container style that will be centered as a whole
	containerStyle := lipgloss.NewStyle().Padding(1, 2)

	// Center the welcome text within the container
	welcomeText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100). // Set a fixed width for the centered elements
		Render("Welcome to gospell!")

	// Center the input field within the container
	inputView := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100).
		Render(m.textInput.View())

		// Center the definition but keep it within the container's width
	definitionText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100).
		Render(wordwrap.String(m.definition, 100)) // Wrap text to fit within the container

	correctionText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100).
		Render(m.correction)

    // debugText := lipgloss.NewStyle().
    //     Align(lipgloss.Center).
    //     Width(100).
    //     Render(fmt.Sprintf("Word: %s", m.word))

	wpmText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100).
		Render(fmt.Sprintf("WPM: %d", wpm.CalculateWpm(m.textInput.Value(), m.initialTime, m.finalTime)))

	// Center the streak counter
	streakText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100).
		Render("Streak: " + fmt.Sprintf("%d", m.streak))

	// Combine all elements with the container style
	content := containerStyle.Render(
        // debugText + "\n" +
		welcomeText + "\n\n" +
			wpmText + "\n" +
			inputView + "\n" +
			definitionText + "\n" +
			correctionText + "\n" +
			streakText,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
