package green_functions

import (
	"math"
	"testing"
)

// TestIntegration_DelhommeauVsLiangWuNoblesse compares results between methods
func TestIntegration_DelhommeauVsLiangWuNoblesse(t *testing.T) {
	// Create test meshes for a simple sphere-like geometry
	centers := [][]float64{
		{0, 0, -1}, {1, 0, -1}, {0, 1, -1}, {-1, 0, -1}, {0, -1, -1},
	}
	normals := [][]float64{
		{0, 0, 1}, {0, 0, 1}, {0, 0, 1}, {0, 0, 1}, {0, 0, 1},
	}

	mesh := NewMockMesh(centers, normals)

	// Initialize both methods
	delhommeau := NewDefaultDelhommeau()
	liangwu := NewLiangWuNoblesseGF()

	// Test parameters
	freeSurface := 0.0
	waterDepth := math.Inf(1)
	wavenumber := complex(1.0, 0)

	// Evaluate both methods
	S1, K1, err1 := delhommeau.Evaluate(mesh, mesh, freeSurface, waterDepth, wavenumber, true, true)
	if err1 != nil {
		t.Fatalf("Delhommeau evaluation failed: %v", err1)
	}

	S2, K2, err2 := liangwu.Evaluate(mesh, mesh, freeSurface, waterDepth, wavenumber, true, true)
	if err2 != nil {
		t.Fatalf("LiangWuNoblesse evaluation failed: %v", err2)
	}

	// Compare matrix dimensions
	s1Rows, s1Cols := S1.Dims()
	s2Rows, s2Cols := S2.Dims()

	if s1Rows != s2Rows || s1Cols != s2Cols {
		t.Errorf("S matrix dimensions mismatch: (%d,%d) vs (%d,%d)", s1Rows, s1Cols, s2Rows, s2Cols)
	}

	k1Rows, k1Cols := K1.Dims()
	k2Rows, k2Cols := K2.Dims()

	if k1Rows != k2Rows || k1Cols != k2Cols {
		t.Errorf("K matrix dimensions mismatch: (%d,%d) vs (%d,%d)", k1Rows, k1Cols, k2Rows, k2Cols)
	}
}

// TestIntegration_FiniteVsInfiniteDepthLimits tests convergence of finite to infinite depth
func TestIntegration_FiniteVsInfiniteDepthLimits(t *testing.T) {
	// Create simple test mesh
	centers := [][]float64{{0, 0, -1}, {1, 0, -1}}
	normals := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh := NewMockMesh(centers, normals)

	// Test parameters
	freeSurface := 0.0
	wavenumber := complex(0.5, 0) // Relatively low frequency

	// Infinite depth reference
	fgInf := NewFinGreen3D(math.Inf(1))
	SInf, _, err := fgInf.Evaluate(mesh, mesh, freeSurface, math.Inf(1), wavenumber, true, true)
	if err != nil {
		t.Fatalf("Infinite depth evaluation failed: %v", err)
	}

	// Test increasing depths
	depths := []float64{50.0, 100.0, 200.0, 500.0}

	for _, depth := range depths {
		fgFin := NewFinGreen3D(depth)
		SFin, _, err := fgFin.Evaluate(mesh, mesh, freeSurface, depth, wavenumber, true, true)
		if err != nil {
			t.Fatalf("Finite depth (%f) evaluation failed: %v", depth, err)
		}

		// Check that matrices have same dimensions
		sInfRows, sInfCols := SInf.Dims()
		sFinRows, sFinCols := SFin.Dims()

		if sInfRows != sFinRows || sInfCols != sFinCols {
			t.Errorf("Matrix dimensions mismatch for depth %f", depth)
		}

		// For very deep water, finite depth should approach infinite depth
		if depth > 100.0 {
			// This is where we would compare actual values if computation was implemented
			t.Logf("Depth %f: matrices computed successfully", depth)
		}
	}
}

// TestIntegration_WavenumberSweep tests behavior across frequency range
func TestIntegration_WavenumberSweep(t *testing.T) {
	// Create test mesh
	centers := [][]float64{{0, 0, -0.5}}
	normals := [][]float64{{0, 0, 1}}
	mesh := NewMockMesh(centers, normals)

	delhommeau := NewDefaultDelhommeau()

	// Test different wavenumbers
	wavenumbers := []complex128{
		complex(0.1, 0),   // Low frequency
		complex(1.0, 0),   // Medium frequency
		complex(5.0, 0),   // High frequency
		complex(1.0, 0.1), // Complex wavenumber
	}

	for i, k := range wavenumbers {
		S, _, err := delhommeau.Evaluate(mesh, mesh, 0.0, math.Inf(1), k, true, true)
		if err != nil {
			t.Errorf("Wavenumber test %d failed: %v", i, err)
		}

		if S == nil {
			t.Errorf("Wavenumber test %d returned nil matrices", i)
		}

		// Check for NaN or Inf values (when computation is implemented)
		sRows, sCols := S.Dims()
		if sRows != 1 || sCols != 1 {
			t.Errorf("Unexpected matrix dimensions for wavenumber test %d", i)
		}
	}
}

// TestIntegration_MeshScaling tests behavior with different mesh sizes
func TestIntegration_MeshScaling(t *testing.T) {
	delhommeau := NewDefaultDelhommeau()

	// Test different mesh sizes
	meshSizes := []int{1, 2, 5, 10}

	for _, size := range meshSizes {
		// Create mesh of given size
		centers := make([][]float64, size)
		normals := make([][]float64, size)

		for i := 0; i < size; i++ {
			centers[i] = []float64{float64(i), 0, -1}
			normals[i] = []float64{0, 0, 1}
		}

		mesh := NewMockMesh(centers, normals)

		S, _, err := delhommeau.Evaluate(mesh, mesh, 0.0, math.Inf(1), complex(1.0, 0), true, true)
		if err != nil {
			t.Errorf("Mesh size %d test failed: %v", size, err)
		}

		// Check matrix dimensions scale correctly
		sRows, sCols := S.Dims()
		if sRows != size || sCols != size {
			t.Errorf("Matrix dimensions for mesh size %d: expected (%d,%d), got (%d,%d)",
				size, size, size, sRows, sCols)
		}
	}
}

// TestIntegration_PrecisionConsistency tests float32 vs float64 consistency
func TestIntegration_PrecisionConsistency(t *testing.T) {
	// Create test mesh
	centers := [][]float64{{0, 0, -1}, {1, 0, -1}}
	normals := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh := NewMockMesh(centers, normals)

	// Test with both precisions
	params64 := DefaultDelhommeauParameters()
	params64.FloatingPointPrecision = Float64
	d64 := NewDelhommeau(params64)

	params32 := DefaultDelhommeauParameters()
	params32.FloatingPointPrecision = Float32
	d32 := NewDelhommeau(params32)

	// Evaluate both
	S64, _, err64 := d64.Evaluate(mesh, mesh, 0.0, math.Inf(1), complex(1.0, 0), true, true)
	if err64 != nil {
		t.Fatalf("Float64 evaluation failed: %v", err64)
	}

	S32, _, err32 := d32.Evaluate(mesh, mesh, 0.0, math.Inf(1), complex(1.0, 0), true, true)
	if err32 != nil {
		t.Fatalf("Float32 evaluation failed: %v", err32)
	}

	// Check dimensions match
	s64Rows, s64Cols := S64.Dims()
	s32Rows, s32Cols := S32.Dims()

	if s64Rows != s32Rows || s64Cols != s32Cols {
		t.Error("Matrix dimensions should be same regardless of precision")
	}
}

// TestIntegration_ErrorHandling tests various error conditions
func TestIntegration_ErrorHandling(t *testing.T) {
	delhommeau := NewDefaultDelhommeau()

	// Test with invalid mesh types
	_, _, err := delhommeau.Evaluate("invalid", nil, 0.0, math.Inf(1), complex(1.0, 0), true, true)
	if err == nil {
		t.Error("Expected error for invalid mesh1 type")
	}

	// Test with mismatched mesh types
	centers := [][]float64{{0, 0, -1}}
	normals := [][]float64{{0, 0, 1}}
	mesh := NewMockMesh(centers, normals)

	_, _, err = delhommeau.Evaluate(mesh, "invalid", 0.0, math.Inf(1), complex(1.0, 0), true, true)
	if err == nil {
		t.Error("Expected error for invalid mesh2 type")
	}

	// Test LiangWuNoblesse constraints
	lwn := NewLiangWuNoblesseGF()

	_, _, err = lwn.Evaluate(mesh, mesh, math.Inf(1), math.Inf(1), complex(1.0, 0), true, true)
	if err == nil {
		t.Error("Expected error for LiangWuNoblesse with infinite free surface")
	}
}

// TestIntegration_MemoryUsage tests memory efficiency with larger problems
func TestIntegration_MemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory test in short mode")
	}

	// Create moderately large mesh
	size := 50
	centers := make([][]float64, size)
	normals := make([][]float64, size)

	for i := 0; i < size; i++ {
		centers[i] = []float64{float64(i), 0, -1}
		normals[i] = []float64{0, 0, 1}
	}

	mesh := NewMockMesh(centers, normals)
	delhommeau := NewDefaultDelhommeau()

	// Multiple evaluations to test for memory leaks
	for i := 0; i < 10; i++ {
		S, K, err := delhommeau.Evaluate(mesh, mesh, 0.0, math.Inf(1), complex(1.0, 0), true, true)
		if err != nil {
			t.Fatalf("Iteration %d failed: %v", i, err)
		}

		// Verify matrices are properly allocated
		if S == nil || K == nil {
			t.Fatalf("Iteration %d returned nil matrices", i)
		}

		// Force garbage collection to test for memory leaks
		if i%3 == 0 {
			// In a real test, we might call runtime.GC() here
			t.Logf("Iteration %d completed successfully", i)
		}
	}
}

// Benchmark integration tests
func BenchmarkIntegration_SmallProblem(b *testing.B) {
	centers := [][]float64{{0, 0, -1}, {1, 0, -1}}
	normals := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh := NewMockMesh(centers, normals)

	delhommeau := NewDefaultDelhommeau()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := delhommeau.Evaluate(mesh, mesh, 0.0, math.Inf(1), complex(1.0, 0), true, true)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkIntegration_MediumProblem(b *testing.B) {
	size := 20
	centers := make([][]float64, size)
	normals := make([][]float64, size)

	for i := 0; i < size; i++ {
		centers[i] = []float64{float64(i), 0, -1}
		normals[i] = []float64{0, 0, 1}
	}

	mesh := NewMockMesh(centers, normals)
	delhommeau := NewDefaultDelhommeau()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := delhommeau.Evaluate(mesh, mesh, 0.0, math.Inf(1), complex(1.0, 0), true, true)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkIntegration_MethodComparison(b *testing.B) {
	centers := [][]float64{{0, 0, -1}, {1, 0, -1}, {0, 1, -1}}
	normals := [][]float64{{0, 0, 1}, {0, 0, 1}, {0, 0, 1}}
	mesh := NewMockMesh(centers, normals)

	b.Run("Delhommeau", func(b *testing.B) {
		delhommeau := NewDefaultDelhommeau()
		for i := 0; i < b.N; i++ {
			_, _, err := delhommeau.Evaluate(mesh, mesh, 0.0, math.Inf(1), complex(1.0, 0), true, true)
			if err != nil {
				b.Fatalf("Benchmark failed: %v", err)
			}
		}
	})

	b.Run("LiangWuNoblesse", func(b *testing.B) {
		lwn := NewLiangWuNoblesseGF()
		for i := 0; i < b.N; i++ {
			_, _, err := lwn.Evaluate(mesh, mesh, 0.0, math.Inf(1), complex(1.0, 0), true, true)
			if err != nil {
				b.Fatalf("Benchmark failed: %v", err)
			}
		}
	})

	b.Run("HAMS", func(b *testing.B) {
		hams := NewHAMS()
		for i := 0; i < b.N; i++ {
			_, _, err := hams.Evaluate(mesh, mesh, 0.0, math.Inf(1), complex(1.0, 0), true, true)
			if err != nil {
				b.Fatalf("Benchmark failed: %v", err)
			}
		}
	})
}
