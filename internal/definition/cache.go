package definition

import (
	"embed"
	"encoding/gob"
	"fmt"
	"log"
)

//go:embed wordmap.gob
var fs embed.FS

// LoadCache loads the wordmap cache from a file.
// It returns a Dictionary, which is a map of words to an Entry slice.
func LoadCache() (Dictionary, err) {
	f, err := fs.Open("wordmap.gob")
	if err != nil {
		return nil, fmt.Errorf("Unable to open wordmap. %v", err)
	}

	cache := make(Dictionary)
	if err := gob.NewDecoder(f).Decode(&cache); err != nil {
		return nil, fmt.Errorf("Error decoding cache.  %v", err)
	}

	return cache, nil
}
