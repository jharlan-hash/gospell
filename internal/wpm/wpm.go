package wpm

import (
	"math"
	"time"
	"unicode/utf8"
)

var wpm int

// CalculateWpm calculates the words per minute (WPM) based on the given word and time interval.
// It takes a word (string), the initial time (time.Time), and the final time (time.Time) as input.
// It returns the calculated WPM as an integer.
func CalculateWpm(word string, initialTime time.Time, finalTime time.Time) int {
	// If the word is empty, return the last calculated wpm
	if word == "" {
		return wpm
	}

	timeElapsed := finalTime.Sub(initialTime)

	// Convert to minutes, and handle very small time intervals
	elapsedMinutes := timeElapsed.Minutes()
	if elapsedMinutes < 0.001 { // Avoid division by very small numbers
		return 0
	}

	numWords := float64(utf8.RuneCountInString(word)) / 5.0

    return int(math.Round(numWords / elapsedMinutes))
}
