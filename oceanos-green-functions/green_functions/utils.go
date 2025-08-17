// Package green_functions - Utility functions and constants
// Copyright (C) 2025 Capytaine Contributors
// See LICENSE file at <https://github.com/capytaine/capytaine>

package green_functions

import (
	"math"
	"math/cmplx"
)

// Mathematical constants
const (
	Gravity = 9.81 // m/s^2
)

// ComputeDistance calculates the Euclidean distance between two 3D points
func ComputeDistance(p1, p2 [3]float64) float64 {
	dx := p1[0] - p2[0]
	dy := p1[1] - p2[1]
	dz := p1[2] - p2[2]
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// ComputeHorizontalDistance calculates the horizontal distance between two 3D points
func ComputeHorizontalDistance(p1, p2 [3]float64) float64 {
	dx := p1[0] - p2[0]
	dy := p1[1] - p2[1]
	return math.Sqrt(dx*dx + dy*dy)
}

// RankineSource computes the Rankine source potential
func RankineSource(r float64) float64 {
	if r == 0 {
		return math.Inf(1)
	}
	return 1.0 / r
}

// RankineSourceGradient computes the gradient of the Rankine source
func RankineSourceGradient(p1, p2 [3]float64) [3]float64 {
	dx := p1[0] - p2[0]
	dy := p1[1] - p2[1]
	dz := p1[2] - p2[2]
	r := math.Sqrt(dx*dx + dy*dy + dz*dz)

	if r == 0 {
		return [3]float64{math.Inf(1), math.Inf(1), math.Inf(1)}
	}

	r3 := r * r * r
	return [3]float64{-dx / r3, -dy / r3, -dz / r3}
}

// ComputeWaveNumber calculates the wave number from frequency
func ComputeWaveNumber(omega, waterDepth float64) complex128 {
	if math.IsInf(waterDepth, 1) {
		// Deep water approximation: k = omega^2/g
		k := omega * omega / Gravity
		return complex(k, 0)
	}

	// Finite depth: solve omega^2 = g*k*tanh(k*h)
	// Use Newton-Raphson iteration
	v := omega * omega / Gravity
	k := v // Initial guess

	for i := 0; i < 100; i++ {
		tanh_kh := math.Tanh(k * waterDepth)
		sech2_kh := 1.0 - tanh_kh*tanh_kh

		f := k*tanh_kh - v
		df := tanh_kh + k*waterDepth*sech2_kh

		if math.Abs(df) < 1e-15 {
			break
		}

		k_new := k - f/df

		if math.Abs(k_new-k) < 1e-12 {
			break
		}
		k = k_new
	}

	return complex(k, 0)
}

// DispersionRelation evaluates the dispersion relation
func DispersionRelation(k complex128, omega, waterDepth float64) complex128 {
	v := complex(omega*omega/Gravity, 0)
	kh := k * complex(waterDepth, 0)

	if math.IsInf(waterDepth, 1) {
		// Deep water: omega^2 = g*k
		return k - v
	}

	// Finite depth: omega^2 = g*k*tanh(k*h)
	return k*cmplx.Tanh(kh) - v
}

// PronyDecomposition represents Prony decomposition parameters
type PronyDecomposition struct {
	Coefficients []complex128
	Exponents    []complex128
}

// NewPronyDecomposition creates a new Prony decomposition
func NewPronyDecomposition(coeffs, exps []complex128) *PronyDecomposition {
	return &PronyDecomposition{
		Coefficients: coeffs,
		Exponents:    exps,
	}
}

// Evaluate evaluates the Prony decomposition at a given point
func (pd *PronyDecomposition) Evaluate(x float64) complex128 {
	var result complex128 = 0
	for i := range pd.Coefficients {
		result += pd.Coefficients[i] * cmplx.Exp(pd.Exponents[i]*complex(x, 0))
	}
	return result
}

// SingularityTreatment handles singularities in Green function computation
type SingularityTreatment int

const (
	HighFrequency SingularityTreatment = iota
	LowFrequency
	LowFrequencyWithRankine
)

// TabulationCache represents a cache for tabulated Green function values
type TabulationCache struct {
	RRange    []float64
	ZRange    []float64
	Values    [][]complex128
	IsValid   bool
	Precision FloatingPointPrecision
}

// NewTabulationCache creates a new tabulation cache
func NewTabulationCache(rRange, zRange []float64, precision FloatingPointPrecision) *TabulationCache {
	rows := len(zRange)
	cols := len(rRange)
	values := make([][]complex128, rows)
	for i := range values {
		values[i] = make([]complex128, cols)
	}

	return &TabulationCache{
		RRange:    rRange,
		ZRange:    zRange,
		Values:    values,
		IsValid:   false,
		Precision: precision,
	}
}

// Interpolate performs bilinear interpolation in the tabulation cache
func (tc *TabulationCache) Interpolate(r, z float64) (complex128, error) {
	if !tc.IsValid {
		return 0, &GreenFunctionEvaluationError{"Tabulation cache is not valid"}
	}

	// Find indices for interpolation
	rIdx := tc.findIndex(tc.RRange, r)
	zIdx := tc.findIndex(tc.ZRange, z)

	if rIdx < 0 || rIdx >= len(tc.RRange)-1 || zIdx < 0 || zIdx >= len(tc.ZRange)-1 {
		return 0, &GreenFunctionEvaluationError{"Point outside tabulation range"}
	}

	// Bilinear interpolation
	r1, r2 := tc.RRange[rIdx], tc.RRange[rIdx+1]
	z1, z2 := tc.ZRange[zIdx], tc.ZRange[zIdx+1]

	fr := (r - r1) / (r2 - r1)
	fz := (z - z1) / (z2 - z1)

	v11 := tc.Values[zIdx][rIdx]
	v12 := tc.Values[zIdx+1][rIdx]
	v21 := tc.Values[zIdx][rIdx+1]
	v22 := tc.Values[zIdx+1][rIdx+1]

	// Interpolate in r direction with proper type conversion
	v1 := v11*complex(1-fr, 0) + v21*complex(fr, 0)
	v2 := v12*complex(1-fr, 0) + v22*complex(fr, 0)

	// Interpolate in z direction
	result := v1*complex(1-fz, 0) + v2*complex(fz, 0)

	return result, nil
}

// findIndex finds the appropriate index for interpolation
func (tc *TabulationCache) findIndex(arr []float64, val float64) int {
	for i := 0; i < len(arr)-1; i++ {
		if val >= arr[i] && val <= arr[i+1] {
			return i
		}
	}
	return -1
}

// ValidateMatrixDimensions checks if matrix dimensions are compatible
func ValidateMatrixDimensions(rows, cols int) error {
	if rows <= 0 || cols <= 0 {
		return &GreenFunctionEvaluationError{"Matrix dimensions must be positive"}
	}
	if rows > 1000000 || cols > 1000000 {
		return &GreenFunctionEvaluationError{"Matrix dimensions too large"}
	}
	return nil
}
