package green_functions

import (
	"testing"
)

func BenchmarkGreenFunctions(b *testing.B) {
	// Contoh parameter (ganti sesuai fungsi yang mau dites)
	x := 1.0
	y := 2.0
	z := 3.0

	b.ResetTimer() // reset timer sebelum loop benchmark
	for i := 0; i < b.N; i++ {
		_ = green_functions(x, y, z) // ganti dengan nama fungsi sebenarnya
	}
}
