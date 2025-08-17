package green_functions

import (
	"math"
	"reflect"
	"testing"
)

func TestDefaultDelhommeauParameters(t *testing.T) {
	params := DefaultDelhommeauParameters()

	expectedParams := DelhommeauParameters{
		TabulationNr:                        676,
		TabulationRmax:                      100.0,
		TabulationNz:                        372,
		TabulationZmin:                      -251.0,
		TabulationNbIntegrationPoints:       1001,
		TabulationGridShape:                 ScaledNemoh3,
		FiniteDepthMethod:                   NewerMethod,
		FiniteDepthPronyDecompositionMethod: PythonMethod,
		FloatingPointPrecision:              Float64,
		GfSingularities:                     LowFreq,
	}

	if !reflect.DeepEqual(params, expectedParams) {
		t.Errorf("Default parameters don't match expected values")
	}
}

func TestNewDelhommeau(t *testing.T) {
	params := DefaultDelhommeauParameters()
	d := NewDelhommeau(params)

	if d == nil {
		t.Fatal("NewDelhommeau returned nil")
	}

	if d.GetFloatingPointPrecision() != Float64 {
		t.Errorf("Expected Float64 precision, got %v", d.GetFloatingPointPrecision())
	}

	if d.parameters.TabulationNr != 676 {
		t.Errorf("Expected TabulationNr 676, got %d", d.parameters.TabulationNr)
	}
}

func TestNewDefaultDelhommeau(t *testing.T) {
	d := NewDefaultDelhommeau()

	if d == nil {
		t.Fatal("NewDefaultDelhommeau returned nil")
	}

	expectedParams := DefaultDelhommeauParameters()
	if !reflect.DeepEqual(d.parameters, expectedParams) {
		t.Error("Default Delhommeau doesn't use default parameters")
	}
}

func TestDelhommeau_String(t *testing.T) {
	// Test with default parameters
	d1 := NewDefaultDelhommeau()
	str1 := d1.String()
	if str1 != "Delhommeau()" {
		t.Errorf("Expected 'Delhommeau()' for default params, got '%s'", str1)
	}

	// Test with modified parameters
	params := DefaultDelhommeauParameters()
	params.TabulationNr = 500
	params.TabulationRmax = 200.0
	d2 := NewDelhommeau(params)
	str2 := d2.String()

	// Should show non-default values
	if str2 == "Delhommeau()" {
		t.Error("Expected non-default string representation for modified parameters")
	}
}

func TestDelhommeau_Hash(t *testing.T) {
	d1 := NewDefaultDelhommeau()
	d2 := NewDefaultDelhommeau()

	// Same parameters should give same hash
	if d1.Hash() != d2.Hash() {
		t.Error("Expected same hash for identical parameters")
	}

	// Different parameters should give different hash
	params := DefaultDelhommeauParameters()
	params.TabulationNr = 500
	d3 := NewDelhommeau(params)

	if d1.Hash() == d3.Hash() {
		t.Error("Expected different hash for different parameters")
	}
}

func TestDelhommeau_GridShapeIndices(t *testing.T) {
	params := DefaultDelhommeauParameters()

	// Test Legacy grid shape
	params.TabulationGridShape = Legacy
	d1 := NewDelhommeau(params)
	if d1.tabulationGridShapeIndex != 0 {
		t.Errorf("Expected Legacy grid shape index 0, got %d", d1.tabulationGridShapeIndex)
	}

	// Test ScaledNemoh3 grid shape
	params.TabulationGridShape = ScaledNemoh3
	d2 := NewDelhommeau(params)
	if d2.tabulationGridShapeIndex != 1 {
		t.Errorf("Expected ScaledNemoh3 grid shape index 1, got %d", d2.tabulationGridShapeIndex)
	}
}

func TestDelhommeau_FiniteDepthMethodIndices(t *testing.T) {
	params := DefaultDelhommeauParameters()

	// Test Legacy method
	params.FiniteDepthMethod = LegacyMethod
	d1 := NewDelhommeau(params)
	if d1.finiteDepthMethodIndex != 0 {
		t.Errorf("Expected Legacy method index 0, got %d", d1.finiteDepthMethodIndex)
	}

	// Test Newer method
	params.FiniteDepthMethod = NewerMethod
	d2 := NewDelhommeau(params)
	if d2.finiteDepthMethodIndex != 1 {
		t.Errorf("Expected Newer method index 1, got %d", d2.finiteDepthMethodIndex)
	}
}

func TestDelhommeau_GFSingularitiesIndices(t *testing.T) {
	params := DefaultDelhommeauParameters()

	// Test HighFreq
	params.GfSingularities = HighFreq
	d1 := NewDelhommeau(params)
	if d1.gfSingularitiesIndex != 0 {
		t.Errorf("Expected HighFreq index 0, got %d", d1.gfSingularitiesIndex)
	}

	// Test LowFreq
	params.GfSingularities = LowFreq
	d2 := NewDelhommeau(params)
	if d2.gfSingularitiesIndex != 1 {
		t.Errorf("Expected LowFreq index 1, got %d", d2.gfSingularitiesIndex)
	}

	// Test LowFreqWithRankinePart
	params.GfSingularities = LowFreqWithRankinePart
	d3 := NewDelhommeau(params)
	if d3.gfSingularitiesIndex != 2 {
		t.Errorf("Expected LowFreqWithRankinePart index 2, got %d", d3.gfSingularitiesIndex)
	}
}

func TestDelhommeau_Evaluate(t *testing.T) {
	d := NewDefaultDelhommeau()

	// Create test meshes
	centers1 := [][]float64{{0, 0, 0}, {1, 0, 0}}
	normals1 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	centers2 := [][]float64{{0, 1, 0}, {1, 1, 0}}
	normals2 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh2 := NewMockMesh(centers2, normals2)

	// Test evaluation
	S, K, err := d.Evaluate(mesh1, mesh2, 0.0, math.Inf(1), complex(1.0, 0), true, true)
	if err != nil {
		t.Fatalf("Unexpected error in Evaluate: %v", err)
	}

	if S == nil || K == nil {
		t.Error("Expected non-nil S and K matrices")
	}

	// Check matrix dimensions
	sRows, sCols := S.Dims()
	if sRows != 2 || sCols != 2 {
		t.Errorf("Expected S matrix dimensions (2, 2), got (%d, %d)", sRows, sCols)
	}

	kRows, kCols := K.Dims()
	expectedKCols := 2 * 1 // early dot product = true
	if kRows != 2 || kCols != expectedKCols {
		t.Errorf("Expected K matrix dimensions (2, %d), got (%d, %d)", expectedKCols, kRows, kCols)
	}
}

func TestDelhommeau_GetParameters(t *testing.T) {
	params := DefaultDelhommeauParameters()
	params.TabulationNr = 500
	d := NewDelhommeau(params)

	retrievedParams := d.GetParameters()
	if retrievedParams.TabulationNr != 500 {
		t.Errorf("Expected TabulationNr 500, got %d", retrievedParams.TabulationNr)
	}
}

// Benchmark tests for Delhommeau
func BenchmarkNewDelhommeau(b *testing.B) {
	params := DefaultDelhommeauParameters()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := NewDelhommeau(params)
		_ = d
	}
}

func BenchmarkDelhommeau_Evaluate_Small(b *testing.B) {
	d := NewDefaultDelhommeau()

	// Create small test meshes
	centers1 := [][]float64{{0, 0, 0}, {1, 0, 0}}
	normals1 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	centers2 := [][]float64{{0, 1, 0}, {1, 1, 0}}
	normals2 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh2 := NewMockMesh(centers2, normals2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := d.Evaluate(mesh1, mesh2, 0.0, math.Inf(1), complex(1.0, 0), true, true)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkDelhommeau_Evaluate_Medium(b *testing.B) {
	d := NewDefaultDelhommeau()

	// Create medium-sized test meshes
	nFaces := 50
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := d.Evaluate(mesh1, mesh2, 0.0, math.Inf(1), complex(1.0, 0), true, true)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
