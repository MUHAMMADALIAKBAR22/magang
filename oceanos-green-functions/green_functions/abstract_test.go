package green_functions

import (
	"gonum.org/v1/gonum/mat"
	"testing"
)

// MockMesh implements MeshLike interface for testing
type MockMesh struct {
	facesCenters *mat.Dense
	facesNormals *mat.Dense
	nbFaces      int
}

func NewMockMesh(centers, normals [][]float64) *MockMesh {
	nFaces := len(centers)
	centersData := make([]float64, 0, nFaces*3)
	normalsData := make([]float64, 0, nFaces*3)

	for i := 0; i < nFaces; i++ {
		centersData = append(centersData, centers[i]...)
		normalsData = append(normalsData, normals[i]...)
	}

	return &MockMesh{
		facesCenters: mat.NewDense(nFaces, 3, centersData),
		facesNormals: mat.NewDense(nFaces, 3, normalsData),
		nbFaces:      nFaces,
	}
}

func (m *MockMesh) GetFacesCenters() *mat.Dense { return m.facesCenters }
func (m *MockMesh) GetFacesNormals() *mat.Dense { return m.facesNormals }
func (m *MockMesh) GetNbFaces() int             { return m.nbFaces }

// TestBaseGreenFunction tests the base functionality
func TestBaseGreenFunction_Creation(t *testing.T) {
	bgf := NewBaseGreenFunction()

	if bgf.GetFloatingPointPrecision() != Float64 {
		t.Errorf("Expected default precision to be Float64, got %v", bgf.GetFloatingPointPrecision())
	}

	bgf.SetFloatingPointPrecision(Float32)
	if bgf.GetFloatingPointPrecision() != Float32 {
		t.Errorf("Expected precision to be Float32 after setting, got %v", bgf.GetFloatingPointPrecision())
	}
}

func TestBaseGreenFunction_GetColocationPointsAndNormals(t *testing.T) {
	bgf := NewBaseGreenFunction()

	// Test with MeshLike objects
	centers1 := [][]float64{{0, 0, 0}, {1, 0, 0}}
	normals1 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh1 := NewMockMesh(centers1, normals1)

	centers2 := [][]float64{{0, 1, 0}, {1, 1, 0}}
	normals2 := [][]float64{{0, 0, 1}, {0, 0, 1}}
	mesh2 := NewMockMesh(centers2, normals2)

	// Test non-adjoint case
	colocs, normals, err := bgf.getColocationPointsAndNormals(mesh1, mesh2, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if colocs == nil || normals == nil {
		t.Error("Expected non-nil collocation points and normals")
	}

	rows, cols := colocs.Dims()
	if rows != 2 || cols != 3 {
		t.Errorf("Expected collocation points dimensions (2, 3), got (%d, %d)", rows, cols)
	}
}

func TestBaseGreenFunction_InitMatrices(t *testing.T) {
	bgf := NewBaseGreenFunction()

	rows, cols := 3, 4

	// Test with early dot product
	S, K, err := bgf.initMatrices(rows, cols, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	sRows, sCols := S.Dims()
	if sRows != rows || sCols != cols {
		t.Errorf("Expected S matrix dimensions (%d, %d), got (%d, %d)", rows, cols, sRows, sCols)
	}

	kRows, kCols := K.Dims()
	expectedKCols := cols * 1 // early dot product = true
	if kRows != rows || kCols != expectedKCols {
		t.Errorf("Expected K matrix dimensions (%d, %d), got (%d, %d)", rows, expectedKCols, kRows, kCols)
	}

	// Test without early dot product
	_, K2, err := bgf.initMatrices(rows, cols, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	k2Rows, k2Cols := K2.Dims()
	expectedK2Cols := cols * 3 // early dot product = false
	if k2Rows != rows || k2Cols != expectedK2Cols {
		t.Errorf("Expected K2 matrix dimensions (%d, %d), got (%d, %d)", rows, expectedK2Cols, k2Rows, k2Cols)
	}
}

func TestBaseGreenFunction_InvalidInput(t *testing.T) {
	bgf := NewBaseGreenFunction()

	// Test with invalid mesh1 type
	_, _, err := bgf.getColocationPointsAndNormals("invalid", nil, false)
	if err == nil {
		t.Error("Expected error for invalid mesh1 type")
	}

	// Test with point array but invalid mesh2
	points := mat.NewDense(2, 3, []float64{0, 0, 0, 1, 1, 1})
	_, _, err = bgf.getColocationPointsAndNormals(points, "invalid", false)
	if err == nil {
		t.Error("Expected error for invalid mesh2 type with point array mesh1")
	}

	// Test with invalid point array dimensions
	invalidPoints := mat.NewDense(2, 2, []float64{0, 0, 1, 1})
	_, _, err = bgf.getColocationPointsAndNormals(invalidPoints, nil, false)
	if err == nil {
		t.Error("Expected error for invalid point array dimensions")
	}
}

// TestGreenFunctionEvaluationError tests custom error type
func TestGreenFunctionEvaluationError(t *testing.T) {
	err := &GreenFunctionEvaluationError{"test error message"}

	if err.Error() != "test error message" {
		t.Errorf("Expected error message 'test error message', got '%s'", err.Error())
	}
}

// Benchmark tests
func BenchmarkBaseGreenFunction_InitMatrices(b *testing.B) {
	bgf := NewBaseGreenFunction()
	rows, cols := 100, 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := bgf.initMatrices(rows, cols, true)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkBaseGreenFunction_GetColocationPoints(b *testing.B) {
	bgf := NewBaseGreenFunction()

	// Create larger test mesh
	centers := make([][]float64, 100)
	normals := make([][]float64, 100)
	for i := 0; i < 100; i++ {
		centers[i] = []float64{float64(i), 0, 0}
		normals[i] = []float64{0, 0, 1}
	}

	mesh1 := NewMockMesh(centers, normals)
	mesh2 := NewMockMesh(centers, normals)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := bgf.getColocationPointsAndNormals(mesh1, mesh2, false)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func TestGreenFunctionConsistency(t *testing.T) {
	// Create simple test for consistency
	bgf := NewBaseGreenFunction()

	// Test matrix initialization consistency
	S1, K1, err := bgf.initMatrices(2, 2, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Use variables to avoid unused variable error
	_ = S1
	_ = K1
}
