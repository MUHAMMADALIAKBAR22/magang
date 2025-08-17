package green_functions

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestComputeDistance(t *testing.T) {
	p1 := [3]float64{0, 0, 0}
	p2 := [3]float64{3, 4, 0}

	expected := 5.0
	result := ComputeDistance(p1, p2)

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("Expected distance %f, got %f", expected, result)
	}

	// Test with same points
	result = ComputeDistance(p1, p1)
	if result != 0 {
		t.Errorf("Expected distance 0 for same points, got %f", result)
	}
}

func TestComputeHorizontalDistance(t *testing.T) {
	p1 := [3]float64{0, 0, 5}
	p2 := [3]float64{3, 4, 10}

	expected := 5.0 // Only x,y components matter
	result := ComputeHorizontalDistance(p1, p2)

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("Expected horizontal distance %f, got %f", expected, result)
	}
}

func TestRankineSource(t *testing.T) {
	// Test normal case
	r := 2.0
	expected := 0.5
	result := RankineSource(r)

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("Expected Rankine source %f, got %f", expected, result)
	}

	// Test zero distance (singularity)
	result = RankineSource(0)
	if !math.IsInf(result, 1) {
		t.Error("Expected positive infinity for zero distance")
	}
}

func TestRankineSourceGradient(t *testing.T) {
	p1 := [3]float64{1, 0, 0}
	p2 := [3]float64{0, 0, 0}

	result := RankineSourceGradient(p1, p2)

	// For points (1,0,0) and (0,0,0), gradient should be (-1,0,0)
	expected := [3]float64{-1, 0, 0}

	for i := 0; i < 3; i++ {
		if math.Abs(result[i]-expected[i]) > 1e-10 {
			t.Errorf("Expected gradient component %d: %f, got %f", i, expected[i], result[i])
		}
	}

	// Test singularity case
	result = RankineSourceGradient(p1, p1)
	for i := 0; i < 3; i++ {
		if !math.IsInf(result[i], 0) {
			t.Errorf("Expected infinity for gradient component %d at singularity", i)
		}
	}
}

func TestComputeWaveNumber_DeepWater(t *testing.T) {
	omega := 2.0
	waterDepth := math.Inf(1)

	expected := complex(omega*omega/Gravity, 0)
	result := ComputeWaveNumber(omega, waterDepth)

	if cmplx.Abs(result-expected) > 1e-10 {
		t.Errorf("Expected wave number %v, got %v", expected, result)
	}
}

func TestComputeWaveNumber_FiniteDepth(t *testing.T) {
	omega := 2.0
	waterDepth := 10.0

	result := ComputeWaveNumber(omega, waterDepth)

	// For finite depth, result should be real and positive
	if imag(result) != 0 {
		t.Error("Expected real wave number for finite depth")
	}

	if real(result) <= 0 {
		t.Error("Expected positive wave number")
	}

	// Verify it satisfies dispersion relation
	k := real(result)
	v := omega * omega / Gravity
	dispersionCheck := k*math.Tanh(k*waterDepth) - v

	if math.Abs(dispersionCheck) > 1e-10 {
		t.Errorf("Wave number doesn't satisfy dispersion relation, error: %e", dispersionCheck)
	}
}

func TestDispersionRelation_DeepWater(t *testing.T) {
	k := complex(1.0, 0)
	omega := 2.0
	waterDepth := math.Inf(1)

	result := DispersionRelation(k, omega, waterDepth)
	expected := k - complex(omega*omega/Gravity, 0)

	if cmplx.Abs(result-expected) > 1e-10 {
		t.Errorf("Expected dispersion relation result %v, got %v", expected, result)
	}
}

func TestDispersionRelation_FiniteDepth(t *testing.T) {
	k := complex(1.0, 0)
	omega := 2.0
	waterDepth := 10.0

	result := DispersionRelation(k, omega, waterDepth)

	// Result should be complex
	if result == 0 {
		t.Error("Expected non-zero dispersion relation result")
	}
}

func TestNewPronyDecomposition(t *testing.T) {
	coeffs := []complex128{complex(1, 0), complex(2, 1)}
	exps := []complex128{complex(-1, 0), complex(-2, 0)}

	pd := NewPronyDecomposition(coeffs, exps)

	if pd == nil {
		t.Fatal("NewPronyDecomposition returned nil")
	}

	if len(pd.Coefficients) != len(coeffs) {
		t.Error("Coefficients length mismatch")
	}

	if len(pd.Exponents) != len(exps) {
		t.Error("Exponents length mismatch")
	}
}

func TestPronyDecomposition_Evaluate(t *testing.T) {
	// Simple case: f(x) = exp(-x)
	coeffs := []complex128{complex(1, 0)}
	exps := []complex128{complex(-1, 0)}

	pd := NewPronyDecomposition(coeffs, exps)

	// Test at x = 0
	result := pd.Evaluate(0)
	expected := complex(1, 0)

	if cmplx.Abs(result-expected) > 1e-10 {
		t.Errorf("Expected %v at x=0, got %v", expected, result)
	}

	// Test at x = 1
	result = pd.Evaluate(1)
	expected = complex(math.Exp(-1), 0)

	if cmplx.Abs(result-expected) > 1e-10 {
		t.Errorf("Expected %v at x=1, got %v", expected, result)
	}
}

func TestNewTabulationCache(t *testing.T) {
	rRange := []float64{0, 1, 2}
	zRange := []float64{-2, -1, 0}

	tc := NewTabulationCache(rRange, zRange, Float64)

	if tc == nil {
		t.Fatal("NewTabulationCache returned nil")
	}

	if len(tc.RRange) != len(rRange) {
		t.Error("R range length mismatch")
	}

	if len(tc.ZRange) != len(zRange) {
		t.Error("Z range length mismatch")
	}

	if tc.IsValid {
		t.Error("Expected cache to be invalid initially")
	}

	if tc.Precision != Float64 {
		t.Error("Precision mismatch")
	}

	// Check matrix dimensions
	if len(tc.Values) != len(zRange) {
		t.Error("Values matrix rows mismatch")
	}

	if len(tc.Values[0]) != len(rRange) {
		t.Error("Values matrix columns mismatch")
	}
}

func TestTabulationCache_Interpolate_Invalid(t *testing.T) {
	rRange := []float64{0, 1, 2}
	zRange := []float64{-2, -1, 0}

	tc := NewTabulationCache(rRange, zRange, Float64)

	// Should fail when cache is invalid
	_, err := tc.Interpolate(0.5, -1.5)
	if err == nil {
		t.Error("Expected error for invalid cache")
	}
}

func TestTabulationCache_Interpolate_OutOfRange(t *testing.T) {
	rRange := []float64{0, 1, 2}
	zRange := []float64{-2, -1, 0}

	tc := NewTabulationCache(rRange, zRange, Float64)
	tc.IsValid = true

	// Test out of range values
	_, err := tc.Interpolate(5.0, -1.5) // r out of range
	if err == nil {
		t.Error("Expected error for r out of range")
	}

	_, err = tc.Interpolate(0.5, 5.0) // z out of range
	if err == nil {
		t.Error("Expected error for z out of range")
	}
}

func TestTabulationCache_Interpolate_Valid(t *testing.T) {
	rRange := []float64{0, 1, 2}
	zRange := []float64{-2, -1, 0}

	tc := NewTabulationCache(rRange, zRange, Float64)
	tc.IsValid = true

	// Set up test values
	tc.Values[0][0] = complex(1, 0) // r=0, z=-2
	tc.Values[0][1] = complex(2, 0) // r=1, z=-2
	tc.Values[1][0] = complex(3, 0) // r=0, z=-1
	tc.Values[1][1] = complex(4, 0) // r=1, z=-1

	// Test interpolation at center
	result, err := tc.Interpolate(0.5, -1.5)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should be average of the four corner values
	expected := complex(2.5, 0)
	if cmplx.Abs(result-expected) > 1e-10 {
		t.Errorf("Expected interpolated value %v, got %v", expected, result)
	}
}

func TestTabulationCache_FindIndex(t *testing.T) {
	rRange := []float64{0, 1, 2, 3}
	zRange := []float64{-3, -2, -1, 0}

	tc := NewTabulationCache(rRange, zRange, Float64)

	// Test valid cases
	idx := tc.findIndex(rRange, 1.5)
	if idx != 1 {
		t.Errorf("Expected index 1 for value 1.5, got %d", idx)
	}

	idx = tc.findIndex(rRange, 0.0)
	if idx != 0 {
		t.Errorf("Expected index 0 for value 0.0, got %d", idx)
	}

	// Test out of range
	idx = tc.findIndex(rRange, 5.0)
	if idx != -1 {
		t.Errorf("Expected index -1 for out of range value, got %d", idx)
	}
}

func TestValidateMatrixDimensions(t *testing.T) {
	// Valid dimensions
	err := ValidateMatrixDimensions(10, 20)
	if err != nil {
		t.Errorf("Unexpected error for valid dimensions: %v", err)
	}

	// Invalid dimensions - zero
	err = ValidateMatrixDimensions(0, 10)
	if err == nil {
		t.Error("Expected error for zero rows")
	}

	err = ValidateMatrixDimensions(10, 0)
	if err == nil {
		t.Error("Expected error for zero columns")
	}

	// Invalid dimensions - negative
	err = ValidateMatrixDimensions(-1, 10)
	if err == nil {
		t.Error("Expected error for negative rows")
	}

	// Invalid dimensions - too large
	err = ValidateMatrixDimensions(2000000, 10)
	if err == nil {
		t.Error("Expected error for too large dimensions")
	}
}

// Benchmark tests
func BenchmarkComputeDistance(b *testing.B) {
	p1 := [3]float64{0, 0, 0}
	p2 := [3]float64{3, 4, 5}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ComputeDistance(p1, p2)
	}
}

func BenchmarkRankineSource(b *testing.B) {
	r := 2.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RankineSource(r)
	}
}

func BenchmarkComputeWaveNumber_FiniteDepth(b *testing.B) {
	omega := 2.0
	waterDepth := 10.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ComputeWaveNumber(omega, waterDepth)
	}
}

func BenchmarkTabulationCache_Interpolate(b *testing.B) {
	rRange := []float64{0, 1, 2, 3, 4, 5}
	zRange := []float64{-5, -4, -3, -2, -1, 0}

	tc := NewTabulationCache(rRange, zRange, Float64)
	tc.IsValid = true

	// Fill with test data
	for i := 0; i < len(zRange); i++ {
		for j := 0; j < len(rRange); j++ {
			tc.Values[i][j] = complex(float64(i+j), 0)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use values within range for interpolation
		r := 2.5  // Between 2 and 3
		z := -2.5 // Between -2 and -3
		_, err := tc.Interpolate(r, z)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
