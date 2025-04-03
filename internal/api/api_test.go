package api

import "testing"

func TestGetRandomLineFromWordlist(t *testing.T) {
	word := RandomWord()
	t.Logf("Random word: %v", word)

	if len(word) == 0 {
		t.Errorf("Expected a non-empty word, got empty string")
	}
	if len(word) > 29 {
		t.Errorf("Expected a word of length <= 29, got %d", len(word))
	}
}

func TestSplitWords(t *testing.T) {
    words := splitWords(fileString)
    if len(words) == 0 {
        t.Errorf("Expected a non-empty list of words, got empty list")
    }
    if len(words) != 90632 {
        t.Errorf("Expected 90632 words, got %d", len(words))
    }
}

func BenchmarkRuntimeRandomWord(b *testing.B) {
	var word string
	for b.Loop() {
		word = RandomWord()
	}
	b.Logf("Random word: %v", word)
}

func BenchmarkSplitWords(b *testing.B) {
	for b.Loop() {
		_ = splitWords(fileString)
	}
}
