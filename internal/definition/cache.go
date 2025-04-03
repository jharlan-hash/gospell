package definition

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"fmt"
)

//go:embed wordmap.gob
var fileBytes []byte

// LoadCache loads the wordmap cache from a file.
// It returns a Dictionary, which is a map of words to an Entry slice.
func LoadCache() Dictionary {
	read := bytes.NewReader(fileBytes)

	enc := gob.NewDecoder(read)
	cache := make(Dictionary)

	if err := enc.Decode(&cache); err != nil {
		fmt.Println("Error decoding cache:", err)
	}

	return cache
}
