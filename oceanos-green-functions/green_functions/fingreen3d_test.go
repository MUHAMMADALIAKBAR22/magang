package green_functions

import (
	"math"
	"testing"
)

func TestNewFinGreen3D(t *testing.T) {
	waterDepth := 10.0
	fg := NewFinGreen3D(waterDepth)

	if fg == nil {
		t.Fatal("NewFinGreen3D returned nil")
	}

	if fg.waterDepth != waterDepth {
		t.Errorf("Expected water depth %f, got %f", waterDepth, fg.waterDepth)
	}

	if fg.GetFloatingPointPrecision() != Float64 {
		t.Errorf("Expected Float64 precision, got %v", fg.GetFloatingPointPrecision())
	}
}

func TestFinGreen3D_String(t *testing.T) {
	fg := NewFinGreen3D(15.5)
	str := fg.String()

	expected := "FinGreen3D(water_depth=15.50)"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

func TestFinGreen3D_SetWaveNumber(t *testing.T) {
	fg := NewFinGreen3D(10.0)
	wavenumber := complex(1.0, 0.1)

	err := fg.SetWaveNumber(wavenumber)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if fg.waveNumber != wavenumber {
		t.Errorf("Expected wavenumber %v, got %v", wavenumber, fg.waveNumber)
	}

	if len(fg.dispersionRoots) == 0 {
		t.Error("Expected non-empty dispersion roots")
	}
}

func TestFinGreen3D_ComputeDispersionRoots(t *testing.T) {
	fg := NewFinGreen3D(10.0)
	wavenumber := complex(1.0, 0)

	roots, err := fg.computeDispersionRoots(wavenumber)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(roots) == 0 {
		t.Error("Expected at least one dispersion root")
	}

	// First root should be the input wavenumber
	if roots[0] != wavenumber {
		t.Errorf("Expected first root to be %v, got %v", wavenumber, roots[0])
	}

	// For finite depth, should have additional imaginary roots
	if len(roots) <= 1 {
		t.Error("Expected multiple roots for finite depth")
	}
}

func TestFinGreen3D_ComputeDispersionRoots_InfiniteDepth(t *testing.T) {
	fg := NewFinGreen3D(math.Inf(1))
	wavenumber := complex(1.0, 0)

	roots, err := fg.computeDispersionRoots(wavenumber)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// For infinite depth, should still get roots but different pattern
	if len(roots) == 0 {
		t.Error("Expected at least one root")
	}
}

func TestFinGreen3D_Evaluate(t *testing.T) {
	fg := NewFinGreen3D(20.0)

	// Create test meshes
	centers1 := [][]float64{{0, 0, -5}, {1, 0, -5}}
	normals1 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	centers2 := [][]float64{{0, 1, -5}, {1, 1, -5}}
	normals2 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh2 := NewMockMesh(centers2, normals2)

	// Test with finite depth
	S, K, err := fg.Evaluate(mesh1, mesh2, 0.0, 20.0, complex(1.0, 0), true, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if S == nil || K == nil {
		t.Error("Expected non-nil matrices")
	}

	// Check matrix dimensions
	sRows, sCols := S.Dims()
	if sRows != 2 || sCols != 2 {
		t.Errorf("Expected S matrix dimensions (2, 2), got (%d, %d)", sRows, sCols)
	}
}

func TestFinGreen3D_Evaluate_WaterDepthChange(t *testing.T) {
	fg := NewFinGreen3D(10.0)

	centers1 := [][]float64{{0, 0, -2}}
	normals1 := [][]float64{{0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	centers2 := [][]float64{{0, 1, -2}}
	normals2 := [][]float64{{0, 0, 1}}
	mesh2 := NewMockMesh(centers2, normals2)

	// Test with different water depth than initialized
	newDepth := 15.0
	_, _, err := fg.Evaluate(mesh1, mesh2, 0.0, newDepth, complex(1.0, 0), true, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should update internal water depth
	if fg.waterDepth != newDepth {
		t.Errorf("Expected water depth to update to %f, got %f", newDepth, fg.waterDepth)
	}
}

func TestFinGreen3D_ComputeGreenFunction3D_InfiniteDepth(t *testing.T) {
	fg := NewFinGreen3D(math.Inf(1))
	fg.SetWaveNumber(complex(1.0, 0))

	// Test points
	rr := 1.0  // horizontal distance
	zf := -1.0 // field point depth
	zp := -2.0 // source point depth

	result, err := fg.computeGreenFunction3D(rr, zf, zp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Result should be complex number
	if result == 0 {
		t.Error("Expected non-zero Green function value")
	}

	// Should have real part from Rankine terms
	if real(result) == 0 {
		t.Error("Expected non-zero real part from Rankine terms")
	}
}

func TestFinGreen3D_ComputeGreenFunction3D_FiniteDepth(t *testing.T) {
	fg := NewFinGreen3D(10.0)
	fg.SetWaveNumber(complex(1.0, 0))

	// Test points
	rr := 1.0  // horizontal distance
	zf := -1.0 // field point depth
	zp := -2.0 // source point depth

	result, err := fg.computeGreenFunction3D(rr, zf, zp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Result should be complex number
	if result == 0 {
		t.Error("Expected non-zero Green function value")
	}
}

func TestFinGreen3D_ComputeWaveTerm(t *testing.T) {
	fg := NewFinGreen3D(10.0)

	rr := 1.0
	zf := -1.0
	zp := -2.0
	k := complex(1.0, 0)

	// Test propagating mode
	term1 := fg.computeWaveTerm(rr, zf, zp, k, true)
	if term1 == 0 {
		t.Error("Expected non-zero wave term for propagating mode")
	}

	// Test evanescent mode
	kn := complex(0, 1.0)
	term2 := fg.computeWaveTerm(rr, zf, zp, kn, false)
	if term2 == 0 {
		t.Error("Expected non-zero wave term for evanescent mode")
	}
}

func TestFinGreen3D_ComputeInfiniteDepthGF_ZeroWavenumber(t *testing.T) {
	fg := NewFinGreen3D(math.Inf(1))
	fg.SetWaveNumber(complex(0, 0))

	rr := 1.0
	zf := -1.0
	zp := -2.0

	result, err := fg.computeInfiniteDepthGF(rr, zf, zp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// With zero wavenumber, should only have Rankine part
	expectedRankine := 1.0/math.Sqrt(rr*rr+(zf-zp)*(zf-zp)) + 1.0/math.Sqrt(rr*rr+(zf+zp)*(zf+zp))

	if math.Abs(real(result)-expectedRankine) > 1e-10 {
		t.Errorf("Expected Rankine value %f, got %f", expectedRankine, real(result))
	}

	if imag(result) != 0 {
		t.Errorf("Expected zero imaginary part, got %f", imag(result))
	}
}

// Benchmark tests
func BenchmarkFinGreen3D_SetWaveNumber(b *testing.B) {
	fg := NewFinGreen3D(10.0)
	wavenumber := complex(1.0, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := fg.SetWaveNumber(wavenumber)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkFinGreen3D_ComputeGreenFunction3D(b *testing.B) {
	fg := NewFinGreen3D(10.0)
	fg.SetWaveNumber(complex(1.0, 0))

	rr := 1.0
	zf := -1.0
	zp := -2.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := fg.computeGreenFunction3D(rr, zf, zp)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkFinGreen3D_Evaluate(b *testing.B) {
	fg := NewFinGreen3D(10.0)

	centers1 := [][]float64{{0, 0, -1}, {1, 0, -1}}
	normals1 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	centers2 := [][]float64{{0, 1, -1}, {1, 1, -1}}
	normals2 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh2 := NewMockMesh(centers2, normals2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := fg.Evaluate(mesh1, mesh2, 0.0, 10.0, complex(1.0, 0), true, true)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkFinGreen3D_InfiniteVsFiniteDepth(b *testing.B) {
	fgInf := NewFinGreen3D(math.Inf(1))
	fgFin := NewFinGreen3D(10.0)

	rr := 1.0
	zf := -1.0
	zp := -2.0

	fgInf.SetWaveNumber(complex(1.0, 0))
	fgFin.SetWaveNumber(complex(1.0, 0))

	b.Run("InfiniteDepth", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := fgInf.computeInfiniteDepthGF(rr, zf, zp)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
		}
	})

	b.Run("FiniteDepth", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := fgFin.computeFiniteDepthGF(rr, zf, zp)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
		}
	})
}
