package resolver

import (
	"fmt"
	"testing"

	"github.com/envlens/internal/parser"
)

// BenchmarkResolve_Large measures resolution throughput for a file with many
// cross-referencing entries.
func BenchmarkResolve_Large(b *testing.B) {
	const n = 200
	entries := make([]parser.Entry, 0, n)
	// First half: plain values.
	for i := 0; i < n/2; i++ {
		entries = append(entries, parser.Entry{
			Key:   fmt.Sprintf("BASE_%d", i),
			Value: fmt.Sprintf("value_%d", i),
		})
	}
	// Second half: references to the first half.
	for i := 0; i < n/2; i++ {
		entries = append(entries, parser.Entry{
			Key:   fmt.Sprintf("REF_%d", i),
			Value: fmt.Sprintf("prefix_${BASE_%d}_suffix", i),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Resolve(entries, StrategyStrict)
		if err != nil {
			b.Fatal(err)
		}
	}
}
