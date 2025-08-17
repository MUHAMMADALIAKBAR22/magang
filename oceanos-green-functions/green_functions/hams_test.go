package green_functions

import (
	"math"
	"testing"
)

func TestNewHAMS(t *testing.T) {
	hams := NewHAMS()

	if hams == nil {
		t.Fatal("NewHAMS returned nil")
	}

	if hams.GetFloatingPointPrecision() != Float64 {
		t.Errorf("Expected Float64 precision, got %v", hams.GetFloatingPointPrecision())
	}

	expectedSettings := map[string]interface{}{
		"green_function": "HAMS",
	}

	if hams.exportableSettings["green_function"] != expectedSettings["green_function"] {
		t.Error("Exportable settings don't match expected values")
	}
}

func TestHAMS_String(t *testing.T) {
	hams := NewHAMS()
	str := hams.String()

	expected := "HAMS()"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

func TestHAMS_Evaluate(t *testing.T) {
	hams := NewHAMS()

	// Create test meshes
	centers1 := [][]float64{{0, 0, 0}, {1, 0, 0}}
	normals1 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	centers2 := [][]float64{{0, 1, 0}, {1, 1, 0}}
	normals2 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh2 := NewMockMesh(centers2, normals2)

	// Test evaluation with various parameters
	testCases := []struct {
		freeSurface float64
		waterDepth  float64
		wavenumber  complex128
	}{
		{0.0, math.Inf(1), complex(1.0, 0)},   // Infinite depth
		{0.0, 10.0, complex(1.0, 0)},          // Finite depth
		{0.0, math.Inf(1), complex(0.5, 0.1)}, // Complex wavenumber
	}

	for i, tc := range testCases {
		S, K, err := hams.Evaluate(mesh1, mesh2, tc.freeSurface, tc.waterDepth, tc.wavenumber, true, true)
		if err != nil {
			t.Errorf("Test case %d: unexpected error: %v", i, err)
		}

		if S == nil || K == nil {
			t.Errorf("Test case %d: expected non-nil matrices", i)
		}

		// Check matrix dimensions
		sRows, sCols := S.Dims()
		if sRows != 2 || sCols != 2 {
			t.Errorf("Test case %d: expected S matrix dimensions (2, 2), got (%d, %d)", i, sRows, sCols)
		}
	}
}

func TestHAMS_Evaluate_InvalidMesh(t *testing.T) {
	hams := NewHAMS()

	centers1 := [][]float64{{0, 0, 0}}
	normals1 := [][]float64{{0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	// Test with invalid mesh2
	_, _, err := hams.Evaluate(mesh1, "invalid", 0.0, math.Inf(1), complex(1.0, 0), true, true)
	if err == nil {
		t.Error("Expected error for invalid mesh2")
	}
}

// Benchmark tests
func BenchmarkHAMS_Evaluate(b *testing.B) {
	hams := NewHAMS()

	centers1 := [][]float64{{0, 0, 0}, {1, 0, 0}}
	normals1 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	centers2 := [][]float64{{0, 1, 0}, {1, 1, 0}}
	normals2 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh2 := NewMockMesh(centers2, normals2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := hams.Evaluate(mesh1, mesh2, 0.0, math.Inf(1), complex(1.0, 0), true, true)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkHAMS_vs_LiangWuNoblesse(b *testing.B) {
	hams := NewHAMS()
	lwn := NewLiangWuNoblesseGF()

	// Create larger test case
	nFaces := 20
	centers1 := make([][]float64, nFaces)
	normals1 := make([][]float64, nFaces)
	centers2 := make([][]float64, nFaces)
	normals2 := make([][]float64, nFaces)

	for i := 0; i < nFaces; i++ {
		centers1[i] = []float64{float64(i), 0, 0}
		normals1[i] = []float64{0, 0, 1}
		centers2[i] = []float64{float64(i), 1, 0}
		normals2[i] = []float64{0, 0, 1}
	}

	mesh1 := NewMockMesh(centers1, normals1)
	mesh2 := NewMockMesh(centers2, normals2)

	b.Run("HAMS", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := hams.Evaluate(mesh1, mesh2, 0.0, math.Inf(1), complex(1.0, 0), true, true)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
		}
	})

	b.Run("LiangWuNoblesse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := lwn.Evaluate(mesh1, mesh2, 0.0, math.Inf(1), complex(1.0, 0), true, true)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
		}
	})
} 