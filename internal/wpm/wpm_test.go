package wpm_test

import (
	"github.com/jharlan-hash/gospell/internal/wpm"
	"testing"
	"time"
)

func TestCalculateWpm(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		word        string
		initialTime time.Time
		finalTime   time.Time
		want        int
	}{
		{"Test with 1 word in 1 minute", "hello", time.Now().Add(-1 * time.Minute), time.Now(), 1},
		{"Test with 1 word in 30 seconds", "hello", time.Now().Add(-30 * time.Second), time.Now(), 2},
		{"Test with 5 words in 1 minute", "hello world this is a test", time.Now().Add(-1 * time.Minute), time.Now(), 5},
		{"Test with 5 words in 30 seconds", "hello world this is a test", time.Now().Add(-30 * time.Second), time.Now(), 10},
		{"Test with 0 words", "", time.Now().Add(-1 * time.Minute), time.Now(), 0},                                                                                                                                                       // Edge case: empty word
		{"Test with very small time interval", "hello", time.Now().Add(-1 * time.Millisecond), time.Now(), 0},                                                                                                                            // Edge case: very small time interval
		{"Test with negative time interval", "hello", time.Now().Add(1 * time.Minute), time.Now(), 0},                                                                                                                                    // Edge case: final time before initial time
		{"Test with 10 words in 2 minutes", "this is a test of the emergency broadcast system", time.Now().Add(-2 * time.Minute), time.Now(), 5},                                                                                         // 10 words in 2 minutes
		{"Test with 15 words in 3 minutes", "this is a test of the emergency broadcast system for testing purposes only", time.Now().Add(-3 * time.Minute), time.Now(), 5},                                                               // 15 words in 3 minutes
		{"Test with 20 words in 4 minutes", "this is a test of the emergency broadcast system for testing purposes only and more words added here", time.Now().Add(-4 * time.Minute), time.Now(), 5},                                     // 20 words in 4 minutes
		{"Test with 25 words in 5 minutes", "this is a test of the emergency broadcast system for testing purposes only and more words added here to increase the count", time.Now().Add(-5 * time.Minute), time.Now(), 5},               // 25 words in 5 minutes
		{"Test with 30 words in 6 minutes", "this is a test of the emergency broadcast system for testing purposes only and more words added here to increase the count significantly", time.Now().Add(-6 * time.Minute), time.Now(), 5}, // 30 words in 6 minutes
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wpm.CalculateWpm(tt.word, tt.initialTime, tt.finalTime)
			if got != tt.want {
				t.Errorf("CalculateWpm() = %v, want %v", got, tt.want)
			}
		})
	}
}
