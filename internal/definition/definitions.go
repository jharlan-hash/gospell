package definition

import (
	"fmt"
)

type State struct {
	Cache       Dictionary
	Word        string
	Index       int
	Definitions []string
}

// getDefinitionList returns a list of definitions for a given word from the cache.
// It populates the definitions field in the State struct.
// This function is called internally by GetDefinition to initialize the definitions list.
func (s *State) GetDefinitionList() {
	definitions := s.Cache[s.Word]
	list := make([]string, 0)

	for _, definition := range definitions {
		list = append(list,
			fmt.Sprintf(
				"(%d of %d) %s: %s",
				definition.DefinitionIndex,
				definition.NumDefinitions,
				definition.PartOfSpeech,
				definition.Definition,
			),
		)
	}

	s.Definitions = list // store the definitions in the state
}

// NextDefinition retrieves the next definition of a word from the cache.
// If the user requests a definition past the last one, it returns the last definition.
func (s *State) NextDefinition() string {
	if s.Index+1 >= len(s.Definitions) { // if user requests something past the end of the definition list
		return s.Definitions[len(s.Definitions)-1] // return the last definition
	} else { // increment index & change definition
		s.Index++
		return s.Definitions[s.Index]
	}
}

// PrevDefinition retrieves the previous definition of a word from the cache.
// If the user requests a definition before the first one, it returns the first definition.
func (s *State) PrevDefinition() string {
	if s.Index-1 < 0 { // if user requests something before the beginning of the definition list
		return s.Definitions[0] // return the first definition
	} else { // decrement index & change definition
		s.Index--
		return s.Definitions[s.Index]
	}
}

// GetDefinition retrieves the first definition of a word from the cache.
//
// How To Use:
//
// 1. Create a new State instance.
//
// 2. Call GetDefinition with the word you want to look up.
//
// 3. Use NextDefinition() and PrevDefinition() to navigate through the definitions.
//
// Example:
//
//	state := &definition.State{}
//	firstDef := state.GetDefinition("example") // retrieves the first definition
//	fmt.Println(firstDef) // prints the first definition
//	nextDef := state.NextDefinition() // retrieves the next definition
//	fmt.Println(nextDef) // prints the next definition
//	prevDef := state.PrevDefinition() // retrieves the previous definition
//	fmt.Println(prevDef) // prints the previous definition
func (m *State) GetDefinition(word string) string {
	if m.Cache == nil { // Only load the cache if it's not already loaded.
		m.Cache = LoadCache()
	}

	m.Word = word
	m.Index = 0
	m.GetDefinitionList() // populate the definitions list
	return m.Definitions[m.Index] // return the first definition
}
