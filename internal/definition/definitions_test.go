package definition_test

import (
	"testing"

	"github.com/jharlan-hash/gospell/internal/definition"
)

func TestState_GetDefinition(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		word string
		want string
	}{
		// TODO: Add test cases.
		{"TestWithExample", "example", "(1 of 6) noun: an item of information that is typical of a class or group"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m definition.State
			got := m.GetDefinition(tt.word)

			if got != tt.want {
				t.Errorf("GetDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGetDefinition(b *testing.B) {
	m := &definition.State{}
	for b.Loop() {
		m.GetDefinition("example")
	}
}

func TestState_NextDefinition(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{"TestWithExample", "(2 of 6) noun: a representative form or pattern"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s = &definition.State{}
			s.GetDefinition("example")

			got := s.NextDefinition()
			if got != tt.want {
				t.Errorf("NextDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNextDefinition(b *testing.B) {
	s := &definition.State{}
	s.GetDefinition("example")
	for b.Loop() {
		s.NextDefinition()
	}
}

func TestState_PrevDefinition(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{"TestWithExample", "(1 of 6) noun: an item of information that is typical of a class or group"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s = &definition.State{}
			s.GetDefinition("example")
			s.NextDefinition()

			got := s.PrevDefinition()

			if got != tt.want {
				t.Errorf("PrevDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkPrevDefinition(b *testing.B) {
	s := &definition.State{}
	s.GetDefinition("example")
	s.NextDefinition()
	for b.Loop() {
		s.PrevDefinition()
	}
}
