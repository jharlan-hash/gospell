package definition_test

import(
	"github.com/jharlan-hash/gospell/internal/definition"
	"testing"
)

func TestLoadCache(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want definition.Dictionary
	}{
        {"TestDictionaryLoading", definition.Dictionary{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := definition.LoadCache()

			if got["example"][0].Definition != "an item of information that is typical of a class or group" {
				t.Errorf("LoadCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkLoadCache(b *testing.B) {
    for b.Loop() {
        _ = definition.LoadCache()
    }
}
