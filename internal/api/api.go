package api

import (
	"bufio"
	"embed"
	"log"
	"math/rand"
	"time"
)

// GetAcceptableWord returns a random word from the wordlist.
func GetAcceptableWord() string {
	return getRandomLineFromWordlist()
}

//go:embed wordlist.txt
var fs embed.FS

// From The Art of Computer Programming, Volume 2, Section 3.4.2, by Donald E. Knuth.
// This is a reservoir sampling algorithm that picks a random line from a file.
func getRandomLineFromWordlist() string {
	file, err := fs.Open("wordlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	randsource := rand.NewSource(time.Now().UnixNano())
	randgenerator := rand.New(randsource)

	lineNum := 1
	var pick string
	for scanner.Scan() {
		line := scanner.Text()
		// Instead of 1 to N it's 0 to N-1
		roll := randgenerator.Intn(lineNum)
		if roll == 0 {
			// We pick this line
			pick = line
		}

		lineNum += 1
	}

	return pick
}
