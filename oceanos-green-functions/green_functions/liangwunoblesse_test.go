package green_functions

import (
	"math"
	"testing"
)

func TestNewLiangWuNoblesseGF(t *testing.T) {
	lwn := NewLiangWuNoblesseGF()

	if lwn == nil {
		t.Fatal("NewLiangWuNoblesseGF returned nil")
	}

	if lwn.GetFloatingPointPrecision() != Float64 {
		t.Errorf("Expected Float64 precision, got %v", lwn.GetFloatingPointPrecision())
	}

	expectedSettings := map[string]interface{}{
		"green_function": "LiangWuNoblesseGF",
	}

	if lwn.exportableSettings["green_function"] != expectedSettings["green_function"] {
		t.Error("Exportable settings don't match expected values")
	}
}

func TestLiangWuNoblesseGF_String(t *testing.T) {
	lwn := NewLiangWuNoblesseGF()
	str := lwn.String()

	expected := "LiangWuNoblesseGF()"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

func TestLiangWuNoblesseGF_Evaluate_Constraints(t *testing.T) {
	lwn := NewLiangWuNoblesseGF()

	// Create test meshes
	centers1 := [][]float64{{0, 0, 0}, {1, 0, 0}}
	normals1 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	centers2 := [][]float64{{0, 1, 0}, {1, 1, 0}}
	normals2 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh2 := NewMockMesh(centers2, normals2)

	// Test with infinite free surface - should fail
	_, _, err := lwn.Evaluate(mesh1, mesh2, math.Inf(1), math.Inf(1), complex(1.0, 0), true, true)
	if err == nil {
		t.Error("Expected error for infinite free surface")
	}

	// Test with finite water depth - should fail
	_, _, err = lwn.Evaluate(mesh1, mesh2, 0.0, 10.0, complex(1.0, 0), true, true)
	if err == nil {
		t.Error("Expected error for finite water depth")
	}

	// Test with valid conditions - should succeed
	_, _, err = lwn.Evaluate(mesh1, mesh2, 0.0, math.Inf(1), complex(1.0, 0), true, true)
	if err != nil {
		t.Errorf("Unexpected error for valid conditions: %v", err)
	}
}

func TestLiangWuNoblesseGF_Evaluate_WavenumberHandling(t *testing.T) {
	lwn := NewLiangWuNoblesseGF()

	// Create test meshes
	centers1 := [][]float64{{0, 0, 0}}
	normals1 := [][]float64{{0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	centers2 := [][]float64{{0, 1, 0}}
	normals2 := [][]float64{{0, 0, 1}}
	mesh2 := NewMockMesh(centers2, normals2)

	// Test with finite wavenumber
	S1, K1, err := lwn.Evaluate(mesh1, mesh2, 0.0, math.Inf(1), complex(1.0, 0), true, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if S1 == nil || K1 == nil {
		t.Error("Expected non-nil matrices")
	}

	// Test with infinite wavenumber
	S2, K2, err := lwn.Evaluate(mesh1, mesh2, 0.0, math.Inf(1), complex(math.Inf(1), 0), true, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if S2 == nil || K2 == nil {
		t.Error("Expected non-nil matrices")
	}
}

// Benchmark tests
func BenchmarkLiangWuNoblesseGF_Evaluate(b *testing.B) {
	lwn := NewLiangWuNoblesseGF()

	centers1 := [][]float64{{0, 0, 0}, {1, 0, 0}}
	normals1 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	centers2 := [][]float64{{0, 1, 0}, {1, 1, 0}}
	normals2 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh2 := NewMockMesh(centers2, normals2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := lwn.Evaluate(mesh1, mesh2, 0.0, math.Inf(1), complex(1.0, 0), true, true)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
} 